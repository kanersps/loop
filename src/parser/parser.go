package parser

import (
	"bufio"
	"fmt"
	"github.com/kanersps/loop/parser/tokenizer"
	"os"
	"strings"
)

var variables = map[string]int{}

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

		line = strings.Replace(line, "\t", "", -1)
		line = strings.Replace(line, " ", "", -1)

		if len(line) != 0 {
			tokenizer.Tokenize(scanner.Text(), lineNumber)

			lineNumber++
		}
	}
}
