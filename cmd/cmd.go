package cmd

import (
	"flag"
	"fmt"
	"github.com/kanersps/loop/parser"
	"github.com/kanersps/loop/repl"
	"os"
)

func Execute() {
	executeFile := flag.String("file", "none-provided", "The file you want to interpret")

	flag.Parse()

	if *executeFile == "none-provided" {
		repl.Console(os.Stdin, os.Stdout)
	} else {
		fmt.Println(*executeFile)

		parser.ParseFile(*executeFile)
	}
}
