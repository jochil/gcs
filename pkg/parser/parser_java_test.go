package parser_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/candidate"
	"github.com/jochil/dlth/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestJava(t *testing.T) {
	path := "testdata/java/test.java"
	candidates := parser.NewParser(path, parser.Java).Parse()

	class := "Foo"
	packageName := "org.example"

	simpleReturn := func(t string) []*candidate.Parameter {
		return []*candidate.Parameter{
			{Name: parser.NoName, Type: t},
		}
	}

	tests := []struct {
		name         string
		params       []*candidate.Parameter
		returnValues []*candidate.Parameter
		visibility   string
	}{
		{
			name: "A",
			params: []*candidate.Parameter{
				{Name: "a", Type: "String"},
			},
			returnValues: []*candidate.Parameter{},
			visibility:   parser.VisibilityPublic,
		},
		{
			name:         "B",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("String"),
			visibility:   parser.VisibilityPublic,
		},
		{
			name:         "C",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   parser.VisibilityPrivate,
		},
		{
			name: "D",
			params: []*candidate.Parameter{
				{Name: "d", Type: "int"},
				{Name: "e", Type: "String"},
			},
			returnValues: simpleReturn("String"),
			visibility:   parser.VisibilityProtected,
		},
		{
			name:         "E",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("int"),
			visibility:   parser.VisibilityPublic,
		},
		{
			name:         "F",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("float"),
			visibility:   parser.VisibilityPublic,
		},
		{
			name:         "G",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("char"),
			visibility:   parser.VisibilityPublic,
		},
		{
			name:         "H",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("double"),
			visibility:   parser.VisibilityPublic,
		},
		{
			name:         "I",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("boolean"),
			visibility:   parser.VisibilityPublic,
		},
		{
			name:         "J",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("byte"),
			visibility:   parser.VisibilityPublic,
		},
		{
			name:         "K",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("long"),
			visibility:   parser.VisibilityPublic,
		},
		{
			name:         "L",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("long[]"),
			visibility:   parser.VisibilityPublic,
		},
		{
			name:         "M",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("String[]"),
			visibility:   parser.VisibilityPublic,
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
		params       []*candidate.Parameter
		returnValues []*candidate.Parameter
	}{
		{
			name:         "A",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
		},
		{
			name: "A",
			params: []*candidate.Parameter{
				{Name: "a", Type: "String"},
			},
			returnValues: []*candidate.Parameter{},
		},
		{
			name: "A",
			params: []*candidate.Parameter{
				{Name: "a", Type: "int"},
			},
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
