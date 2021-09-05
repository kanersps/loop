package cmd

import (
	"flag"
	"fmt"
	"github.com/kanersps/loop/evaluator"
	"github.com/kanersps/loop/object"
	"github.com/kanersps/loop/parser"
	"github.com/kanersps/loop/parser/lexer"
	"github.com/kanersps/loop/repl"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func Execute() {
	executeFile := flag.String("file", "none-provided", "The file you want to interpret")

	flag.Parse()

	if *executeFile == "none-provided" {
		repl.Console(os.Stdin, os.Stdout)
	} else {
		fmt.Println(*executeFile)

		env := object.NewEnvironment()

		input, err := ioutil.ReadFile(*executeFile)

		if err != nil {
			log.Fatal(err)
		}

		l := lexer.Create(string(input))
		p := parser.Create(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(os.Stdout, p.Errors())
			log.Fatal()
		}

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			io.WriteString(os.Stdout, evaluated.Inspect())
			io.WriteString(os.Stdout, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
