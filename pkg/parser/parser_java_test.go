package parser_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/parser"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/stretchr/testify/assert"
)

func TestJava(t *testing.T) {
	path := "testdata/java/test.java"
	candidates := parser.NewParser(path, java.GetLanguage()).Parse()

	class := "Foo"
	packageName := "org.example"

	tests := []struct {
		name         string
		params       []*parser.Parameter
		returnValues []*parser.Parameter
	}{
		{
			name: "A",
			params: []*parser.Parameter{
				{Name: "a", Type: "String"},
			},
			returnValues: []*parser.Parameter{},
		},
		{
			name:   "B",
			params: []*parser.Parameter{},
			returnValues: []*parser.Parameter{
				{Name: parser.NoName, Type: "String"},
			},
		},
		{
			name:         "C",
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
