# Go Code Viewer with Bubble Tea

This Go code demonstrates a simple code viewer application built using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework. 
The application allows you to view and interact with Go code snippets from various functions provided by your own Go modules. You can easily integrate your own code by calling your embedded module into the `main.go` file.

## Features

- View and switch between different Go code snippets.
- Navigate through code snippets using keyboard shortcuts.
- Toggling between code view and execution view.
- Supports code with functions that return deferred functions for cleanup.

Execution is essentially calling and piping the output of a function from stdout to a string


## Overview

The `model` package is responsible for defining and managing the application's user interface (UI) model. It provides the underlying data structure and logic for managing and displaying Go code snippets.

It encapsulates the state of the application, including information about code tabs, currently selected code snippets, and user interactions.

### Key Components

The `model` package consists of the following key components:

- **Model:** The `Model` struct defines the application's main model. It contains properties such as `Tabs`, `TabContent`, `activeTab`, `codeView`, `Stack`, `viewport`, `SourceCode`, and `MethodContainer`.

  - `Tabs`: An array of strings representing different method names as tabs.
  - `TabContent`: An array of strings representing the content of each code tab.
  - `Stack`: A stack data structure that stores deferred functions for cleanup.
  - `SourceCode`: The source code of the Go module containing the code snippets.
  - `MethodContainer`: A placeholder for a struct with method receivers that define the code snippets.


## Prerequisites

Before running this code, make sure you have Go installed on your system.

## Usage

To use the `model` package in your Go Code Viewer application, follow these steps:
- Create your own module 
- Import the `model` package into your main file:

   ```go
   import "github.com/lazarusking/gocodetea/model"
   ```
- In some other module write this at the top of the `.go` file containing the methods you want to display, using the go embed directive(after any imports, the `var` can have any name) 
   - (the learning folder serves as an example for how it could be structured.)
   ```go
   //go:embed filename.go
   var CodeSrc string 

   //this embeds the src code into a string CodeSrc
   ```

   in your `main.go` file: initialize your model with this
   ```go
   func main() {

   var learn example.Go_Struct // rename this to the package with the Struct containing method receivers

	//getting the types of the methods
	t := reflect.TypeOf(learn)

	var learning_tabs []string
	for i := 0; i < t.NumMethod(); i++ {
		methodName := t.Method(i).Name
		learning_tabs = append(learning_tabs, methodName)
	}
	//populate model with initial fields
	m := model.Model{Tabs: learning_tabs,
		TabContent:      make([]string, len(learning_tabs)),
		Stack:   model.NewStack(),
		SourceCode:      learning.CodeSrc,
		MethodContainer: learn}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
   }
   ```
   A sample of the other module
   ```go
   package example

   import (
      _ "embed"
      "fmt"
   )

   //go:embed code.go
   var CodeSrc string

   type Go_Struct struct{}

   func (Go_Struct) Variables() {
      fmt.Println("\nVariable Declaration and Initialization:")
      var foo string = "Hello"
      const bar string = `World`
      foo2 := "Shorthand"
      tabs := []string{"some", "test"}

      fmt.Println(foo, bar, tabs)
      fmt.Println(foo2)
      percent := (7.0 / 9) * 100
      fmt.Printf("%.2f %x\n", percent, 10)
      s := fmt.Sprintf("hex:%x bin:%b", 10, 10)
      fmt.Println(s)
   }
   ```

   ```
   someothermodule/
   │   └── code.go  # file containing methods to display
   
   mainproject/
   ├── main.go   # This is the main program that uses the model package
   ├── go.mod
   └── go.sum
   ```


## Cloning
To test how the repo works or for contribution:
1. Clone the repository or download the code files to your local machine.

2. Navigate to the directory containing the code files in your terminal.

3. Run the application using the following command:

   ```bash
   go run .
   ```
   or perhaps 
   ```bash 
   go build
   ```

4. The application will launch, and you will see a list of tabs representing different Go code snippets.

## Keybindings 
 Use the following keyboard shortcuts to navigate and interact with the application:

   - **Right Arrow / Tab:** Switch to the next tab.
   - **Left Arrow / Shift+Tab:** Switch to the previous tab.
   - **c:** Toggle between viewing the code and executing it.
   - **Ctrl+C / q:** Quit the application.

6. Explore the Go code snippets, execute them, and view their output.

## Code Structure

- `main.go`: The main entry point of the application.
- `model.go`: Defines the application model and update functions.
- `learning`: An example module package that the main file interacts with its functions.


## License

This code is open-source and available under the [MIT License](LICENSE). Feel free to modify and use it for your projects.

If you have any questions or feedback, please don't hesitate to reach out. Enjoy exploring and experimenting with Go code!