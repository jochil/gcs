package parser_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetrics(t *testing.T) {
	tests := map[string]struct {
		path  string
		cc    int
		lines int
	}{
		"no_control":        {path: "testdata/cyclo/a.go", cc: 1, lines: 4},
		"simple_if":         {path: "testdata/cyclo/b.go", cc: 2, lines: 6},
		"else_if":           {path: "testdata/cyclo/c.go", cc: 4, lines: 12},
		"switch_no_default": {path: "testdata/cyclo/d.go", cc: 2, lines: 6},
		"switch_default":    {path: "testdata/cyclo/e.go", cc: 3, lines: 10},
		"simple_for":        {path: "testdata/cyclo/f.go", cc: 2, lines: 5},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			candidates := parser.NewParser(parser.GuessLanguage(tc.path)).Parse()
			require.Len(t, candidates, 1)
			c := candidates[0]
			cc, err := c.CyclomaticComplexity()
			require.NoError(t, err)
			assert.Equal(t, tc.cc, cc, "wrong cyclic complexity for function")
			assert.Equal(t, tc.lines, c.Lines)
		})
	}
}
