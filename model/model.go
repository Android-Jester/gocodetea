package model

import (
	"bytes"
	"errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"

	// learning "example.com/learning"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	// activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	activeTabBorder  = lipgloss.OuterHalfBlockBorder()
	docStyle         = lipgloss.NewStyle().Padding(1, 2)
	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle   = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	// windowStyle      = lipgloss.NewStyle().BorderForeground(highlightColor).Align(lipgloss.Left)
	// .Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

type Stack struct {
	mutex sync.Mutex
	s     []func()
}

func NewStack() *Stack {
	return &Stack{mutex: sync.Mutex{}, s: make([]func(), 0)}
}
func (s *Stack) Push(data func()) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.s = append(s.s, data)
}
func (s *Stack) Pop() (func(), error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	l := len(s.s)
	if l == 0 {
		return nil, errors.New("empty stack")
	}
	last := s.s[l-1]
	s.s = s.s[:l-1]
	return last, nil
}
func (s *Stack) IsEmpty() bool {
	return len(s.s) == 0
}
func (s *Stack) Top() (func(), error) {
	l := len(s.s)
	if l == 0 {
		return nil, errors.New("no element")
	}
	last := s.s[l-1]
	return last, nil
}

// bubbletea Model for UI
type Model struct {
	Tabs            []string
	TabContent      []string
	activeTab       int
	codeView        bool
	Stack           *Stack //func returned from functions that run when you change tabs
	viewport        viewport.Model
	SourceCode      string
	MethodContainer interface{} //contains a struct with method receivers that the model will execute
}

func (m Model) renderTabs() string {
	var renderedTabs []string
	for i, t := range m.Tabs {
		var style lipgloss.Style
		_, _, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		// } else if isLast && isActive {
		// 	// border.BottomRight = "│"
		// }
		// else if isLast && !isActive {
		// 	border.BottomRight = "┤"
		// }
		border, _, _, _, _ := style.GetBorder()
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

// Return a string containing the tabs that will serve as header
func (m Model) headerView() string {
	doc := strings.Builder{}
	row := m.renderTabs()
	doc.WriteString(row)
	doc.WriteString("\n")
	return docStyle.Render(doc.String())
}

func (m Model) View() string {
	doc := strings.Builder{}
	row := m.renderTabs()
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(m.viewport.View())
	// doc.WriteString(m.TabContent[m.activeTab])
	// doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.TabContent[m.activeTab]))
	return docStyle.Render(doc.String())
}
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// equivalent to value of say learning.GoStruct{}
	val := reflect.ValueOf(m.MethodContainer)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			for { //pop every item in the stack and run items until stack is empty
				f, err := m.Stack.Pop()
				// fmt.Println(&f, err)
				if err != nil {
					// fmt.Println(&f, err)
					break
				}
				f()
			}
			// fmt.Println("popping now", len(m.deferredFuncs.s))
			return m, tea.Quit
		case "c": //change view to show code or not
			m.codeView = !m.codeView
			m.setTabContentToFnOutput(val)
			m.viewport.SetContent(m.TabContent[m.activeTab])
			return m, nil
		case "right", "tab":
			// fmt.Println(m.Tabs[m.activeTab])
			//check last item in stack, if fn is nil implies that
			//tabContent didn't return a cleanup function,
			//codeView shouldnt execute cleanup functions as well
			// fmt.Println(m.deferredFuncs.s)
			if fn, _ := m.Stack.Top(); fn != nil && !m.codeView {
				// m.deferredFuncs()
				// fmt.Println(err)
				// fmt.Println("This runnn")
				fn()
				m.Stack.Pop()
			}
			// fmt.Println("Always runs")
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			m.setTabContentToFnOutput(val)
			m.viewport.SetContent(m.TabContent[m.activeTab])
			return m, nil
		case "left", "shift+tab":
			// fmt.Println(m.Tabs[m.activeTab])
			if fn, _ := m.Stack.Top(); fn != nil && !m.codeView {
				fn()
				m.Stack.Pop()
			}
			m.activeTab = max(m.activeTab-1, 0)
			m.setTabContentToFnOutput(val)
			m.viewport.SetContent(m.TabContent[m.activeTab])
			return m, nil
		default:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		s := m.headerView()
		headerHeight := lipgloss.Height(s)
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - headerHeight
	}
	//trying to set initialContent
	m.setTabContentToFnOutput(val)
	m.viewport.SetContent(m.TabContent[m.activeTab])
	return m, nil
}

// Sets the model's TabContent property based on whether the codeView
// is enabled or not. If codeView is true, it sets the TabContent to the formatted body of
// the selected function. Otherwise, it sets it to the output of the function.
func (m *Model) setTabContentToFnOutput(val reflect.Value) {
	tab := m.Tabs[m.activeTab]
	val = val.MethodByName(tab)

	if m.codeView {
		f, fs := getFuncAST(tab, "", m.SourceCode)
		body := getFuncBodyStr(f, fs)
		m.TabContent[m.activeTab] = body
	} else {
		content := m.funcToStdOut(val)
		m.TabContent[m.activeTab] = content
	}
}
func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

// Takes a reflect.Value representing a function, executes it, and captures its
// standard output. It returns the standard output as a string.
// If the function being executed returns a closure, it stores that closure in the model.
func (m *Model) funcToStdOut(f reflect.Value) string {
	oldStdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	if f.Kind() != reflect.Func {
		panic("Not a function")
	}

	switch f.Type().String() {
	case "func() func()":
		newf := f.Interface().(func() func())
		closure := newf()
		// m.deferredFuncs[m.activeTab] = closure
		m.Stack.Push(closure)

	default:
		f.Call([]reflect.Value{})
		// m.deferredFuncs[m.activeTab] = func() {}
		m.Stack.Push(func() {})

	}

	outC := make(chan string)

	go func() { //concurrent writing of stdout to buf then sending along the channel
		var buf strings.Builder
		io.Copy(&buf, r)
		outC <- buf.String()
		buf.Reset()
	}()
	w.Close()
	os.Stdout = oldStdOut
	out := <-outC
	// fmt.Print(out, outC, oldStdOut)
	// fmt.Println("previous output:")
	// fmt.Print(out)
	return out
}

// Takes an *ast.FuncDecl and a *token.FileSet and returns the body of the
// function as a string. It uses the go/format package to format the body.
func getFuncBodyStr(f *ast.FuncDecl, fs *token.FileSet) string {
	var buf bytes.Buffer
	// printer := printer.Config{Mode: printer.RawFormat}
	if err := format.Node(&buf, fs, f.Body); err != nil {
		panic(err)
	}
	return buf.String()
}

// Parses a function AST from a source code or filename.
// Uses src if provided or fallback to filename
// It looks for a function with the specified name and
// returns its *ast.FuncDecl along with the corresponding *token.FileSet.
func getFuncAST(funcname string, filename string, src any) (*ast.FuncDecl, *token.FileSet) {
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, filename, src, 0)
	if err != nil {
		panic(err)
	}
	//loop through all func decl to check if
	//your func exists
	for _, d := range file.Decls {
		if f, ok := d.(*ast.FuncDecl); ok && f.Name.Name == funcname {
			return f, fs
		}
	}
	panic("function not found")

}
