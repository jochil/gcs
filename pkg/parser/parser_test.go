package parser_test

import (
	"testing"

	"github.com/jochil/gcs/pkg/candidate"
	"github.com/jochil/gcs/pkg/parser"
	"github.com/jochil/gcs/pkg/types"
	"github.com/stretchr/testify/assert"
)

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
	static       bool
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
	assert.Len(t, actual, len(expected), "invalid parameter amount")
	for i, p := range actual {
		assert.Equal(t, expected[i].Name, p.Name, "invalid parameter name")
		assert.Equal(t, expected[i].Type, p.Type, "invalid parameter type")
	}
}

func assertCandidate(t *testing.T, tc candidateTestCase, c *candidate.Candidate) {
	assert.Equal(t, tc.name, c.Function.Name, "invalid function name")
	if tc.class == "" {
		assert.Nil(t, c.Class, "invalid class")
	} else {
		assert.Equal(t, tc.class, c.Class.Name, "invalid class")
	}
	assert.Equal(t, tc.packageName, c.Package, "invalid package")
	assert.Equal(t, tc.visibility, c.Function.Visibility, "invalid visibility")
	assert.Equal(t, tc.static, c.Function.Static, "invalid static modifier")

	assertParams(t, tc.params, c.Function.Parameters)
	assertParams(t, tc.returnValues, c.Function.ReturnValues)
}
