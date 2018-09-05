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
	// scanner.Split(bufio.ScanWords)

	buffersize := 10
	words := make([]string, buffersize)
	index := 0
	existerr := false

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			existerr = true
			break
		}

		words[index] = scanner.Text()
		index++

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
		fmt.Println(line)
	}

	if len(words) > 0 {
		l := lexer.New(words[0])
		fmt.Printf("-----------------------------------\n")
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}

		fmt.Printf("-----------------------------------\n")
	}

	return true
}
