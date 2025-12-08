package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "axe",
	Short: "Axe chops down squash-merged Git branches",
	Long: `Axe is a CLI tool that identifies and removes local Git branches
that have been squash-merged on GitHub but still exist locally.

It uses the GitHub CLI (gh) to check the merge status of pull requests
associated with your local branches.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("repo", "r", "", "Repository path (defaults to current directory)")
}
