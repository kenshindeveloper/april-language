package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/kenshindeveloper/april/file"
	"github.com/kenshindeveloper/april/repl"
)

const major int = 0
const minor int = 0
const micro int = 0

func main() {

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	switch len(os.Args) {
	case 1:
		//REPL (bucle lectura-evaluacion-impresion)
		fmt.Printf("April interpreter version: %d.%d.%d-alpha %s\n", major, minor, micro, user.Username)
		repl.Start(os.Stdin, os.Stdout)
	case 2:
		fmt.Printf("April interpreter version: %d.%d.%d-alpha %s\n\n", major, minor, micro, user.Username)
		file.Start(os.Args[1])
	default:
		fmt.Println("incorrect number of parametres")
	}
}
