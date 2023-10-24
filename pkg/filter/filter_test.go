package filter_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/filter"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	tests := map[string]struct {
		path       string
		extensions []string
		result     bool
	}{
		"go":      {path: "foo.go", extensions: []string{}, result: true},
		"go_no":   {path: "foo.go", extensions: []string{".java"}, result: false},
		"go_test": {path: "foo_test.go", extensions: []string{}, result: false},
		"java":    {path: "foo.java", extensions: []string{}, result: true},
		"js":      {path: "foo.js", extensions: []string{}, result: true},
		"c":       {path: "foo.c", extensions: []string{}, result: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := filter.Valid(tc.path, tc.extensions)
			assert.Equal(t, tc.result, result)
		})
	}

}
