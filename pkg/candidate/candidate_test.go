package candidate_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/candidate"
	"github.com/jochil/dlth/pkg/helper"
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
		"no_control":        {path: "../cfg/testdata/cyclo/golang/a.go", cc: 1, lines: 4},
		"simple_if":         {path: "../cfg/testdata/cyclo/golang/b.go", cc: 2, lines: 6},
		"else_if":           {path: "../cfg/testdata/cyclo/golang/c.go", cc: 4, lines: 12},
		"switch_no_default": {path: "../cfg/testdata/cyclo/golang/d.go", cc: 2, lines: 6},
		"switch_default":    {path: "../cfg/testdata/cyclo/golang/e.go", cc: 3, lines: 10},
		"simple_for":        {path: "../cfg/testdata/cyclo/golang/f.go", cc: 2, lines: 5},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			candidates := parser.NewParser(helper.GuessLanguage(tc.path)).Parse()
			candidate.CalcScore(candidates)
			require.Len(t, candidates, 1)
			c := candidates[0]
			assert.Equal(t, tc.cc, c.Metrics.CyclomaticComplexity, "wrong cyclic complexity for function")
			assert.Equal(t, tc.lines, c.Metrics.LinesOfCode)
		})
	}
}
