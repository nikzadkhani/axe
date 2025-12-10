package cmd

import (
	"fmt"
	"os"

	"github.com/nikzadkhani/axe/pkg/branch"
	"github.com/nikzadkhani/axe/pkg/git"
	"github.com/nikzadkhani/axe/pkg/github"
	"github.com/nikzadkhani/axe/pkg/output"
	"github.com/nikzadkhani/axe/pkg/progress"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "branches",
	Aliases: []string{"list", "ls"},
	Short:   "List branches ready to axe 🪓",
	Long: `Find all local Git branches that have been squash-merged on GitHub
but are still hanging around locally.

These branches are ready to be chopped! Use 'axe chop' to remove them.`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("verbose", "v", false, "Show verbose output including PR numbers")
	listCmd.Flags().BoolP("all", "a", false, "Show all branches with their PR status (open, closed, no PR, etc.)")
}

func runList(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	showAll, _ := cmd.Flags().GetBool("all")
	repoPath, _ := cmd.Flags().GetString("repo")

	if repoPath == "" {
		repoPath = "."
	}

	// Create dependencies
	gitClient := git.NewDefaultClient()
	githubClient := github.NewDefaultClient()
	branchService := branch.NewService(gitClient, githubClient)

	// Create formatter based on --no-color flag
	noColor, _ := cmd.Flags().GetBool("no-color")
	var formatter output.Formatter
	if noColor {
		formatter = output.NewPlainFormatter(os.Stdout)
	} else {
		formatter = output.NewColoredFormatter(os.Stdout)
	}

	reporter := progress.NewSpinnerReporter(os.Stdout)

	// Validate repository
	if err := gitClient.ValidateRepository(repoPath); err != nil {
		formatter.PrintError(err.Error())
		return err
	}

	// Show all branch statuses or just merged branches
	if showAll {
		// Get all branch statuses
		statusMap, err := branchService.GetAllBranchStatuses(repoPath, reporter)
		if err != nil {
			formatter.PrintError(fmt.Sprintf("Failed to get branch statuses: %v", err))
			return err
		}

		fmt.Println() // Add spacing after spinner

		// Check if there are any branches
		totalBranches := 0
		for _, branches := range statusMap {
			totalBranches += len(branches)
		}

		if totalBranches == 0 {
			formatter.PrintInfo("No branches found! 🪓")
			return nil
		}

		// Display all statuses
		formatter.PrintBranchStatuses(statusMap)
	} else {
		// Get merged branches only (original behavior)
		mergedBranches, err := branchService.GetMergedBranches(repoPath, reporter)
		if err != nil {
			formatter.PrintError(fmt.Sprintf("Failed to get local branches: %v", err))
			return err
		}

		fmt.Println() // Add spacing after spinner

		// Display results
		if len(mergedBranches) == 0 {
			formatter.PrintInfo("No branches to axe! All clean 🪓")
			return nil
		}

		formatter.PrintHeader(fmt.Sprintf("🪓 Found %d branch(es) to axe:", len(mergedBranches)))
		for _, mb := range mergedBranches {
			if verbose {
				formatter.PrintBranchWithPR(mb.Name, mb.PR)
			} else {
				formatter.PrintBranch(mb.Name)
			}
		}
	}

	return nil
}
