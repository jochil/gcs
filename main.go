package main

import (
	"fmt"

	"github.com/jochil/test-helper/pkg/parser"

	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
)

func main() {
	fmt.Println("Go candidates", parser.Parse("examples/test.go", golang.GetLanguage()))
	fmt.Println("JavaScript candidates", parser.Parse("examples/test.js", javascript.GetLanguage()))
}
