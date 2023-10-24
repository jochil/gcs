package cmd

import (
	"encoding/json"
	"io/fs"
	"path/filepath"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jochil/dlth/internal/tui"
	"github.com/jochil/dlth/pkg/metric"
	"github.com/jochil/dlth/pkg/parser"
	"github.com/spf13/cobra"
)

var (
	printJSON bool

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
}

func run(cmd *cobra.Command, args []string) error {
	// TODO validate args
	srcPath, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}

	// walk over the given path and all child directories, parse the supported source code files
	// and collect possible candidates
	candidates := []*parser.Candidate{}
	err = filepath.WalkDir(srcPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if language, ok := parser.SupportedExt[filepath.Ext(path)]; ok {
				nc := parser.NewParser(path, language).Parse()
				candidates = append(candidates, nc...)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	metric.CalcScore(candidates)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

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

func startTUI(candidates []*parser.Candidate, srcPath string) error {
	state, err := tui.NewCandidateModel(candidates, srcPath)
	if err != nil {
		return err
	}
	if _, err := tea.NewProgram(state).Run(); err != nil {
		return err
	}
	return nil
}
