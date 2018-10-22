package repl

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"

	"github.com/kenshindeveloper/april/evaluator"
	"github.com/kenshindeveloper/april/lexer"
	"github.com/kenshindeveloper/april/libs"
	"github.com/kenshindeveloper/april/object"
	"github.com/kenshindeveloper/april/parser"
)

const PROMPT = ">> "
const PROMPTBLOCK = "... "

const APRIL_ERROR = `
               _ _ _     
             /        \
      (\/)  <  ERROR!  |
      ( . .) \ _ _  _ /
 __\/_C(")(")__\/___\//____
`

//Start es una funcion que recibe dos parametros, de escritura y lectura. Ejecuta el programa en modo interactivo
func Start(w io.Writer, r io.Reader, version string) {
	info, err := user.Current()
	if err != nil {
		log.Fatal("problema al leer la informacion del usuario.")
	}
	cleanConsole()
	fmt.Printf("Hello %s! This is the April programming language version %s!\n", info.Name, version)
	fmt.Printf("Feel free to type in commands.\n")
	env := object.NewEnvironment()
	stackBlock := libs.NewStack()
	inputBlock := ""

	for {
		if stackBlock.Len() == 0 {
			fmt.Printf("%s", PROMPT)
		} else {
			fmt.Printf("%s", PROMPTBLOCK)
		}
		reader := bufio.NewScanner(r)
		if !reader.Scan() {
			return
		}

		input := reader.Text()
		if ok, _ := regexp.MatchString("^(.*)?\\{( )*$", input); ok {
			stackBlock.Push("{")
		} else if ok, _ := regexp.MatchString("^(.*)?\\}( |.)*$", input); ok {
			stackBlock.Pop()
		}

		switch stackBlock.Len() {
		case 0:
			if inputBlock != "" {
				input = inputBlock + input
				inputBlock = ""
			}

			l := lexer.New(input)
			p := parser.New(l)
			program := p.ParserProgram()
			if len(p.Error()) > 0 {
				PrintParseError(w, p.Error())
				continue
			}

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
		default:
			if evalTextInput(input) {
				input += ";"
			}

			inputBlock += input
		}
	}
}

func evalTextInput(text string) bool {
	if ok, _ := regexp.MatchString("^(.*)?(\\{|\\}|;)( )*$", text); ok {
		return false
	}
	return true
}

func PrintParseError(w io.Writer, errors []string) {
	io.WriteString(w, APRIL_ERROR+"\n\n")
	if len(errors) > 1 {
		io.WriteString(w, "shit! there are error.\n")
	} else {
		io.WriteString(w, "shit! there is a error.\n")
	}
	io.WriteString(w, "parser errors:\n")

	for _, err := range errors {
		io.WriteString(w, "\t- "+err+"\n")
	}
}

func cleanConsole() {
	command := exec.Command("cmd", "/c", "cls")
	command.Stdout = os.Stdout
	command.Run()
}
