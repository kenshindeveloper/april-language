package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/kenshindeveloper/april/repl"
)

const MAJOR int = 0
const MINOR int = 0
const MICRO int = 0

func main() {

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("April interpreter version: %d.%d.%d-alpha %s\n", MAJOR, MINOR, MICRO, user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
