package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kenshindeveloper/april/evaluator"
	"github.com/kenshindeveloper/april/lexer"
	"github.com/kenshindeveloper/april/parser"
)

const PROMPT string = ">> "

const APRIL_LOGO = `
               _ _ _     
             /        \
      (\/)  <  ERROR!  
      ( . .) \ _ _  _ /
 __\/_C(")(")__\/___\//____

`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParserProgram()
		if len(p.Error()) != 0 {
			printParserErrors(out, p.Error())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, APRIL_LOGO)
	if len(errors) > 1 {
		io.WriteString(out, "Shit! there are errors:\n")
	} else if len(errors) == 1 {
		io.WriteString(out, "Shit! there is a error:\n")
	}
	for _, msg := range errors {
		io.WriteString(out, "\t- "+msg+"\n\n")
	}
}
