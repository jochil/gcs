package parser_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestJavaScript(t *testing.T) {
	path := "testdata/test.js"
	candidates := parser.NewParser(path, parser.JavaScript).Parse()

	assert.Equal(t, "a", candidates[0].Function.Name)
	assert.Equal(t, path, candidates[0].Path)

	assert.Equal(t, "b", candidates[1].Function.Name)
	assert.Equal(t, path, candidates[1].Path)

	assert.Equal(t, "c", candidates[2].Function.Name)
	assert.Equal(t, path, candidates[2].Path)
}

func TestC(t *testing.T) {
	path := "testdata/test.c"
	candidates := parser.NewParser(path, parser.C).Parse()

	assert.Equal(t, "main", candidates[0].Function.Name)
	assert.Equal(t, path, candidates[0].Path)

	assert.Equal(t, "a", candidates[1].Function.Name)
	assert.Equal(t, path, candidates[1].Path)
}
