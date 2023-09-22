package parser_test

import (
	"testing"

	"github.com/jochil/test-helper/pkg/parser"
)

func TestGraph(t *testing.T) {
	candidates := parser.NewParser(parser.GuessLanguage("testdata/cyclo/a.go")).Parse()
  candidates[0].SaveGraph()
}
