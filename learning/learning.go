package learning

import (
	_ "embed"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed learning.go
var CodeSrc string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

type geometry interface {
	area() float64
	perim() float64
}
type GoLearner interface {
	Variables()
	Conditionals(rect)
	Loops()
	Functions()
	Closures()
}

type Go_Struct struct{}

func (Go_Struct) Variables() {
	colorFn("\nVariable Declaration and Initialization:")

	// Assignments and Initializations
	// OR
	// var (
	// 	foo = "Hello"
	// 	bar string = "World"
	// )
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
func (Go_Struct) Conditionals() {
	conditionals()
}
func (Go_Struct) Loops() {
	loops()
}
func (Go_Struct) Functions() {
	for i := 0; i < 4; i++ {
		fmt.Println(fib(i))
	}
	fmt.Println(myFunction("Hello", "World"))
}
func (Go_Struct) Closures() { // closure
	colorFn("\nClosures:")
	add := add_closure()
	add(10)
	add(10)
	fmt.Println(add(2))
}
func (receiver Go_Struct) FileIO() func() {
	os.Mkdir("test", 0666)
	os.MkdirAll("test/bruh/some/people/", 0755)
	tempf, _ := os.Create("newf.txt")
	folder, _ := os.ReadDir(".")

	dirHeader := []interface{}{"Mode", "Time", "Size", "Name"}
	var dash []interface{}
	for i := 0; i < len(dirHeader); i++ {
		dash = append(dash, "----")
	}

	fmt.Printf("%-10s \t %21s \t %6s \t %6s\n", dirHeader...)
	fmt.Printf("%-10s \t %21s \t %6s \t %6s\n", dash...)

	// filepath.WalkDir("learning", func(path string, d fs.DirEntry, err error) error {
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return err
	// 	}
	// 	if d.IsDir() {
	// 		return nil
	// 	}
	// 	f, _ := d.Info()
	// 	fmt.Printf("%6s \t %6v \t %6v \t %6s\n", f.Mode(), f.IsDir(), f.Size(), d.Name())
	// 	return nil
	// })

	for _, d := range folder {
		if strings.HasPrefix(d.Name(), ".") {
			continue
		}
		printDirContent("", d, "")
	}
	return func() {
		//function that is exec when you switch tabs
		// fmt.Println("run deferred func")
		defer os.RemoveAll("test")

		defer func() {
			// 	tempf.Close()
			if err := os.Remove(tempf.Name()); err != nil {
				fmt.Println(err)
			}
			// 	check(os.Remove(tempf.Name()))
		}()
		defer tempf.Close()
	}
}
func printDirContent(path string, d os.DirEntry, indent string) {
	f, _ := d.Info()
	t := fmt.Sprintf("%v", time.Time.Format(f.ModTime(), "01/02/2006 03:04:05PM"))
	strfmt := func(s string) string {
		c := []string{colorBlue, colorRed, colorGreen, colorPurple, colorYellow, colorCyan}
		rndcol := c[rand.Intn(len(c))]
		return fmt.Sprint(rndcol, s, colorReset)
	}
	dirSymbol := ""
	newPath := filepath.Join(path, d.Name())
	before, _, IsNested := strings.Cut(newPath, string(filepath.Separator))
	// ss := fmt.Sprintln(before, "b", after, "a", "p", newPath, " bap")
	// colorFn(ss)
	if d.IsDir() {
		var newIndent strings.Builder
		newIndent.WriteString(indent)
		if IsNested {
			dirSymbol = strfmt("└")
		}
		if before == newPath {
			// dirSymbol = ""
			// fmt.Println("before", before, path)
			newIndent.Reset()
			newIndent.WriteString(" ")
			// fmt.Println("It is nested")
		} else {
			// newIndent.WriteString(indent)
			newIndent.WriteString(" ")
		}
		fmt.Printf("%6s \t %-21v \t %6v \t%s%s%-6s folder loop\n", f.Mode(), t, f.Size(), newIndent.String(), dirSymbol, (colorCyan + d.Name() + colorReset))
		// fmt.Println(len(indent))

		dir, err := os.ReadDir(newPath)
		if err != nil {
			fmt.Println(err)
		}
		for _, newdir := range dir {
			printDirContent(newPath, newdir, newIndent.String())
		}

	} else {
		var newIndent strings.Builder
		newIndent.WriteString(indent)

		if IsNested {
			dirSymbol = strfmt("└")
		} else if before == newPath {
			dirSymbol = ""
			// fmt.Println("before", before, path)
			newIndent.WriteString("")
		} else {
			newIndent.WriteString(indent + " ")
		}
		fmt.Printf("%6s \t %6v \t %6v \t %s%s%-6s \n", f.Mode(), t, f.Size(), newIndent.String(), dirSymbol, d.Name())
	}
}
func fib(n int) int {
	if n < 2 {
		return n
	} else {
		// return fib(n-1) + fib(n-2)
		first, second := 0, 1
		for i := 2; i < n+1; i++ {
			third := first + second
			first, second = second, third
			// fmt.Println(first, second, third)
		}
		return second
	}
}

type rect struct {
	width, height float64
}
type circle struct {
	radius float64
}

func (r rect) area() float64 {
	return r.width * r.height
}
func (r rect) perim() float64 {
	return 2*r.width + 2*r.height
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}
func (c circle) perim() float64 {
	return 2 * math.Pi * c.radius
}

func measure(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.perim())
}

func Features(r rect) {
	lg := Go_Struct{}
	lg.Variables()
	lg.Loops()
	lg.Functions()
	lg.Closures()
	lg.Conditionals()
}

func main() {
	r := rect{width: 3, height: 4}
	c := circle{radius: 5}
	// lg := go_struct{}

	measure(r)
	measure(c)
	Features(r)
}

type colorFnConfig struct {
	color string
	text  string
}
type colorFnOption func(*colorFnConfig)

func withcolor(color string) colorFnOption {
	return func(cfg *colorFnConfig) {
		cfg.color = color
	}
}

func colorFn(text string, options ...colorFnOption) any {
	cfg := &colorFnConfig{
		color: colorRed,
		text:  text,
	}
	for _, option := range options {
		option(cfg)
	}
	s := fmt.Sprintln(cfg.color, cfg.text, string(colorReset))
	fmt.Print(s)
	return s
}

func loops() {
	colorFn("\nLooping Structures:")
	for i := 0; i < 10; i++ {
		if i < 2 {
			continue
		}
		s := fmt.Sprintf("%d", i)
		fmt.Print(s)
		if i > 5 {
			break
		}
	}
	i := 0
	for i < 10 { // no init nor post statement; just condition
		fmt.Print(i)
		i++
	}
	func() { //anonymous function
		fmt.Println("\nIt is Finished")
	}()
}

func conditionals() {
	r := rect{width: 3, height: 4}

	colorFn("\nConditional Statements:")
	if r.area() > 30 {
		fmt.Println("area is > 30")
	} else {
		fmt.Println("area is < 30")
	}

	if x := 10; x > 5 { //shorthand
		fmt.Println("x is >5")
	}

	switch day := "monday"; day { //shorthand
	case "monday":
		fmt.Println("Let's learn Golang")
	case "friday":
		fmt.Println("let's party")
	default:
		fmt.Println("browse memes")
	}
}

func closure() { // closure
	colorFn("\nClosures:")
	add := add_closure()
	add(10)
	add(10)
	fmt.Println(add(2))
}

func myFunction(s1, s2 string) (string, string) {
	colorFn("\nFunctions and Anonymous Functions:")

	//multiple arguments variadic
	sum := add(1, 2, 3, 4, 5)
	fmt.Println("sum is", sum)

	s := fmt.Sprintf("%s %s", s1, s2)
	fn := func() string {
		return "...."
	}
	return s, fn()
}
func add_closure() func(int) int {
	sum := 0
	return func(v int) int { //sum is scoped to this func
		sum += v
		return sum
	}
}
func add(values ...int) int {
	sum := 0
	for _, v := range values {
		sum += v
	}
	return sum
}
