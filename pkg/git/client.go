//go:generate mockgen -source=client.go -destination=mock_client.go -package=git

package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Client provides an interface for git operations
type Client interface {
	// ValidateRepository checks if the path is a valid git repository
	ValidateRepository(repoPath string) error
	// GetLocalBranches returns all local branch names
	GetLocalBranches(repoPath string) ([]string, error)
	// DeleteBranch force-deletes a branch
	DeleteBranch(repoPath, branch string) error
}

// DefaultClient implements Client using git commands
type DefaultClient struct{}

// NewDefaultClient creates a new DefaultClient
func NewDefaultClient() *DefaultClient {
	return &DefaultClient{}
}

func (c *DefaultClient) ValidateRepository(repoPath string) error {
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "--git-dir")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not a git repository: %s", repoPath)
	}
	return nil
}

func (c *DefaultClient) GetLocalBranches(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "-C", repoPath, "branch", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get local branches: %w", err)
	}

	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	return branches, nil
}

func (c *DefaultClient) DeleteBranch(repoPath, branch string) error {
	cmd := exec.Command("git", "-C", repoPath, "branch", "-D", branch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete branch %q: %w", branch, err)
	}
	return nil
}
