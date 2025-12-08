//go:generate mockgen -source=client.go -destination=mock_client.go -package=git

package git

import (
	"os/exec"
	"strings"
)

// Client provides an interface for git operations
type Client interface {
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

func (c *DefaultClient) GetLocalBranches(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "-C", repoPath, "branch", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	return branches, nil
}

func (c *DefaultClient) DeleteBranch(repoPath, branch string) error {
	cmd := exec.Command("git", "-C", repoPath, "branch", "-D", branch)
	return cmd.Run()
}
