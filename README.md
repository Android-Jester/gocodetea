# Go Code Viewer with Bubble Tea

This Go code demonstrates a simple code viewer application built using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework. The application allows you to view and interact with Go code snippets from various functions provided by the `learning.Go_Struct` struct defined in the `example.com/learning` package.
Execution is essentially calling and piping the output of a function from stdout to a string

## Features

- View and switch between different Go code snippets.
- Navigate through code snippets using keyboard shortcuts.
- Toggling between code view and execution view.
- Supports code with functions that return deferred functions for cleanup.

## Prerequisites

Before running this code, make sure you have Go installed on your system.

## Usage

1. Clone the repository or download the code files to your local machine.

2. Navigate to the directory containing the code files in your terminal.

3. Run the application using the following command:

   ```bash
   go run .
   ```
   or perhaps 
   ```bash 
   go build -o appname
   ```

4. The application will launch, and you will see a list of tabs representing different Go code snippets.

5. Use the following keyboard shortcuts to navigate and interact with the application:

   - **Right Arrow / Tab:** Switch to the next tab.
   - **Left Arrow / Shift+Tab:** Switch to the previous tab.
   - **c:** Toggle between viewing the code and executing it.
   - **Ctrl+C / q:** Quit the application.

6. Explore the Go code snippets, execute them, and view their output.

## Code Structure

- `main.go`: The main entry point of the application.
- `model.go`: Defines the application model and update functions.
- `learning.go`: Imports the `example.com/learning` package and interacts with its functions.


## License

This code is open-source and available under the [MIT License](LICENSE). Feel free to modify and use it for your projects.

If you have any questions or feedback, please don't hesitate to reach out. Enjoy exploring and experimenting with Go code!