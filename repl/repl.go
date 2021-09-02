package repl

import (
	"bufio"
	"fmt"
	"github.com/kanersps/loop/parser"
	"github.com/kanersps/loop/parser/tokens"
	"io"
)

func Console(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)

	for {
		fmt.Printf(">> ")
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		lexer := parser.ParseValue(line)

		for tok := lexer.FindToken(); tok.TokenType != tokens.EOF; tok = lexer.FindToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
