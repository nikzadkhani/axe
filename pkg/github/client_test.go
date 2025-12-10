package github

import (
	"testing"
)

func TestGetPRStatus(t *testing.T) {
	tests := []struct {
		name     string
		repoPath string
		branch   string
		wantErr  bool
	}{
		{
			name:     "valid branch with PR",
			repoPath: ".",
			branch:   "feature/test",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewDefaultClient()
			_, err := client.GetPRStatus(tt.repoPath, tt.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPRStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
