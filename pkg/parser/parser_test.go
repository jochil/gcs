package parser_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/candidate"
	"github.com/jochil/dlth/pkg/parser"
	"github.com/jochil/dlth/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestJavaScript(t *testing.T) {
	path := "testdata/javascript/declaration.js"
	candidates := parser.NewParser(path, types.JavaScript).Parse()

	assert.Equal(t, "a", candidates[0].Function.Name)
	assert.Equal(t, path, candidates[0].Path)

	assert.Equal(t, "b", candidates[1].Function.Name)
	assert.Equal(t, path, candidates[1].Path)

	assert.Equal(t, "c", candidates[2].Function.Name)
	assert.Equal(t, path, candidates[2].Path)
}

func TestC(t *testing.T) {
	path := "testdata/test.c"
	candidates := parser.NewParser(path, types.C).Parse()

	assert.Equal(t, "main", candidates[0].Function.Name)
	assert.Equal(t, path, candidates[0].Path)

	assert.Equal(t, "a", candidates[1].Function.Name)
	assert.Equal(t, path, candidates[1].Path)
}

func simpleReturn(t *testing.T, typeName string) []*candidate.Parameter {
	t.Helper()
	return []*candidate.Parameter{
		{Name: types.NoName, Type: typeName},
	}
}

type candidateTestCase struct {
	name         string
	params       []*candidate.Parameter
	returnValues []*candidate.Parameter
	visibility   string
	class        string
	packageName  string
}

func runParserTests(t *testing.T, tests []candidateTestCase, path string, language types.Language) {
	candidates := parser.NewParser(path, language).Parse()
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assertCandidate(t, tc, candidates[i])
		})
	}
}

func assertParams(t *testing.T, expected []*candidate.Parameter, actual []*candidate.Parameter) {
	t.Helper()
	assert.Len(t, actual, len(expected))
	for i, p := range actual {
		assert.Equal(t, expected[i].Name, p.Name, "invalid parameter name")
		assert.Equal(t, expected[i].Type, p.Type, "invalid parameter type")
	}
}

func assertCandidate(t *testing.T, tc candidateTestCase, c *candidate.Candidate) {
	assert.Equal(t, tc.name, c.Function.Name, "invalid function name")
	assert.Equal(t, tc.class, c.Class, "invalid class")
	assert.Equal(t, tc.packageName, c.Package, "invalid package")
	assert.Equal(t, tc.visibility, c.Function.Visibility, "invalid visibility")

	assertParams(t, tc.params, c.Function.Parameters)
	assertParams(t, tc.returnValues, c.Function.ReturnValues)
}
