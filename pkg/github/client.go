//go:generate mockgen -source=client.go -destination=mock_client.go -package=github

package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// PRInfo represents information about a pull request
type PRInfo struct {
	Number  int    `json:"number"`
	State   string `json:"state"`
	Title   string `json:"title"`
	IsDraft bool   `json:"isDraft"`
}

// Client provides an interface for GitHub operations
type Client interface {
	// GetMergedPR returns PR info for a merged PR associated with the branch
	// Returns nil if no merged PR is found
	GetMergedPR(repoPath, branch string) (*PRInfo, error)

	// GetPRStatus returns PR info for any PR associated with the branch
	// Returns nil if no PR is found, otherwise returns the most recent PR
	GetPRStatus(repoPath, branch string) (*PRInfo, error)
}

// DefaultClient implements Client using gh CLI
type DefaultClient struct{}

// NewDefaultClient creates a new DefaultClient
func NewDefaultClient() *DefaultClient {
	return &DefaultClient{}
}

func (c *DefaultClient) GetMergedPR(repoPath, branch string) (*PRInfo, error) {
	cmd := exec.Command("gh", "pr", "list",
		"--state", "merged",
		"--head", branch,
		"--json", "number,state,title,isDraft",
		"--limit", "1")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to check PR for branch %q: %w", branch, err)
	}

	var prs []PRInfo
	if err := json.Unmarshal(output, &prs); err != nil {
		return nil, fmt.Errorf("failed to parse PR data for branch %q: %w", branch, err)
	}

	if len(prs) == 0 {
		return nil, nil
	}

	return &prs[0], nil
}

func (c *DefaultClient) GetPRStatus(repoPath, branch string) (*PRInfo, error) {
	cmd := exec.Command("gh", "pr", "list",
		"--state", "all",
		"--head", branch,
		"--json", "number,state,title,isDraft",
		"--limit", "1")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to check PR status for branch %q: %w", branch, err)
	}

	var prs []PRInfo
	if err := json.Unmarshal(output, &prs); err != nil {
		return nil, fmt.Errorf("failed to parse PR data for branch %q: %w", branch, err)
	}

	if len(prs) == 0 {
		return nil, nil
	}

	return &prs[0], nil
}

