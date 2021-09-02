package parser

import (
	"bufio"
	"fmt"
	"github.com/kanersps/loop/parser/lexer"
	"os"
)

var variables = map[string]int{}

func parse(value string, lineNumber int) *lexer.Lexer {
	return lexer.Create(value)
}

func ParseValue(value string) *lexer.Lexer {
	return parse(value, 1)
}

func ParseFile(path string) {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) != 0 {
			parse(scanner.Text(), lineNumber)

			lineNumber++
		}
	}
}
