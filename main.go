package main

import (
	_ "embed"
	"fmt"
	"os"
	"reflect"

	learning "example.com/learning"
	"github.com/charmbracelet/bubbles/viewport"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var width = 78
	var height = 20
	// viewport size
	vp := viewport.New(width, height)

	// declare zero-valued struct
	var learn learning.Go_Struct
	// var ln = &learning.Go_Struct{} //using pointers

	//getting the types of the methods
	t := reflect.TypeOf(learn)
	// t := reflect.TypeOf((*learning.GoLearner)(nil)).Elem()
	// t := reflect.TypeOf(ln)
	// t := reflect.TypeOf(struct{learning.Go_Struct}{})

	var learning_tabs []string
	// var tabContent = make([]string, 0)
	//append names of interface/struct methods to array
	for i := 0; i < t.NumMethod(); i++ {
		methodName := t.Method(i).Name
		// f, fs := getFuncAST(methodName, "learning/learning.go", learning.CodeSrc)
		// body := getFuncBodyStr(f, fs)
		// val := reflect.ValueOf(learn).MethodByName(methodName)
		// fmt.Println(val)
		learning_tabs = append(learning_tabs, methodName)
		// s, _ := funcToStdOut(val)
		// fmt.Println(s)
		// tabContent = append(tabContent, s)
		if i == 0 {
			vp.SetContent("")
		}
	}

	// tabs := []string{"Lip Gloss", "Blush", "Eye Shadow", "Mascara", "Foundation"}
	// tabContent := []string{funcToStdOut(learn.Closures), "Blush Tab", funcToStdOut(learn.Functions), "Mascara Tab", "Foundation Tab"}
	// tabContent := []string{"Lip Gloss Tab", "Blush Tab", "Eye Shadow Tab", "Mascara Tab", "Foundation Tab"}

	//populate model with initial fields
	m := model{Tabs: learning_tabs, TabContent: make([]string, len(learning_tabs)), viewport: vp, deferredFuncs: NewStack()}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
