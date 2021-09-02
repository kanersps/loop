package repl

import (
	"bufio"
	"fmt"
	"github.com/kanersps/loop/evaluator"
	"github.com/kanersps/loop/parser"
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
		parser := parser.Create(lexer)
		program := parser.ParseProgram()

		for _, err := range parser.Errors() {
			fmt.Println(err)
		}

		evaluated := evaluator.Eval(program)
		fmt.Println(evaluated)
		if evaluated != nil {
			io.WriteString(output, evaluated.Inspect())
			io.WriteString(output, "\n")
		}
	}
}
