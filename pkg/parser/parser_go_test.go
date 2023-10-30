package parser_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestGo_SimpleFunction(t *testing.T) {
	path := "testdata/golang/function.go"
	candidates := parser.NewParser(path, parser.Go).Parse()

	packageName := "examples"

	tests := []struct {
		name         string
		params       []*parser.Parameter
		returnValues []*parser.Parameter
		visibility   string
	}{
		{
			name: "A",
			params: []*parser.Parameter{
				{Name: "a", Type: "string"},
			},
			returnValues: []*parser.Parameter{
				{Name: parser.NoName, Type: "int8"},
			},
			visibility: parser.VisibilityPublic,
		},
		{
			name:   "B",
			params: []*parser.Parameter{},
			returnValues: []*parser.Parameter{
				{Name: "err", Type: "error"},
			},
			visibility: parser.VisibilityPublic,
		},
		{
			name:         "C",
			params:       []*parser.Parameter{},
			returnValues: []*parser.Parameter{},
			visibility:   parser.VisibilityPublic,
		},
		{
			name: "D",
			params: []*parser.Parameter{
				{Name: "d", Type: "string"},
			},
			returnValues: []*parser.Parameter{},
			visibility:   parser.VisibilityPublic,
		},
		{
			name:         "e",
			params:       []*parser.Parameter{},
			returnValues: []*parser.Parameter{},
			visibility:   parser.VisibilityPrivate,
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
	candidates := parser.NewParser(path, parser.Go).Parse()

	class := "*MyStruct"
	packageName := "examples"

	tests := []struct {
		name         string
		params       []*parser.Parameter
		returnValues []*parser.Parameter
	}{
		{
			name: "A",
			params: []*parser.Parameter{
				{Name: "a", Type: "int"},
				{Name: "b", Type: "uint"},
			},
			returnValues: []*parser.Parameter{
				{Name: "c", Type: "string"},
				{Name: "err", Type: "error"},
			},
		},
		{
			name: "B",
			params: []*parser.Parameter{
				{Name: "a", Type: "int"},
				{Name: "b", Type: "uint"},
			},
			returnValues: []*parser.Parameter{},
		},
		{
			name:   "C",
			params: []*parser.Parameter{},
			returnValues: []*parser.Parameter{
				{Name: "c", Type: "string"},
				{Name: "err", Type: "error"},
			},
		},
		{
			name:   "D",
			params: []*parser.Parameter{},
			returnValues: []*parser.Parameter{
				{Name: parser.NoName, Type: "error"},
			},
		},
		{
			name: "E",
			params: []*parser.Parameter{
				{Name: "a", Type: "int"},
			},
			returnValues: []*parser.Parameter{
				{Name: parser.NoName, Type: "string"},
				{Name: parser.NoName, Type: "error"},
			},
		},
		{
			name:         "F",
			params:       []*parser.Parameter{},
			returnValues: []*parser.Parameter{},
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

func testParams(t *testing.T, expected []*parser.Parameter, actual []*parser.Parameter) {
	t.Helper()
	assert.Len(t, actual, len(expected))
	for i, p := range actual {
		assert.Equal(t, p.Name, actual[i].Name)
		assert.Equal(t, p.Type, expected[i].Type)
	}
}
