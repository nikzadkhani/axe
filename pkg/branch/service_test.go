package branch

import (
	"errors"
	"testing"

	"github.com/nikzadk/axe/pkg/git"
	"github.com/nikzadk/axe/pkg/github"
	"go.uber.org/mock/gomock"
)

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
			branches, err := service.GetMergedBranches(tt.repoPath)

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
			deleted, failed := service.DeleteBranches(".", tt.branches)

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
