package main

import (
	"fmt"
	"os"

	"github.com/kenshindeveloper/april/file"
	"github.com/kenshindeveloper/april/repl"
)

const mayor = 1
const minus = 0
const micro = 0

func main() {
	switch len(os.Args) {
	case 1:
		repl.Start(os.Stdout, os.Stdin, fmt.Sprintf("%d.%d.%d", mayor, minus, micro))
	case 2:
		file.Start(os.Stdout, os.Stdin, os.Args[1])
	default:
		fmt.Printf("number parameters is incorrect.")
	}
}
