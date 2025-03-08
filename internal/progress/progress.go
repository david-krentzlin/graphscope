package progress

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// ProgressBar displays progress when total is known
type ProgressBar struct {
	current     int64
	total       int64
	width       int
	description string
	mu          sync.Mutex
	done        bool
}

// Spinner displays progress when total is unknown
type Spinner struct {
	frames      []string
	description string
	done        chan bool
	mu          sync.Mutex
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int64, width int) *ProgressBar {
	return &ProgressBar{
		total:       total,
		width:       width,
		description: "",
		done:        false,
	}
}

// Update updates the progress bar
func (p *ProgressBar) Update(current int64, description string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current = current
	p.description = formatDescription(description, 80)

	if !p.done {
		p.render()
	}
}

// render displays the current progress
func (p *ProgressBar) render() {
	percentage := float64(p.current) / float64(p.total)
	filled := int(float64(p.width) * percentage)

	// Ensure filled doesn't exceed width
	if filled > p.width {
		filled = p.width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", p.width-filled)
	fmt.Fprintf(os.Stderr, "\r[%s] %.2f%% %s", bar, percentage*100, p.description)
}

// Finish completes the progress bar
func (p *ProgressBar) Finish() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.done = true
	fmt.Fprintln(os.Stderr)
}

// NewSpinner creates a new spinner
func NewSpinner() *Spinner {
	return &Spinner{
		frames:      []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		description: "",
		done:        make(chan bool),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	go func() {
		i := 0
		for {
			select {
			case <-s.done:
				return
			default:
				s.mu.Lock()
				frame := s.frames[i%len(s.frames)]
				fmt.Fprintf(os.Stderr, "\r%s %s", frame, s.description)
				s.mu.Unlock()

				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// SetDescription updates the spinner text
func (s *Spinner) SetDescription(description string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.description = formatDescription(description, 80)
}

// Stop ends the spinner
func (s *Spinner) Stop() {
	s.done <- true
	fmt.Fprintln(os.Stderr)
}

func formatDescription(description string, width int) string {
	descLen := utf8.RuneCountInString(description)

	if descLen > width {
		runeSlice := []rune(description)
		return string(runeSlice[:width-3]) + "..."
	} else if descLen < width {
		// Pad with spaces if too short
		return description + strings.Repeat(" ", width-descLen)
	}

	return description
}
