package file

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/kenshindeveloper/april/evaluator"
	"github.com/kenshindeveloper/april/lexer"
	"github.com/kenshindeveloper/april/object"
	"github.com/kenshindeveloper/april/parser"
	"github.com/kenshindeveloper/april/repl"
)

func Start(name string) bool {

	file, _ := ioutil.ReadFile(name)

	input := string(file)

	// fmt.Println("------------------------------------")
	// fmt.Println(input)
	// fmt.Println("------------------------------------\n")

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParserProgram()
	if len(p.Error()) != 0 {
		repl.PrintParserErrors(os.Stdout, p.Error())
	}
	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(os.Stdout, evaluated.Inspect())
	}

	return true
}
