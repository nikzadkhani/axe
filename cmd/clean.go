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

var cleanCmd = &cobra.Command{
	Use:     "chop",
	Aliases: []string{"clean", "delete", "rm"},
	Short:   "Chop down squash-merged branches 🪓",
	Long: `Chop down all local Git branches that have been squash-merged on GitHub
but are still taking up space locally.

This command checks GitHub for merged PRs before swinging the axe.
Branches are force-deleted since squash-merged commits don't show up
in git's history.`,
	RunE: runClean,
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolP("dry-run", "n", false, "Show what would be chopped without actually chopping")
	cleanCmd.Flags().BoolP("force", "f", false, "Skip confirmation and start chopping")
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

	// Get merged branches
	mergedBranches, err := branchService.GetMergedBranches(repoPath, reporter)
	if err != nil {
		formatter.PrintError(fmt.Sprintf("Failed to get local branches: %v", err))
		return err
	}

	fmt.Println() // Add spacing after spinner

	if len(mergedBranches) == 0 {
		formatter.PrintInfo("No branches to chop! All clean 🪓")
		return nil
	}

	// Display what will be chopped
	formatter.PrintHeader(fmt.Sprintf("🪓 Found %d branch(es) ready to chop:", len(mergedBranches)))
	for _, mb := range mergedBranches {
		formatter.PrintBranch(mb.Name)
	}
	fmt.Println()

	if dryRun {
		formatter.PrintWarning("Dry run - no branches were chopped")
		return nil
	}

	// Confirm deletion unless force flag is set
	if !force {
		fmt.Print("🪓 Chop these branches? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			formatter.PrintInfo("Cancelled. No branches were chopped.")
			return nil
		}
	}

	// Extract branch names
	var branchNames []string
	for _, mb := range mergedBranches {
		branchNames = append(branchNames, mb.Name)
	}

	// Delete branches
	deleted, failed := branchService.DeleteBranches(repoPath, branchNames, reporter)

	// Display results
	fmt.Println()
	for _, branch := range deleted {
		formatter.PrintSuccess(fmt.Sprintf("Chopped: %s", branch))
	}
	for _, branch := range failed {
		formatter.PrintError(fmt.Sprintf("Failed to chop: %s", branch))
	}

	fmt.Println()
	if len(failed) > 0 {
		formatter.PrintWarning(fmt.Sprintf("🪓 Chopped %d branch(es), %d failed", len(deleted), len(failed)))
	} else {
		formatter.PrintSuccess(fmt.Sprintf("🪓 Chopped %d branch(es)!", len(deleted)))
	}

	return nil
}
