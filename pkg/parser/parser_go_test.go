package parser_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/candidate"
	"github.com/jochil/dlth/pkg/parser"
	"github.com/jochil/dlth/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestGo_SimpleFunction(t *testing.T) {
	path := "testdata/golang/function.go"
	candidates := parser.NewParser(path, types.Go).Parse()

	packageName := "examples"

	tests := []struct {
		name         string
		params       []*candidate.Parameter
		returnValues []*candidate.Parameter
		visibility   string
	}{
		{
			name: "A",
			params: []*candidate.Parameter{
				{Name: "a", Type: "string"},
			},
			returnValues: []*candidate.Parameter{
				{Name: types.NoName, Type: "int8"},
			},
			visibility: types.VisibilityPublic,
		},
		{
			name:   "B",
			params: []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{
				{Name: "err", Type: "error"},
			},
			visibility: types.VisibilityPublic,
		},
		{
			name:         "C",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
		{
			name: "D",
			params: []*candidate.Parameter{
				{Name: "d", Type: "string"},
			},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "e",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPrivate,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			method := candidates[i]
			assert.Equal(t, tc.name, method.Function.Name)
			assert.Equal(t, packageName, method.Package)
			assert.Equal(t, tc.visibility, method.Function.Visibility)

			assert.Len(t, method.Function.Parameters, len(tc.params))
			for i, p := range tc.params {
				assert.Equal(t, p.Name, method.Function.Parameters[i].Name)
				assert.Equal(t, p.Type, method.Function.Parameters[i].Type)
			}

			assert.Len(t, method.Function.ReturnValues, len(tc.returnValues))
			for i, p := range tc.returnValues {
				assert.Equal(t, p.Name, method.Function.ReturnValues[i].Name)
				assert.Equal(t, p.Type, method.Function.ReturnValues[i].Type)
			}
		})
	}
}

func TestGo_Method(t *testing.T) {
	path := "testdata/golang/method.go"
	candidates := parser.NewParser(path, types.Go).Parse()

	class := "*MyStruct"
	packageName := "examples"

	tests := []struct {
		name         string
		params       []*candidate.Parameter
		returnValues []*candidate.Parameter
	}{
		{
			name: "A",
			params: []*candidate.Parameter{
				{Name: "a", Type: "int"},
				{Name: "b", Type: "uint"},
			},
			returnValues: []*candidate.Parameter{
				{Name: "c", Type: "string"},
				{Name: "err", Type: "error"},
			},
		},
		{
			name: "B",
			params: []*candidate.Parameter{
				{Name: "a", Type: "int"},
				{Name: "b", Type: "uint"},
			},
			returnValues: []*candidate.Parameter{},
		},
		{
			name:   "C",
			params: []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{
				{Name: "c", Type: "string"},
				{Name: "err", Type: "error"},
			},
		},
		{
			name:   "D",
			params: []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{
				{Name: types.NoName, Type: "error"},
			},
		},
		{
			name: "E",
			params: []*candidate.Parameter{
				{Name: "a", Type: "int"},
			},
			returnValues: []*candidate.Parameter{
				{Name: types.NoName, Type: "string"},
				{Name: types.NoName, Type: "error"},
			},
		},
		{
			name:         "F",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			method := candidates[i]
			assert.Equal(t, tc.name, method.Function.Name)
			assert.Equal(t, class, method.Class)
			assert.Equal(t, packageName, method.Package)

			testParams(t, tc.params, method.Function.Parameters)
			testParams(t, tc.returnValues, method.Function.ReturnValues)
		})
	}
}

func testParams(t *testing.T, expected []*candidate.Parameter, actual []*candidate.Parameter) {
	t.Helper()
	assert.Len(t, actual, len(expected))
	for i, p := range actual {
		assert.Equal(t, p.Name, actual[i].Name)
		assert.Equal(t, p.Type, expected[i].Type)
	}
}
