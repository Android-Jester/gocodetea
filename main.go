package main

import (
	_ "embed"
	"fmt"
	"os"
	"reflect"

	"github.com/lazarusking/gocodetea/learning"
	"github.com/lazarusking/gocodetea/model"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
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
		learning_tabs = append(learning_tabs, methodName)
	}
	// tabs := []string{"Lip Gloss", "Blush", "Eye Shadow", "Mascara", "Foundation"}
	// tabContent := []string{"Lip Gloss Tab", "Blush Tab", "Eye Shadow Tab", "Mascara Tab", "Foundation Tab"}

	//populate model with initial fields
	m := model.Model{Tabs: learning_tabs,
		TabContent:      make([]string, len(learning_tabs)),
		DeferredFuncs:   model.NewStack(),
		SourceCode:      learning.CodeSrc,
		MethodContainer: learn}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
