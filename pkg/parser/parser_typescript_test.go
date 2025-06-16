package parser_test

import (
	"testing"

	"github.com/jochil/gcs/pkg/candidate"
	"github.com/jochil/gcs/pkg/types"
)

func TestTypeScript_Functions(t *testing.T) {
	tests := []candidateTestCase{
		{
			name: "A",
			params: []*candidate.Parameter{
				{Name: "a", Type: "number"},
				{Name: "b", Type: "string"},
			},
			returnValues: simpleReturn(t, "number"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "B",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
		{
			name: "C",
			params: []*candidate.Parameter{
				{Name: "a", Type: "object | null"},
			},
			returnValues: simpleReturn(t, "string"),
			visibility:   types.VisibilityPublic,
		},
	}

	runParserTests(t, tests, "testdata/typescript/function.ts", types.TypeScript)
}

func TestTypeScript_Methods(t *testing.T) {
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
			visibility:   types.VisibilityPrivate,
		},
		{
			name:         "C",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPrivate,
		},
		{
			name:         "D",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityProtected,
		},
	}

	runParserTests(t, tests, "testdata/typescript/method.ts", types.TypeScript)
}
