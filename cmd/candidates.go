package cmd

import (
	"encoding/json"
	"path/filepath"

	"github.com/jochil/dlth/internal/tui"
	"github.com/jochil/dlth/pkg/candidate"
	"github.com/jochil/dlth/pkg/search"
	tea "github.com/charmbracelet/bubbletea"
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
	versionCmd.Flags().BoolVar(&printJSON, "json", false, "print results as json to stdout")
	versionCmd.Flags().IntVarP(&limit, "limit", "l", 0, "limit the amount of candidates (after sorting by score)")
	versionCmd.Flags().StringArrayVar(&extensions, "ext", []string{}, "only parse files with listed extension, flag can be used multiple times")
	rootCmd.AddCommand(versionCmd)
}

func run(cmd *cobra.Command, args []string) error {
	// TODO validate args
	srcPath, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}

	candidates, err := search.SearchWithOptions(srcPath, search.Options{
		Limit:      limit,
		Extensions: extensions,
	})

	if err != nil {
		return err
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

func startTUI(candidates candidate.Candidates, srcPath string) error {
	state, err := tui.NewCandidateModel(candidates, srcPath)
	if err != nil {
		return err
	}
	if _, err := tea.NewProgram(state).Run(); err != nil {
		return err
	}
	return nil
}
