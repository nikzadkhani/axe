package branch

import (
	"errors"
	"testing"

	"github.com/nikzadkhani/axe/pkg/git"
	"github.com/nikzadkhani/axe/pkg/github"
	"go.uber.org/mock/gomock"
)

// mockReporter implements ProgressReporter for testing
type mockReporter struct{}

func (m *mockReporter) Start(msg string)         {}
func (m *mockReporter) Update(msg string)        {}
func (m *mockReporter) Stop(msg string)          {}
func (m *mockReporter) StopWithError(msg string) {}

func TestService_GetMergedBranches(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(*git.MockClient, *github.MockClient)
		repoPath      string
		wantBranches  int
		wantErr       bool
		expectedNames []string
	}{
		{
			name: "returns merged branches successfully",
			setupMocks: func(gitMock *git.MockClient, ghMock *github.MockClient) {
				gitMock.EXPECT().
					GetLocalBranches(".").
					Return([]string{"main", "feature-1", "feature-2"}, nil)

				ghMock.EXPECT().
					GetMergedPR(".", "feature-1").
					Return(&github.PRInfo{Number: 1, State: "merged", Title: "Feature 1"}, nil)

				ghMock.EXPECT().
					GetMergedPR(".", "feature-2").
					Return(nil, nil)
			},
			repoPath:      ".",
			wantBranches:  1,
			wantErr:       false,
			expectedNames: []string{"feature-1"},
		},
		{
			name: "filters out main and master branches",
			setupMocks: func(gitMock *git.MockClient, ghMock *github.MockClient) {
				gitMock.EXPECT().
					GetLocalBranches(".").
					Return([]string{"main", "master", "feature-1"}, nil)

				ghMock.EXPECT().
					GetMergedPR(".", "feature-1").
					Return(&github.PRInfo{Number: 1, State: "merged", Title: "Feature 1"}, nil)
			},
			repoPath:      ".",
			wantBranches:  1,
			wantErr:       false,
			expectedNames: []string{"feature-1"},
		},
		{
			name: "returns error when git client fails",
			setupMocks: func(gitMock *git.MockClient, ghMock *github.MockClient) {
				gitMock.EXPECT().
					GetLocalBranches(".").
					Return(nil, errors.New("git error"))
			},
			repoPath:     ".",
			wantBranches: 0,
			wantErr:      true,
		},
		{
			name: "returns empty list when no merged branches found",
			setupMocks: func(gitMock *git.MockClient, ghMock *github.MockClient) {
				gitMock.EXPECT().
					GetLocalBranches(".").
					Return([]string{"main", "feature-1"}, nil)

				ghMock.EXPECT().
					GetMergedPR(".", "feature-1").
					Return(nil, nil)
			},
			repoPath:      ".",
			wantBranches:  0,
			wantErr:       false,
			expectedNames: []string{},
		},
		{
			name: "handles empty branch list",
			setupMocks: func(gitMock *git.MockClient, ghMock *github.MockClient) {
				gitMock.EXPECT().
					GetLocalBranches(".").
					Return([]string{}, nil)
			},
			repoPath:      ".",
			wantBranches:  0,
			wantErr:       false,
			expectedNames: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gitMock := git.NewMockClient(ctrl)
			ghMock := github.NewMockClient(ctrl)

			tt.setupMocks(gitMock, ghMock)

			service := NewService(gitMock, ghMock)
			reporter := &mockReporter{}
			branches, err := service.GetMergedBranches(tt.repoPath, reporter)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetMergedBranches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(branches) != tt.wantBranches {
				t.Errorf("GetMergedBranches() got %d branches, want %d", len(branches), tt.wantBranches)
			}

			if !tt.wantErr && tt.expectedNames != nil {
				for i, name := range tt.expectedNames {
					if branches[i].Name != name {
						t.Errorf("GetMergedBranches() branch[%d] = %s, want %s", i, branches[i].Name, name)
					}
				}
			}
		})
	}
}

func TestService_DeleteBranches(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*git.MockClient)
		branches       []string
		wantDeleted    int
		wantFailed     int
		expectedDeleted []string
		expectedFailed  []string
	}{
		{
			name: "deletes all branches successfully",
			setupMocks: func(gitMock *git.MockClient) {
				gitMock.EXPECT().DeleteBranch(".", "feature-1").Return(nil)
				gitMock.EXPECT().DeleteBranch(".", "feature-2").Return(nil)
			},
			branches:        []string{"feature-1", "feature-2"},
			wantDeleted:     2,
			wantFailed:      0,
			expectedDeleted: []string{"feature-1", "feature-2"},
			expectedFailed:  []string{},
		},
		{
			name: "handles partial deletion failure",
			setupMocks: func(gitMock *git.MockClient) {
				gitMock.EXPECT().DeleteBranch(".", "feature-1").Return(nil)
				gitMock.EXPECT().DeleteBranch(".", "feature-2").Return(errors.New("delete error"))
			},
			branches:        []string{"feature-1", "feature-2"},
			wantDeleted:     1,
			wantFailed:      1,
			expectedDeleted: []string{"feature-1"},
			expectedFailed:  []string{"feature-2"},
		},
		{
			name: "handles all deletions failing",
			setupMocks: func(gitMock *git.MockClient) {
				gitMock.EXPECT().DeleteBranch(".", "feature-1").Return(errors.New("delete error"))
			},
			branches:        []string{"feature-1"},
			wantDeleted:     0,
			wantFailed:      1,
			expectedDeleted: []string{},
			expectedFailed:  []string{"feature-1"},
		},
		{
			name:            "handles empty branch list",
			setupMocks:      func(gitMock *git.MockClient) {},
			branches:        []string{},
			wantDeleted:     0,
			wantFailed:      0,
			expectedDeleted: []string{},
			expectedFailed:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gitMock := git.NewMockClient(ctrl)
			ghMock := github.NewMockClient(ctrl)

			tt.setupMocks(gitMock)

			service := NewService(gitMock, ghMock)
			reporter := &mockReporter{}
			deleted, failed := service.DeleteBranches(".", tt.branches, reporter)

			if len(deleted) != tt.wantDeleted {
				t.Errorf("DeleteBranches() deleted %d branches, want %d", len(deleted), tt.wantDeleted)
			}

			if len(failed) != tt.wantFailed {
				t.Errorf("DeleteBranches() failed %d branches, want %d", len(failed), tt.wantFailed)
			}

			for i, name := range tt.expectedDeleted {
				if i >= len(deleted) || deleted[i] != name {
					t.Errorf("DeleteBranches() deleted[%d] = %v, want %s", i, deleted, name)
				}
			}

			for i, name := range tt.expectedFailed {
				if i >= len(failed) || failed[i] != name {
					t.Errorf("DeleteBranches() failed[%d] = %v, want %s", i, failed, name)
				}
			}
		})
	}
}

func TestService_GetAllBranchStatuses(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(*git.MockClient, *github.MockClient)
		repoPath   string
		wantErr    bool
		wantCounts map[string]int
	}{
		{
			name: "categorizes branches by status",
			setupMocks: func(gitMock *git.MockClient, ghMock *github.MockClient) {
				gitMock.EXPECT().
					GetLocalBranches(".").
					Return([]string{"main", "merged-1", "open-1", "draft-1", "closed-1", "no-pr-1"}, nil)

				ghMock.EXPECT().
					GetPRStatus(".", "merged-1").
					Return(&github.PRInfo{Number: 1, State: "MERGED", Title: "Merged PR", IsDraft: false}, nil)

				ghMock.EXPECT().
					GetPRStatus(".", "open-1").
					Return(&github.PRInfo{Number: 2, State: "OPEN", Title: "Open PR", IsDraft: false}, nil)

				ghMock.EXPECT().
					GetPRStatus(".", "draft-1").
					Return(&github.PRInfo{Number: 3, State: "OPEN", Title: "Draft PR", IsDraft: true}, nil)

				ghMock.EXPECT().
					GetPRStatus(".", "closed-1").
					Return(&github.PRInfo{Number: 4, State: "CLOSED", Title: "Closed PR", IsDraft: false}, nil)

				ghMock.EXPECT().
					GetPRStatus(".", "no-pr-1").
					Return(nil, nil)
			},
			repoPath: ".",
			wantErr:  false,
			wantCounts: map[string]int{
				"merged": 1,
				"open":   1,
				"draft":  1,
				"closed": 1,
				"no-pr":  1,
			},
		},
		{
			name: "filters out main and master branches",
			setupMocks: func(gitMock *git.MockClient, ghMock *github.MockClient) {
				gitMock.EXPECT().
					GetLocalBranches(".").
					Return([]string{"main", "master", "feature-1"}, nil)

				ghMock.EXPECT().
					GetPRStatus(".", "feature-1").
					Return(&github.PRInfo{Number: 1, State: "MERGED", Title: "Feature 1", IsDraft: false}, nil)
			},
			repoPath: ".",
			wantErr:  false,
			wantCounts: map[string]int{
				"merged": 1,
				"open":   0,
				"draft":  0,
				"closed": 0,
				"no-pr":  0,
			},
		},
		{
			name: "returns error when git client fails",
			setupMocks: func(gitMock *git.MockClient, ghMock *github.MockClient) {
				gitMock.EXPECT().
					GetLocalBranches(".").
					Return(nil, errors.New("git error"))
			},
			repoPath: ".",
			wantErr:  true,
		},
		{
			name: "handles empty branch list",
			setupMocks: func(gitMock *git.MockClient, ghMock *github.MockClient) {
				gitMock.EXPECT().
					GetLocalBranches(".").
					Return([]string{}, nil)
			},
			repoPath: ".",
			wantErr:  false,
			wantCounts: map[string]int{
				"merged": 0,
				"open":   0,
				"draft":  0,
				"closed": 0,
				"no-pr":  0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gitMock := git.NewMockClient(ctrl)
			ghMock := github.NewMockClient(ctrl)

			tt.setupMocks(gitMock, ghMock)

			service := NewService(gitMock, ghMock)
			reporter := &mockReporter{}
			statusMap, err := service.GetAllBranchStatuses(tt.repoPath, reporter)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllBranchStatuses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				for status, expectedCount := range tt.wantCounts {
					if len(statusMap[status]) != expectedCount {
						t.Errorf("GetAllBranchStatuses() status %s has %d branches, want %d",
							status, len(statusMap[status]), expectedCount)
					}
				}
			}
		})
	}
}
