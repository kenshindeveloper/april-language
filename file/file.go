package file

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/kenshindeveloper/april/evaluator"
	"github.com/kenshindeveloper/april/lexer"
	"github.com/kenshindeveloper/april/object"
	"github.com/kenshindeveloper/april/parser"
	"github.com/kenshindeveloper/april/repl"
)

func Start(w io.Writer, r io.Reader, path string) {
	if len(path) <= len(".april") || path[len(path)-6:len(path)] != ".april" {
		log.Fatalf("incorrect path file.")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error to open file: '%s'", path)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParserProgram()
	if len(p.Error()) > 0 {
		repl.PrintParseError(w, p.Error())
		os.Exit(2)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		switch evaluated.(type) {
		case *object.Nil:
			io.WriteString(w, "")
		default:
			io.WriteString(w, evaluated.Inspect())
			io.WriteString(w, "\n")
		}
	}
}
