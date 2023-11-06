package parser_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/candidate"
	"github.com/jochil/dlth/pkg/parser"
	"github.com/jochil/dlth/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestTypeScript(t *testing.T) {
	path := "testdata/typescript/function.ts"
	candidates := parser.NewParser(path, types.TypeScript).Parse()

	tests := []struct {
		name         string
		params       []*candidate.Parameter
		returnValues []*candidate.Parameter
		visibility   string
	}{
		{
			name: "A",
			params: []*candidate.Parameter{
				{Name: "a", Type: "number"},
				{Name: "b", Type: "string"},
			},
			returnValues: []*candidate.Parameter{
				{Name: types.NoName, Type: "number"},
			},
			visibility: types.VisibilityPublic,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			method := candidates[i]
			assert.Equal(t, tc.name, method.Function.Name)
			assert.Equal(t, tc.visibility, method.Function.Visibility)

			testParams(t, tc.params, method.Function.Parameters)
			testParams(t, tc.returnValues, method.Function.ReturnValues)
		})
	}
}
