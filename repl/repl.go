package repl

import (
	"bufio"
	"fmt"
	"github.com/kanersps/loop/evaluator"
	"github.com/kanersps/loop/object"
	"github.com/kanersps/loop/parser"
	"github.com/kanersps/loop/parser/lexer"
	"io"
)

func Console(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	env := object.NewEnvironment()

	for {
		fmt.Printf(">> ")
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.Create(line)
		parser := parser.Create(l)
		program := parser.ParseProgram()

		if len(parser.Errors()) != 0 {
			printParserErrors(output, parser.Errors())
			continue
		}
		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			io.WriteString(output, evaluated.Inspect())
			io.WriteString(output, "\n")
		}

		//program.PrintAST()
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
