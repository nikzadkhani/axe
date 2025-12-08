package cmd

import (
	"fmt"
	"os"

	"github.com/nikzadk/axe/pkg/branch"
	"github.com/nikzadk/axe/pkg/git"
	"github.com/nikzadk/axe/pkg/github"
	"github.com/nikzadk/axe/pkg/output"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List local branches that were squash-merged on GitHub",
	Long: `List all local Git branches that have associated pull requests
that were merged on GitHub but still exist locally.

This command uses the GitHub CLI (gh) to check PR status.`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("verbose", "v", false, "Show verbose output including PR numbers")
}

func runList(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	repoPath, _ := cmd.Flags().GetString("repo")

	if repoPath == "" {
		repoPath = "."
	}

	// Create dependencies
	gitClient := git.NewDefaultClient()
	githubClient := github.NewDefaultClient()
	branchService := branch.NewService(gitClient, githubClient)
	formatter := output.NewColoredFormatter(os.Stdout)

	// Get merged branches
	mergedBranches, err := branchService.GetMergedBranches(repoPath)
	if err != nil {
		formatter.PrintError(fmt.Sprintf("Failed to get local branches: %v", err))
		return err
	}

	// Display results
	if len(mergedBranches) == 0 {
		formatter.PrintInfo("No squash-merged branches found.")
		return nil
	}

	formatter.PrintHeader(fmt.Sprintf("Found %d squash-merged branch(es):", len(mergedBranches)))
	for _, mb := range mergedBranches {
		if verbose {
			formatter.PrintBranchWithPR(mb.Name, mb.PR)
		} else {
			formatter.PrintBranch(mb.Name)
		}
	}

	return nil
}
