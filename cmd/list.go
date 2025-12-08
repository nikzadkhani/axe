package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

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

type prInfo struct {
	Number int    `json:"number"`
	State  string `json:"state"`
	Title  string `json:"title"`
}

func runList(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
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
	mergedBranches := make(map[string]prInfo)
	for _, branch := range filteredBranches {
		pr, err := getMergedPR(repoPath, branch)
		if err == nil && pr != nil {
			mergedBranches[branch] = *pr
		}
	}

	// Display results
	if len(mergedBranches) == 0 {
		fmt.Println("No squash-merged branches found.")
		return nil
	}

	fmt.Printf("Found %d squash-merged branch(es):\n\n", len(mergedBranches))
	for branch, pr := range mergedBranches {
		if verbose {
			fmt.Printf("  %s (PR #%d: %s)\n", branch, pr.Number, pr.Title)
		} else {
			fmt.Printf("  %s\n", branch)
		}
	}

	return nil
}

func getLocalBranches(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "-C", repoPath, "branch", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	return branches, nil
}

func getMergedPR(repoPath, branch string) (*prInfo, error) {
	// Use gh to check for merged PRs for this branch
	cmd := exec.Command("gh", "pr", "list",
		"--repo", "$(git -C "+repoPath+" remote get-url origin | sed 's/.*github.com[:/]\\(.*\\)\\.git/\\1/')",
		"--state", "merged",
		"--head", branch,
		"--json", "number,state,title",
		"--limit", "1")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var prs []prInfo
	if err := json.Unmarshal(output, &prs); err != nil {
		return nil, err
	}

	if len(prs) == 0 {
		return nil, fmt.Errorf("no merged PR found")
	}

	return &prs[0], nil
}
