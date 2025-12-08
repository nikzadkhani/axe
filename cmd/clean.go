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

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Delete local branches that were squash-merged on GitHub",
	Long: `Delete all local Git branches that have associated pull requests
that were merged on GitHub.

This command uses the GitHub CLI (gh) to check PR status before deletion.
Branches are force-deleted using 'git branch -D' since squash-merged
branches don't show as merged in git's history.`,
	RunE: runClean,
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolP("dry-run", "n", false, "Show what would be deleted without actually deleting")
	cleanCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}

func runClean(cmd *cobra.Command, args []string) error {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	force, _ := cmd.Flags().GetBool("force")
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

	if len(mergedBranches) == 0 {
		formatter.PrintInfo("No squash-merged branches found to delete.")
		return nil
	}

	// Display what will be deleted
	formatter.PrintHeader(fmt.Sprintf("Found %d squash-merged branch(es) to delete:", len(mergedBranches)))
	for _, mb := range mergedBranches {
		formatter.PrintBranch(mb.Name)
	}
	fmt.Println()

	if dryRun {
		formatter.PrintWarning("Dry run - no branches were deleted")
		return nil
	}

	// Confirm deletion unless force flag is set
	if !force {
		fmt.Print("Delete these branches? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			formatter.PrintInfo("Cancelled.")
			return nil
		}
	}

	// Extract branch names
	var branchNames []string
	for _, mb := range mergedBranches {
		branchNames = append(branchNames, mb.Name)
	}

	// Delete branches
	deleted, failed := branchService.DeleteBranches(repoPath, branchNames)

	// Display results
	fmt.Println()
	for _, branch := range deleted {
		formatter.PrintSuccess(fmt.Sprintf("Deleted: %s", branch))
	}
	for _, branch := range failed {
		formatter.PrintError(fmt.Sprintf("Failed to delete: %s", branch))
	}

	fmt.Println()
	if len(failed) > 0 {
		formatter.PrintWarning(fmt.Sprintf("Deleted %d branch(es), %d failed", len(deleted), len(failed)))
	} else {
		formatter.PrintSuccess(fmt.Sprintf("Deleted %d branch(es)", len(deleted)))
	}

	return nil
}
