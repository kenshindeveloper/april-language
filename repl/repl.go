package repl

import (
	"bufio"
	"fmt"
	"io"

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

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")

	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, APRIL_LOGO)
	io.WriteString(out, "Shit! exist errors:\n")
	for i, msg := range errors {
		io.WriteString(out, "\t"+fmt.Sprintf("%d", i)+"- "+msg+"\n\n")
	}
}
