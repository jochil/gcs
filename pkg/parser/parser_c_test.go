package parser_test

import (
	"testing"

	"github.com/jochil/gcs/pkg/candidate"
	"github.com/jochil/gcs/pkg/types"
)

func TestC(t *testing.T) {

	tests := []candidateTestCase{
		{
			name:         "main",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
		{
			name:         "a",
			params:       []*candidate.Parameter{},
			returnValues: []*candidate.Parameter{},
			visibility:   types.VisibilityPublic,
		},
	}
	runParserTests(t, tests, "testdata/c/function.c", types.C)
}
