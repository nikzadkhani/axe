package branch

import (
	"github.com/nikzadk/axe/pkg/git"
	"github.com/nikzadk/axe/pkg/github"
)

// MergedBranch represents a branch with its associated merged PR
type MergedBranch struct {
	Name string
	PR   *github.PRInfo
}

// Service orchestrates git and GitHub operations for branch management
type Service struct {
	gitClient    git.Client
	githubClient github.Client
}

// NewService creates a new branch Service
func NewService(gitClient git.Client, githubClient github.Client) *Service {
	return &Service{
		gitClient:    gitClient,
		githubClient: githubClient,
	}
}

// GetMergedBranches returns all local branches that have been squash-merged on GitHub
func (s *Service) GetMergedBranches(repoPath string) ([]MergedBranch, error) {
	// Get all local branches
	branches, err := s.gitClient.GetLocalBranches(repoPath)
	if err != nil {
		return nil, err
	}

	// Filter out main/master branches
	var filteredBranches []string
	for _, branch := range branches {
		if branch != "main" && branch != "master" && branch != "" {
			filteredBranches = append(filteredBranches, branch)
		}
	}

	// Check each branch for merged PRs
	var mergedBranches []MergedBranch
	for _, branch := range filteredBranches {
		pr, err := s.githubClient.GetMergedPR(repoPath, branch)
		if err == nil && pr != nil {
			mergedBranches = append(mergedBranches, MergedBranch{
				Name: branch,
				PR:   pr,
			})
		}
	}

	return mergedBranches, nil
}

// DeleteBranches deletes the specified branches
func (s *Service) DeleteBranches(repoPath string, branches []string) (deleted []string, failed []string) {
	for _, branch := range branches {
		err := s.gitClient.DeleteBranch(repoPath, branch)
		if err != nil {
			failed = append(failed, branch)
		} else {
			deleted = append(deleted, branch)
		}
	}
	return deleted, failed
}
