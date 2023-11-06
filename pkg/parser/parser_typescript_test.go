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
	}

	runParserTests(t, tests, "testdata/typescript/function.ts", types.TypeScript)
}
