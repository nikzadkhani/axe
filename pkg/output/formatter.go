package output

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/nikzadkhani/axe/pkg/branch"
	"github.com/nikzadkhani/axe/pkg/github"
)

// Formatter provides an interface for formatted output
type Formatter interface {
	// PrintSuccess prints a success message
	PrintSuccess(msg string)
	// PrintError prints an error message
	PrintError(msg string)
	// PrintWarning prints a warning message
	PrintWarning(msg string)
	// PrintInfo prints an info message
	PrintInfo(msg string)
	// PrintBranch prints a branch name
	PrintBranch(branch string)
	// PrintBranchWithPR prints a branch with PR info
	PrintBranchWithPR(branch string, pr *github.PRInfo)
	// PrintHeader prints a header message
	PrintHeader(msg string)
	// PrintBranchStatuses prints branches grouped by status
	PrintBranchStatuses(statusMap map[string][]branch.BranchStatus)
}

// ColoredFormatter implements Formatter with colored output
type ColoredFormatter struct {
	writer io.Writer
}

// NewColoredFormatter creates a new ColoredFormatter
func NewColoredFormatter(w io.Writer) *ColoredFormatter {
	return &ColoredFormatter{writer: w}
}

func (f *ColoredFormatter) PrintSuccess(msg string) {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Fprintf(f.writer, "%s %s\n", green("✓"), msg)
}

func (f *ColoredFormatter) PrintError(msg string) {
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	fmt.Fprintf(f.writer, "%s %s\n", red("✗"), msg)
}

func (f *ColoredFormatter) PrintWarning(msg string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Fprintf(f.writer, "%s %s\n", yellow("⚠"), msg)
}

func (f *ColoredFormatter) PrintInfo(msg string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	fmt.Fprintf(f.writer, "%s %s\n", cyan("ℹ"), msg)
}

func (f *ColoredFormatter) PrintBranch(branch string) {
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	fmt.Fprintf(f.writer, "  %s\n", green(branch))
}

func (f *ColoredFormatter) PrintBranchWithPR(branch string, pr *github.PRInfo) {
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Fprintf(f.writer, "  %s %s %s\n",
		green(branch),
		yellow(fmt.Sprintf("(PR #%d)", pr.Number)),
		cyan(pr.Title))
}

func (f *ColoredFormatter) PrintHeader(msg string) {
	bold := color.New(color.Bold).SprintFunc()
	fmt.Fprintf(f.writer, "\n%s\n", bold(msg))
}

func (f *ColoredFormatter) PrintBranchStatuses(statusMap map[string][]branch.BranchStatus) {
	// Define status order and formatting
	statusInfo := []struct {
		key    string
		emoji  string
		color  func(a ...interface{}) string
		label  string
	}{
		{"merged", "🪓", color.New(color.FgGreen, color.Bold).SprintFunc(), "Merged (ready to axe)"},
		{"open", "📂", color.New(color.FgCyan).SprintFunc(), "Open PR"},
		{"draft", "✏️", color.New(color.FgMagenta).SprintFunc(), "Draft PR"},
		{"closed", "❌", color.New(color.FgRed).SprintFunc(), "Closed (not merged)"},
		{"no-pr", "🔍", color.New(color.FgYellow).SprintFunc(), "No PR"},
	}

	for _, info := range statusInfo {
		branches := statusMap[info.key]
		if len(branches) == 0 {
			continue
		}

		// Print section header
		headerColor := color.New(color.Bold).SprintFunc()
		fmt.Fprintf(f.writer, "\n%s %s: %d branch(es)\n",
			info.emoji,
			headerColor(info.label),
			len(branches))

		// Print branches
		for _, b := range branches {
			if b.PR != nil {
				yellow := color.New(color.FgYellow).SprintFunc()
				dim := color.New(color.Faint).SprintFunc()
				fmt.Fprintf(f.writer, "  %s %s %s\n",
					info.color(b.Name),
					yellow(fmt.Sprintf("(#%d)", b.PR.Number)),
					dim(b.PR.Title))
			} else {
				fmt.Fprintf(f.writer, "  %s\n", info.color(b.Name))
			}
		}
	}
}

// PlainFormatter implements Formatter with plain text output
type PlainFormatter struct {
	writer io.Writer
}

// NewPlainFormatter creates a new PlainFormatter
func NewPlainFormatter(w io.Writer) *PlainFormatter {
	return &PlainFormatter{writer: w}
}

func (f *PlainFormatter) PrintSuccess(msg string) {
	fmt.Fprintf(f.writer, "✓ %s\n", msg)
}

func (f *PlainFormatter) PrintError(msg string) {
	fmt.Fprintf(f.writer, "✗ %s\n", msg)
}

func (f *PlainFormatter) PrintWarning(msg string) {
	fmt.Fprintf(f.writer, "⚠ %s\n", msg)
}

func (f *PlainFormatter) PrintInfo(msg string) {
	fmt.Fprintf(f.writer, "ℹ %s\n", msg)
}

func (f *PlainFormatter) PrintBranch(branch string) {
	fmt.Fprintf(f.writer, "  %s\n", branch)
}

func (f *PlainFormatter) PrintBranchWithPR(branch string, pr *github.PRInfo) {
	fmt.Fprintf(f.writer, "  %s (PR #%d: %s)\n", branch, pr.Number, pr.Title)
}

func (f *PlainFormatter) PrintHeader(msg string) {
	fmt.Fprintf(f.writer, "\n%s\n", msg)
}

func (f *PlainFormatter) PrintBranchStatuses(statusMap map[string][]branch.BranchStatus) {
	// Define status order
	statusInfo := []struct {
		key   string
		emoji string
		label string
	}{
		{"merged", "🪓", "Merged (ready to axe)"},
		{"open", "📂", "Open PR"},
		{"draft", "✏️", "Draft PR"},
		{"closed", "❌", "Closed (not merged)"},
		{"no-pr", "🔍", "No PR"},
	}

	for _, info := range statusInfo {
		branches := statusMap[info.key]
		if len(branches) == 0 {
			continue
		}

		// Print section header
		fmt.Fprintf(f.writer, "\n%s %s: %d branch(es)\n",
			info.emoji,
			info.label,
			len(branches))

		// Print branches
		for _, b := range branches {
			if b.PR != nil {
				fmt.Fprintf(f.writer, "  %s (#%d: %s)\n",
					b.Name,
					b.PR.Number,
					b.PR.Title)
			} else {
				fmt.Fprintf(f.writer, "  %s\n", b.Name)
			}
		}
	}
}

