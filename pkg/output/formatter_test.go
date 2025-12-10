package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nikzadkhani/axe/pkg/branch"
	"github.com/nikzadkhani/axe/pkg/github"
)

func TestColoredFormatter_PrintSuccess(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewColoredFormatter(buf)

	formatter.PrintSuccess("Operation completed")

	output := buf.String()
	if !strings.Contains(output, "Operation completed") {
		t.Errorf("PrintSuccess() output = %q, want it to contain 'Operation completed'", output)
	}
}

func TestColoredFormatter_PrintError(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewColoredFormatter(buf)

	formatter.PrintError("Something went wrong")

	output := buf.String()
	if !strings.Contains(output, "Something went wrong") {
		t.Errorf("PrintError() output = %q, want it to contain 'Something went wrong'", output)
	}
}

func TestColoredFormatter_PrintWarning(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewColoredFormatter(buf)

	formatter.PrintWarning("This is a warning")

	output := buf.String()
	if !strings.Contains(output, "This is a warning") {
		t.Errorf("PrintWarning() output = %q, want it to contain 'This is a warning'", output)
	}
}

func TestColoredFormatter_PrintInfo(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewColoredFormatter(buf)

	formatter.PrintInfo("Information message")

	output := buf.String()
	if !strings.Contains(output, "Information message") {
		t.Errorf("PrintInfo() output = %q, want it to contain 'Information message'", output)
	}
}

func TestColoredFormatter_PrintBranch(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewColoredFormatter(buf)

	formatter.PrintBranch("feature-branch")

	output := buf.String()
	if !strings.Contains(output, "feature-branch") {
		t.Errorf("PrintBranch() output = %q, want it to contain 'feature-branch'", output)
	}
}

func TestColoredFormatter_PrintBranchWithPR(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewColoredFormatter(buf)

	pr := &github.PRInfo{
		Number: 123,
		Title:  "Add new feature",
	}

	formatter.PrintBranchWithPR("feature-branch", pr)

	output := buf.String()
	if !strings.Contains(output, "feature-branch") {
		t.Errorf("PrintBranchWithPR() output = %q, want it to contain 'feature-branch'", output)
	}
	if !strings.Contains(output, "123") {
		t.Errorf("PrintBranchWithPR() output = %q, want it to contain '123'", output)
	}
	if !strings.Contains(output, "Add new feature") {
		t.Errorf("PrintBranchWithPR() output = %q, want it to contain 'Add new feature'", output)
	}
}

func TestColoredFormatter_PrintHeader(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewColoredFormatter(buf)

	formatter.PrintHeader("Section Header")

	output := buf.String()
	if !strings.Contains(output, "Section Header") {
		t.Errorf("PrintHeader() output = %q, want it to contain 'Section Header'", output)
	}
}

func TestPlainFormatter_PrintSuccess(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewPlainFormatter(buf)

	formatter.PrintSuccess("Operation completed")

	output := buf.String()
	if !strings.Contains(output, "Operation completed") {
		t.Errorf("PrintSuccess() output = %q, want it to contain 'Operation completed'", output)
	}
}

func TestPlainFormatter_PrintBranchWithPR(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewPlainFormatter(buf)

	pr := &github.PRInfo{
		Number: 123,
		Title:  "Add new feature",
	}

	formatter.PrintBranchWithPR("feature-branch", pr)

	output := buf.String()
	if !strings.Contains(output, "feature-branch") {
		t.Errorf("PrintBranchWithPR() output = %q, want it to contain 'feature-branch'", output)
	}
	if !strings.Contains(output, "PR #123") {
		t.Errorf("PrintBranchWithPR() output = %q, want it to contain 'PR #123'", output)
	}
	if !strings.Contains(output, "Add new feature") {
		t.Errorf("PrintBranchWithPR() output = %q, want it to contain 'Add new feature'", output)
	}
}

func TestColoredFormatter_PrintBranchStatuses(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewColoredFormatter(buf)

	statusMap := map[string][]branch.BranchStatus{
		"merged": {
			{Name: "merged-branch", Status: "merged", PR: &github.PRInfo{Number: 1, Title: "Merged PR"}},
		},
		"open": {
			{Name: "open-branch", Status: "open", PR: &github.PRInfo{Number: 2, Title: "Open PR"}},
		},
		"draft": {
			{Name: "draft-branch", Status: "draft", PR: &github.PRInfo{Number: 3, Title: "Draft PR"}},
		},
		"closed": {
			{Name: "closed-branch", Status: "closed", PR: &github.PRInfo{Number: 4, Title: "Closed PR"}},
		},
		"no-pr": {
			{Name: "no-pr-branch", Status: "no-pr", PR: nil},
		},
	}

	formatter.PrintBranchStatuses(statusMap)

	output := buf.String()

	// Check for all sections
	if !strings.Contains(output, "Merged (ready to axe)") {
		t.Errorf("PrintBranchStatuses() output should contain 'Merged (ready to axe)'")
	}
	if !strings.Contains(output, "Open PR") {
		t.Errorf("PrintBranchStatuses() output should contain 'Open PR'")
	}
	if !strings.Contains(output, "Draft PR") {
		t.Errorf("PrintBranchStatuses() output should contain 'Draft PR'")
	}
	if !strings.Contains(output, "Closed (not merged)") {
		t.Errorf("PrintBranchStatuses() output should contain 'Closed (not merged)'")
	}
	if !strings.Contains(output, "No PR") {
		t.Errorf("PrintBranchStatuses() output should contain 'No PR'")
	}

	// Check for branch names
	branches := []string{"merged-branch", "open-branch", "draft-branch", "closed-branch", "no-pr-branch"}
	for _, branchName := range branches {
		if !strings.Contains(output, branchName) {
			t.Errorf("PrintBranchStatuses() output should contain branch '%s'", branchName)
		}
	}
}

func TestPlainFormatter_PrintBranchStatuses(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewPlainFormatter(buf)

	statusMap := map[string][]branch.BranchStatus{
		"merged": {
			{Name: "merged-branch", Status: "merged", PR: &github.PRInfo{Number: 1, Title: "Merged PR"}},
		},
		"no-pr": {
			{Name: "no-pr-branch", Status: "no-pr", PR: nil},
		},
	}

	formatter.PrintBranchStatuses(statusMap)

	output := buf.String()

	// Check for sections
	if !strings.Contains(output, "Merged (ready to axe)") {
		t.Errorf("PrintBranchStatuses() output should contain 'Merged (ready to axe)'")
	}
	if !strings.Contains(output, "No PR") {
		t.Errorf("PrintBranchStatuses() output should contain 'No PR'")
	}

	// Check for branch names
	if !strings.Contains(output, "merged-branch") {
		t.Errorf("PrintBranchStatuses() output should contain 'merged-branch'")
	}
	if !strings.Contains(output, "no-pr-branch") {
		t.Errorf("PrintBranchStatuses() output should contain 'no-pr-branch'")
	}

	// Check PR details
	if !strings.Contains(output, "#1") {
		t.Errorf("PrintBranchStatuses() output should contain PR number '#1'")
	}
	if !strings.Contains(output, "Merged PR") {
		t.Errorf("PrintBranchStatuses() output should contain PR title 'Merged PR'")
	}
}
