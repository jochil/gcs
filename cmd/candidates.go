package cmd

import (
	"io/fs"
	"path/filepath"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jochil/dlth/internal/tui"
	"github.com/jochil/dlth/pkg/metric"
	"github.com/jochil/dlth/pkg/parser"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "candidates",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Scans for test candidates",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		state, err := tui.NewCandidateModel(candidates, srcPath)
		if err != nil {
			return err
		}
		if _, err := tea.NewProgram(state).Run(); err != nil {
			return err
		}

		return nil
	},
}
