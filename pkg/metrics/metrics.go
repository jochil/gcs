package metrics

import (
	"errors"
	"strings"

	"github.com/dominikbraun/graph"
)

type Metrics struct {
	LinesOfCode          int
	CyclomaticComplexity int
	FuzzFriendlyName     bool
}

func CountLines(sourceCode string) int {
	// TODO count actual lines.. no comments, no empty ones, ...
	lines := strings.Split(strings.ReplaceAll(sourceCode, "\r\n", "\n"), "\n")
	return len(lines)
}

func CalcCyclomaticComplexity(cfg graph.Graph[int, int]) (cc int, err error) {
	if cfg == nil {
		err = errors.New("no graph found")
		return
	}

	edges, err := cfg.Size()
	if err != nil {
		return
	}
	nodes, err := cfg.Order()
	if err != nil {
		return
	}
	cc = edges - nodes + 2
	return
}

func HasFuzzFriendlyName(name string) bool {
	lcSubstrings := []string{
		"encode",
		"decode",
		"parse",
		"encrypt",
		"decrypt",
		"open",
		"load",
	}
	lcName := strings.ToLower(name)

	for _, lcSubstring := range lcSubstrings {
		if strings.Contains(lcName, lcSubstring) {
			return true
		}
	}
	return false
}
