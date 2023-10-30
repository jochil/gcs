package parser_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestJava(t *testing.T) {
	path := "testdata/java/test.java"
	candidates := parser.NewParser(path, parser.Java).Parse()

	class := "Foo"
	packageName := "org.example"

	tests := []struct {
		name         string
		params       []*parser.Parameter
		returnValues []*parser.Parameter
		visibility   string
	}{
		{
			name: "A",
			params: []*parser.Parameter{
				{Name: "a", Type: "String"},
			},
			returnValues: []*parser.Parameter{},
			visibility:   parser.VisibilityPublic,
		},
		{
			name:   "B",
			params: []*parser.Parameter{},
			returnValues: []*parser.Parameter{
				{Name: parser.NoName, Type: "String"},
			},
			visibility: parser.VisibilityPublic,
		},
		{
			name:         "C",
			params:       []*parser.Parameter{},
			returnValues: []*parser.Parameter{},
			visibility:   parser.VisibilityPrivate,
		},
		{
			name: "D",
			params: []*parser.Parameter{
				{Name: "d", Type: "int"},
				{Name: "e", Type: "String"},
			},
			returnValues: []*parser.Parameter{
				{Name: parser.NoName, Type: "String"},
			},
			visibility: parser.VisibilityProtected,
		},
		{
			name:   "E",
			params: []*parser.Parameter{},
			returnValues: []*parser.Parameter{
				{Name: parser.NoName, Type: "int"},
			},
			visibility: parser.VisibilityPublic,
		},
		{
			name:   "F",
			params: []*parser.Parameter{},
			returnValues: []*parser.Parameter{
				{Name: parser.NoName, Type: "float"},
			},
			visibility: parser.VisibilityPublic,
		},
		{
			name:   "G",
			params: []*parser.Parameter{},
			returnValues: []*parser.Parameter{
				{Name: parser.NoName, Type: "char"},
			},
			visibility: parser.VisibilityPublic,
		},
		{
			name:   "H",
			params: []*parser.Parameter{},
			returnValues: []*parser.Parameter{
				{Name: parser.NoName, Type: "double"},
			},
			visibility: parser.VisibilityPublic,
		},
		{
			name:   "I",
			params: []*parser.Parameter{},
			returnValues: []*parser.Parameter{
				{Name: parser.NoName, Type: "boolean"},
			},
			visibility: parser.VisibilityPublic,
		},
		{
			name:   "J",
			params: []*parser.Parameter{},
			returnValues: []*parser.Parameter{
				{Name: parser.NoName, Type: "byte"},
			},
			visibility: parser.VisibilityPublic,
		},
		{
			name:   "K",
			params: []*parser.Parameter{},
			returnValues: []*parser.Parameter{
				{Name: parser.NoName, Type: "long"},
			},
			visibility: parser.VisibilityPublic,
		},
		{
			name:   "L",
			params: []*parser.Parameter{},
			returnValues: []*parser.Parameter{
				{Name: parser.NoName, Type: "long[]"},
			},
			visibility: parser.VisibilityPublic,
		},
		{
			name:   "M",
			params: []*parser.Parameter{},
			returnValues: []*parser.Parameter{
				{Name: parser.NoName, Type: "String[]"},
			},
			visibility: parser.VisibilityPublic,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			method := candidates[i]
			assert.Equal(t, tc.name, method.Function.Name)
			assert.Equal(t, class, method.Class)
			assert.Equal(t, packageName, method.Package)
			assert.Equal(t, tc.visibility, method.Function.Visibility)

			testParams(t, tc.params, method.Function.Parameters)
			testParams(t, tc.returnValues, method.Function.ReturnValues)
		})
	}
}

func TestJava_Overloading(t *testing.T) {
	path := "testdata/java/overloading.java"
	candidates := parser.NewParser(path, parser.Java).Parse()

	class := "Foo"
	packageName := "org.example"

	tests := []struct {
		name         string
		params       []*parser.Parameter
		returnValues []*parser.Parameter
	}{
		{
			name:         "A",
			params:       []*parser.Parameter{},
			returnValues: []*parser.Parameter{},
		},
		{
			name: "A",
			params: []*parser.Parameter{
				{Name: "a", Type: "String"},
			},
			returnValues: []*parser.Parameter{},
		},
		{
			name: "A",
			params: []*parser.Parameter{
				{Name: "a", Type: "int"},
			},
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
