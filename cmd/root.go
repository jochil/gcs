package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
  Use:   "dlth",
  Short: "Dirty Little Test Helper",
  Long: `Finds candidates for good (fuzz|unit) tests and automatically generates test functions`,
  Run: func(cmd *cobra.Command, args []string) {
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
