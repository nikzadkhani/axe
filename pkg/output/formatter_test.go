package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nikzadk/axe/pkg/github"
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
