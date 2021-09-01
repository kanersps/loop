package main

import (
	"flag"
	"fmt"
	"loop/src/parser"
)

func main() {
	executeFile := flag.String("file", "", "The file you want to interpret")

	flag.Parse()

	fmt.Println(*executeFile)

	parser.ParseFile(*executeFile)
}
