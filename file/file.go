package file

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/kenshindeveloper/april/lexer"
	"github.com/kenshindeveloper/april/token"
)

func Start(name string) bool {

	file, err := os.Open(name)
	if err != nil {
		log.Fatalf("file: %s not exist\n", name)
		return false
	}

	scanner := bufio.NewScanner(file)

	buffersize := 1
	words := make([]string, buffersize)
	index := 0
	existerr := false
	text := ""

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			existerr = true
			break
		}

		text = scanner.Text()
		if len(text) > 0 {
			words[index] = text
			index++
		}

		if index >= len(words) {
			auxwords := make([]string, buffersize)
			words = append(words, auxwords...)
		}
	}

	if existerr {
		log.Fatal("error reading the file")
		return false
	}

	for _, line := range words {
		fmt.Printf("%s\n", line)
	}

	if len(words) > 0 {
		fmt.Printf("-----------------------------------\n")
		for _, line := range words {
			l := lexer.New(line)
			for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
				fmt.Printf("%+v\n", tok)
			}
		}
		fmt.Printf("-----------------------------------\n")
	}

	return true
}
