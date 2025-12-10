package branch

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/nikzadkhani/axe/pkg/git"
	"github.com/nikzadkhani/axe/pkg/github"
)

// MergedBranch represents a branch with its associated merged PR
type MergedBranch struct {
	Name string
	PR   *github.PRInfo
}

// BranchStatus represents a branch with its PR status
type BranchStatus struct {
	Name   string
	Status string // "merged", "open", "closed", "draft", "no-pr"
	PR     *github.PRInfo
}

// ProgressReporter is an interface for reporting progress during operations
type ProgressReporter interface {
	Start(msg string)
	Update(msg string)
	Stop(msg string)
	StopWithError(msg string)
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
func (s *Service) GetMergedBranches(repoPath string, reporter ProgressReporter) ([]MergedBranch, error) {
	// Get all local branches
	reporter.Start("Fetching local branches...")
	branches, err := s.gitClient.GetLocalBranches(repoPath)
	if err != nil {
		reporter.StopWithError(fmt.Sprintf("Failed to fetch local branches: %v", err))
		return nil, err
	}
	reporter.Stop(fmt.Sprintf("Found %d local branches", len(branches)))

	// Filter out main/master branches
	var filteredBranches []string
	for _, branch := range branches {
		if branch != "main" && branch != "master" && branch != "" {
			filteredBranches = append(filteredBranches, branch)
		}
	}

	if len(filteredBranches) == 0 {
		return []MergedBranch{}, nil
	}

	// Check each branch for merged PRs (parallelized)
	reporter.Start(fmt.Sprintf("Looking for branches to chop (%d to check)...", len(filteredBranches)))
	mergedBranches := s.checkBranchesParallel(repoPath, filteredBranches, reporter)
	reporter.Stop(fmt.Sprintf("Found %d branches ready to axe", len(mergedBranches)))

	return mergedBranches, nil
}

// checkBranchesParallel checks multiple branches concurrently using a worker pool
func (s *Service) checkBranchesParallel(repoPath string, branches []string, reporter ProgressReporter) []MergedBranch {
	// Use a worker pool to limit concurrent API calls
	const maxWorkers = 10
	numWorkers := min(maxWorkers, len(branches))

	// Channels for work distribution
	branchChan := make(chan string, len(branches))
	resultChan := make(chan MergedBranch, len(branches))

	// Progress tracking
	var processed atomic.Int32
	total := int32(len(branches))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for branch := range branchChan {
				pr, err := s.githubClient.GetMergedPR(repoPath, branch)

				// Update progress
				count := processed.Add(1)
				reporter.Update(fmt.Sprintf("Checking PR status (%d/%d)", count, total))

				// Only send result if branch has a merged PR
				if err == nil && pr != nil {
					resultChan <- MergedBranch{
						Name: branch,
						PR:   pr,
					}
				}
			}
		}()
	}

	// Send work to workers
	for _, branch := range branches {
		branchChan <- branch
	}
	close(branchChan)

	// Wait for all workers to finish, then close results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	var mergedBranches []MergedBranch
	for result := range resultChan {
		mergedBranches = append(mergedBranches, result)
	}

	return mergedBranches
}

// DeleteBranches deletes the specified branches
func (s *Service) DeleteBranches(repoPath string, branches []string, reporter ProgressReporter) (deleted []string, failed []string) {
	reporter.Start(fmt.Sprintf("Chopping %d branches...", len(branches)))
	for i, branch := range branches {
		reporter.Update(fmt.Sprintf("Chopping (%d/%d): %s", i+1, len(branches), branch))
		err := s.gitClient.DeleteBranch(repoPath, branch)
		if err != nil {
			failed = append(failed, branch)
		} else {
			deleted = append(deleted, branch)
		}
	}
	reporter.Stop(fmt.Sprintf("Chopped %d branches", len(deleted)))
	return deleted, failed
}

// GetAllBranchStatuses returns all local branches with their PR status
func (s *Service) GetAllBranchStatuses(repoPath string, reporter ProgressReporter) (map[string][]BranchStatus, error) {
	// Get all local branches
	reporter.Start("Fetching local branches...")
	branches, err := s.gitClient.GetLocalBranches(repoPath)
	if err != nil {
		reporter.StopWithError(fmt.Sprintf("Failed to fetch local branches: %v", err))
		return nil, err
	}
	reporter.Stop(fmt.Sprintf("Found %d local branches", len(branches)))

	// Filter out main/master branches
	var filteredBranches []string
	for _, branch := range branches {
		if branch != "main" && branch != "master" && branch != "" {
			filteredBranches = append(filteredBranches, branch)
		}
	}

	if len(filteredBranches) == 0 {
		return map[string][]BranchStatus{}, nil
	}

	// Check each branch for PR status (parallelized)
	reporter.Start(fmt.Sprintf("Checking PR status for %d branches...", len(filteredBranches)))
	statuses := s.checkAllBranchesParallel(repoPath, filteredBranches, reporter)
	reporter.Stop(fmt.Sprintf("Completed status check for %d branches", len(filteredBranches)))

	return statuses, nil
}

// checkAllBranchesParallel checks all branches concurrently and categorizes them by status
func (s *Service) checkAllBranchesParallel(repoPath string, branches []string, reporter ProgressReporter) map[string][]BranchStatus {
	// Use a worker pool to limit concurrent API calls
	const maxWorkers = 10
	numWorkers := min(maxWorkers, len(branches))

	// Channels for work distribution
	branchChan := make(chan string, len(branches))
	resultChan := make(chan BranchStatus, len(branches))

	// Progress tracking
	var processed atomic.Int32
	total := int32(len(branches))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for branch := range branchChan {
				pr, err := s.githubClient.GetPRStatus(repoPath, branch)

				// Update progress
				count := processed.Add(1)
				reporter.Update(fmt.Sprintf("Checking PR status (%d/%d)", count, total))

				// Determine status
				var status string
				if err != nil || pr == nil {
					status = "no-pr"
				} else if pr.IsDraft {
					status = "draft"
				} else if pr.State == "MERGED" {
					status = "merged"
				} else if pr.State == "OPEN" {
					status = "open"
				} else if pr.State == "CLOSED" {
					status = "closed"
				} else {
					status = "no-pr"
				}

				resultChan <- BranchStatus{
					Name:   branch,
					Status: status,
					PR:     pr,
				}
			}
		}()
	}

	// Send work to workers
	for _, branch := range branches {
		branchChan <- branch
	}
	close(branchChan)

	// Wait for all workers to finish, then close results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect and categorize results
	statusMap := map[string][]BranchStatus{
		"merged": {},
		"open":   {},
		"closed": {},
		"draft":  {},
		"no-pr":  {},
	}

	for result := range resultChan {
		statusMap[result.Status] = append(statusMap[result.Status], result)
	}

	return statusMap
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

