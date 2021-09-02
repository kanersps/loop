package main

import (
	"flag"
	"fmt"
	"github.com/kanersps/loop/parser"
)

func main() {
	executeFile := flag.String("file", "", "The file you want to interpret")

	flag.Parse()

	fmt.Println(*executeFile)

	parser.ParseFile(*executeFile)
}
