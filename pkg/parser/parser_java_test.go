package parser_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/candidate"
	"github.com/jochil/dlth/pkg/parser"
	"github.com/jochil/dlth/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestJava_Methods(t *testing.T) {
	tests := []candidateTestCase{
		{
			name:        "A",
			packageName: "org.example",
			class:       "Foo",
			params: []*candidate.Parameter{
				{Name: "a", Type: "String"},
			},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
			static:       true,
		},
		{
			name:         "B",
			packageName:  "org.example",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn(t, "String"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "C",
			packageName:  "org.example",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPrivate,
		},
		{
			name:        "D",
			packageName: "org.example",
			class:       "Foo",
			params: []*candidate.Parameter{
				{Name: "d", Type: "int"},
				{Name: "e", Type: "String"},
			},
			returnValues: simpleReturn(t, "String"),
			visibility:   types.VisibilityProtected,
		},
		{
			name:         "E",
			packageName:  "org.example",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn(t, "int"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "F",
			packageName:  "org.example",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn(t, "float"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "G",
			packageName:  "org.example",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn(t, "char"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "H",
			packageName:  "org.example",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn(t, "double"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "I",
			packageName:  "org.example",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn(t, "boolean"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "J",
			packageName:  "org.example",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn(t, "byte"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "K",
			packageName:  "org.example",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn(t, "long"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "L",
			packageName:  "org.example",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn(t, "long[]"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "M",
			packageName:  "org.example",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn(t, "String[]"),
			visibility:   types.VisibilityPublic,
		},
	}
	runParserTests(t, tests, "testdata/java/method.java", types.Java)
}

func TestJava_Overloading(t *testing.T) {
	tests := []candidateTestCase{
		{
			name:         "A",
			packageName:  "org.example",
			class:        "Foo",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
		{
			name:        "A",
			packageName: "org.example",
			class:       "Foo",
			params: []*candidate.Parameter{
				{Name: "a", Type: "String"},
			},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
		{
			name:        "A",
			packageName: "org.example",
			class:       "Foo",
			params: []*candidate.Parameter{
				{Name: "a", Type: "int"},
			},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
	}
	runParserTests(t, tests, "testdata/java/overloading.java", types.Java)
}

func TestJava_Parameter(t *testing.T) {
	tests := []candidateTestCase{
		{
			name:        "Spread",
			packageName: "org.example",
			class:       "Foo",
			params: []*candidate.Parameter{
				{Name: "a", Type: "String..."},
			},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
		{
			name:        "Generic",
			packageName: "org.example",
			class:       "Foo",
			params: []*candidate.Parameter{
				{Name: "a", Type: "Map<String, Integer>"},
			},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
	}

	runParserTests(t, tests, "testdata/java/parameter.java", types.Java)
}

func TestJava_Constructor(t *testing.T) {
	candidates := parser.NewParser("testdata/java/constructor.java", types.Java).Parse()
	require.Len(t, candidates, 1)
	require.Len(t, candidates[0].Class.Constructors, 2)
}
