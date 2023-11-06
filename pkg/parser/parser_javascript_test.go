package parser_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/candidate"
	"github.com/jochil/dlth/pkg/types"
)

func TestJavaScript_Declaration(t *testing.T) {
	tests := []candidateTestCase{
		{
			name:         "a",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "b",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "c",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
	}
	runParserTests(t, tests, "testdata/javascript/declaration.js", types.JavaScript)
}

func TestJavaScript_Method(t *testing.T) {
	tests := []candidateTestCase{
		{
			name:         "A",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "B",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
			// TODO actually this should be true, but it is not showing up in the AST
			static: false,
		},
		{
			name:         "C",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			// TODO actually this should be private, but # is not handled in the AST
			visibility: types.VisibilityPublic,
		},
	}
	runParserTests(t, tests, "testdata/javascript/method.js", types.JavaScript)
}
