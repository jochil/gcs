package parser_test

import (
	"testing"

	"github.com/jochil/gcs/pkg/candidate"
	"github.com/jochil/gcs/pkg/types"
)

func TestGo_SimpleFunction(t *testing.T) {
	tests := []candidateTestCase{
		{
			name:        "A",
			packageName: "examples",
			params: []*candidate.Parameter{
				{Name: "a", Type: "string"},
			},
			returnValues: simpleReturn(t, "int8"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:        "B",
			packageName: "examples",
			params:      []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{
				{Name: "err", Type: "error"},
			},
			visibility: types.VisibilityPublic,
		},
		{
			name:         "C",
			packageName:  "examples",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
		{
			name:        "D",
			packageName: "examples",
			params: []*candidate.Parameter{
				{Name: "d", Type: "string"},
			},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "e",
			packageName:  "examples",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPrivate,
		},
	}

	runParserTests(t, tests, "testdata/golang/function.go", types.Go)
}

func TestGo_Method(t *testing.T) {
	tests := []candidateTestCase{
		{
			name:        "A",
			class:       "MyStruct",
			packageName: "examples",
			params: []*candidate.Parameter{
				{Name: "a", Type: "int"},
				{Name: "b", Type: "uint"},
			},
			returnValues: []*candidate.Parameter{
				{Name: "c", Type: "string"},
				{Name: "err", Type: "error"},
			},
			visibility: types.VisibilityPublic,
		},
		{
			name:        "B",
			class:       "MyStruct",
			packageName: "examples",
			params: []*candidate.Parameter{
				{Name: "a", Type: "int"},
				{Name: "b", Type: "uint"},
			},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
		{
			name:        "C",
			class:       "MyStruct",
			packageName: "examples",
			params:      []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{
				{Name: "c", Type: "string"},
				{Name: "err", Type: "error"},
			},
			visibility: types.VisibilityPublic,
		},
		{
			name:         "D",
			class:        "MyStruct",
			packageName:  "examples",
			params:       []*candidate.Parameter{},
			returnValues: simpleReturn(t, "error"),
			visibility:   types.VisibilityPublic,
		},
		{
			name:        "E",
			class:       "MyStruct",
			packageName: "examples",
			params: []*candidate.Parameter{
				{Name: "a", Type: "int"},
			},
			returnValues: []*candidate.Parameter{
				{Name: types.NoName, Type: "string"},
				{Name: types.NoName, Type: "error"},
			},
			visibility: types.VisibilityPublic,
		},
		{
			name:         "F",
			class:        "MyStruct",
			packageName:  "examples",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "G",
			class:        "MyStruct",
			packageName:  "examples",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
	}

	runParserTests(t, tests, "testdata/golang/method.go", types.Go)
}
