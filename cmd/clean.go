package cmd

import (
	"fmt"
	"os"
	"os/exec"

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

	// Get all local branches
	branches, err := getLocalBranches(repoPath)
	if err != nil {
		return fmt.Errorf("failed to get local branches: %w", err)
	}

	// Filter out main/master branches
	var filteredBranches []string
	for _, branch := range branches {
		if branch != "main" && branch != "master" {
			filteredBranches = append(filteredBranches, branch)
		}
	}

	// Check each branch for merged PRs
	var toDelete []string
	for _, branch := range filteredBranches {
		pr, err := getMergedPR(repoPath, branch)
		if err == nil && pr != nil {
			toDelete = append(toDelete, branch)
		}
	}

	if len(toDelete) == 0 {
		fmt.Println("No squash-merged branches found to delete.")
		return nil
	}

	// Display what will be deleted
	fmt.Printf("Found %d squash-merged branch(es) to delete:\n\n", len(toDelete))
	for _, branch := range toDelete {
		fmt.Printf("  %s\n", branch)
	}
	fmt.Println()

	if dryRun {
		fmt.Println("(Dry run - no branches were deleted)")
		return nil
	}

	// Confirm deletion unless force flag is set
	if !force {
		fmt.Print("Delete these branches? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// Delete branches
	deleted := 0
	failed := 0
	for _, branch := range toDelete {
		cmd := exec.Command("git", "-C", repoPath, "branch", "-D", branch)
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to delete %s: %v\n", branch, err)
			failed++
		} else {
			fmt.Printf("Deleted: %s\n", branch)
			deleted++
		}
	}

	fmt.Printf("\nDeleted %d branch(es)", deleted)
	if failed > 0 {
		fmt.Printf(" (%d failed)", failed)
	}
	fmt.Println()

	return nil
}
