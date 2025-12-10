package progress

import (
	"fmt"
	"io"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

// Reporter provides an interface for reporting progress
type Reporter interface {
	// Start begins a progress operation with a message
	Start(msg string)
	// Update updates the current progress message
	Update(msg string)
	// Stop stops the progress indicator with a final message
	Stop(msg string)
	// StopWithError stops the progress indicator with an error message
	StopWithError(msg string)
}

// SpinnerReporter implements Reporter using a spinner
type SpinnerReporter struct {
	spinner *spinner.Spinner
	writer  io.Writer
}

// NewSpinnerReporter creates a new SpinnerReporter
func NewSpinnerReporter(w io.Writer) *SpinnerReporter {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(w))
	s.Writer = w
	return &SpinnerReporter{
		spinner: s,
		writer:  w,
	}
}

func (r *SpinnerReporter) Start(msg string) {
	r.spinner.Suffix = " " + msg
	r.spinner.Start()
}

func (r *SpinnerReporter) Update(msg string) {
	r.spinner.Suffix = " " + msg
}

func (r *SpinnerReporter) Stop(msg string) {
	r.spinner.Stop()
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Fprintf(r.writer, "%s %s\n", green("✓"), msg)
}

func (r *SpinnerReporter) StopWithError(msg string) {
	r.spinner.Stop()
	red := color.New(color.FgRed).SprintFunc()
	fmt.Fprintf(r.writer, "%s %s\n", red("✗"), msg)
}

// SilentReporter implements Reporter with no output (for testing)
type SilentReporter struct{}

// NewSilentReporter creates a new SilentReporter
func NewSilentReporter() *SilentReporter {
	return &SilentReporter{}
}

func (r *SilentReporter) Start(msg string)          {}
func (r *SilentReporter) Update(msg string)         {}
func (r *SilentReporter) Stop(msg string)           {}
func (r *SilentReporter) StopWithError(msg string)  {}
