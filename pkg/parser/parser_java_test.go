package parser_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/candidate"
	"github.com/jochil/dlth/pkg/parser"
	"github.com/jochil/dlth/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestJava(t *testing.T) {
	path := "testdata/java/test.java"
	candidates := parser.NewParser(path, types.Java).Parse()

	class := "Foo"
	packageName := "org.example"

	simpleReturn := func(t string) []*candidate.Parameter {
		return []*candidate.Parameter{
			{Name: types.NoName, Type: t},
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
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "B",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("String"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "C",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPrivate,
		},
		{
			name: "D",
			params: []*candidate.Parameter{
				{Name: "d", Type: "int"},
				{Name: "e", Type: "String"},
			},
			returnValues: simpleReturn("String"),
			visibility:   types.VisibilityProtected,
		},
		{
			name:         "E",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("int"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "F",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("float"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "G",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("char"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "H",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("double"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "I",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("boolean"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "J",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("byte"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "K",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("long"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "L",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("long[]"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "M",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn("String[]"),
			visibility:   types.VisibilityPublic,
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
	candidates := parser.NewParser(path, types.Java).Parse()

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
