package parser_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/parser"
	"github.com/smacker/go-tree-sitter/c"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/stretchr/testify/assert"
)

func TestGo(t *testing.T) {
	path := "testdata/test.go"
	candidates := parser.NewParser(path, golang.GetLanguage()).Parse()

	assert.Equal(t, "A", candidates[0].Function.Name)
	assert.Equal(t, path, candidates[0].Path)

	assert.Equal(t, "B", candidates[1].Function.Name)
	assert.Equal(t, path, candidates[1].Path)
}

func TestJavaScript(t *testing.T) {
	path := "testdata/test.js"
	candidates := parser.NewParser(path, javascript.GetLanguage()).Parse()

	assert.Equal(t, "a", candidates[0].Function.Name)
	assert.Equal(t, path, candidates[0].Path)

	assert.Equal(t, "b", candidates[1].Function.Name)
	assert.Equal(t, path, candidates[1].Path)

	assert.Equal(t, "c", candidates[2].Function.Name)
	assert.Equal(t, path, candidates[2].Path)
}

func TestJava(t *testing.T) {
	path := "testdata/test.java"
	candidates := parser.NewParser(path, java.GetLanguage()).Parse()

	assert.Equal(t, "A", candidates[0].Function.Name)
	assert.Equal(t, path, candidates[0].Path)

	assert.Equal(t, "B", candidates[1].Function.Name)
	assert.Equal(t, path, candidates[1].Path)
}

func TestC(t *testing.T) {
	path := "testdata/test.c"
	candidates := parser.NewParser(path, c.GetLanguage()).Parse()

	assert.Equal(t, "main", candidates[0].Function.Name)
	assert.Equal(t, path, candidates[0].Path)

	assert.Equal(t, "a", candidates[1].Function.Name)
	assert.Equal(t, path, candidates[1].Path)
}
