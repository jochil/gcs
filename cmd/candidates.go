package cmd

import (
	"encoding/json"
	"io/fs"
	"path/filepath"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jochil/dlth/internal/tui"
	"github.com/jochil/dlth/pkg/candidate"
	"github.com/jochil/dlth/pkg/filter"
	"github.com/jochil/dlth/pkg/helper"
	"github.com/jochil/dlth/pkg/parser"
	"github.com/spf13/cobra"
)

var (
	printJSON  bool
	limit      int
	extensions []string

	versionCmd = &cobra.Command{
		Use:   "candidates",
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Short: "Scans for test candidates",
		RunE:  run,
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&printJSON, "json", false, "print results as json to stdout")
	versionCmd.Flags().IntVarP(&limit, "limit", "l", 0, "limit the amount of candidates (after sorting by score)")
	versionCmd.Flags().StringArrayVar(&extensions, "ext", []string{}, "only parse files with listed extension, flag can be used multiple times")
}

func run(cmd *cobra.Command, args []string) error {
	// TODO validate args
	srcPath, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}

	// walk over the given path and all child directories, parse the supported source code files
	// and collect possible candidates
	candidates := []*candidate.Candidate{}
	err = filepath.WalkDir(srcPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if filter.Valid(path, extensions) {
				nc := parser.NewParser(helper.GuessLanguage(path)).Parse()
				candidates = append(candidates, nc...)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	candidate.CalcScore(candidates)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	if limit > 0 {
		// TODO check for out of bounds
		if limit > len(candidates) {
			limit = len(candidates)
		}
		candidates = candidates[:limit]
	}

	if printJSON {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		err := enc.Encode(candidates)
		if err != nil {
			return err
		}
	} else {
		return startTUI(candidates, srcPath)
	}
	return nil
}

func startTUI(candidates []*candidate.Candidate, srcPath string) error {
	state, err := tui.NewCandidateModel(candidates, srcPath)
	if err != nil {
		return err
	}
	if _, err := tea.NewProgram(state).Run(); err != nil {
		return err
	}
	return nil
}
