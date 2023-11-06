package parser_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/candidate"
	"github.com/jochil/dlth/pkg/types"
)

func TestTypeScript(t *testing.T) {
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
