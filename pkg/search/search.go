package search

import (
	"io/fs"
	"path/filepath"
	"sort"

	"github.com/CodeIntelligenceTesting/dlth/pkg/candidate"
	"github.com/CodeIntelligenceTesting/dlth/pkg/filter"
	"github.com/CodeIntelligenceTesting/dlth/pkg/helper"
	"github.com/CodeIntelligenceTesting/dlth/pkg/parser"
)

type Options struct {
	Filter     func(c *candidate.Candidate) bool
	Extensions []string
	Limit      int
}

func Search(srcPath string) (candidate.Candidates, error) {
	return SearchWithOptions(srcPath, Options{})
}

func SearchWithOptions(srcPath string, opts Options) (candidate.Candidates, error) {
	// walk over the given path and all child directories, parse the supported source code files
	// and collect possible candidates
	candidates := candidate.Candidates{}
	err := filepath.WalkDir(srcPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if filter.Valid(path, opts.Extensions) {
				nc := parser.NewParser(helper.GuessLanguage(path)).Parse()
				candidates = append(candidates, nc...)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	candidates.CalcScore()
	candidates = candidates.Filter(opts.Filter)

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	if opts.Limit > 0 {
		if opts.Limit > len(candidates) {
			opts.Limit = len(candidates)
		}
		candidates = candidates[:opts.Limit]
	}

	return candidates, nil
}
