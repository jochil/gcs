package main

import (
	"github.com/jochil/dlth/pkg/generator"
	"github.com/jochil/dlth/pkg/parser"
)

func main() {
	parser.NewParser(parser.GuessLanguage("examples/test.js")).Parse()
	parser.NewParser(parser.GuessLanguage("examples/test.java")).Parse()
	parser.NewParser(parser.GuessLanguage("examples/test.c")).Parse()

	candidates := parser.NewParser(parser.GuessLanguage("examples/test.go")).Parse()
	generator.CreateGoTest(candidates[0])
	//candidates[0].SaveGraph()
}
