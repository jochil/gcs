package parser_test

import (
	"testing"

	"github.com/jochil/test-helper/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCycloGo(t *testing.T) {
	tests := map[string]struct {
		path   string
		wantCC int
	}{
		"no_control":        {path: "testdata/cyclo/a.go", wantCC: 1},
		"simple_if":         {path: "testdata/cyclo/b.go", wantCC: 2},
		"else_if":           {path: "testdata/cyclo/c.go", wantCC: 4},
		"switch_no_default": {path: "testdata/cyclo/d.go", wantCC: 2},
		"switch_default":    {path: "testdata/cyclo/e.go", wantCC: 3},
		"simple_for":        {path: "testdata/cyclo/f.go", wantCC: 2},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			candidates := parser.NewParser(parser.GuessLanguage(tc.path)).Parse()
			assert.Len(t, candidates, 1)
			cc, err := candidates[0].CyclomaticComplexity()
			require.NoError(t, err)
			assert.Equal(t, tc.wantCC, cc, "wrong cyclic complexity for function")
		})
	}
}
