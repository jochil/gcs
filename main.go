package main

import (
	"github.com/jochil/test-helper/pkg/generator"
	"github.com/jochil/test-helper/pkg/parser"

	"github.com/smacker/go-tree-sitter/c"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
)

func main() {
  candidates := parser.NewParser("examples/test_cyclo.go", golang.GetLanguage()).Parse()
  generator.CreateGoTest(candidates[0])
}

func runExamples() {
	parser.NewParser("examples/test.go", golang.GetLanguage()).Parse()
	parser.NewParser("examples/test.js", javascript.GetLanguage()).Parse()
	parser.NewParser("examples/test.java", java.GetLanguage()).Parse()
	parser.NewParser("examples/test.c", c.GetLanguage()).Parse()
}
