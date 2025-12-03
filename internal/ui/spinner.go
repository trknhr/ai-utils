package ui

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// ASCII art logo for aiu (purple to pink gradient)
const Logo = "\033[38;5;135m █▀█ █ █ █\n\033[38;5;171m █▀█ █ █ █\n\033[38;5;213m ▀ ▀ ▀ ▀▀▀\033[0m"

// Spinner frames for animation
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Spinner represents an animated spinner with logo
type Spinner struct {
	message  string
	stopChan chan struct{}
	doneChan chan struct{}
	mu       sync.Mutex
	running  bool
}

// NewSpinner creates a new spinner with the given message
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message:  message,
		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	// Print logo
	fmt.Fprintln(os.Stderr, Logo)

	go func() {
		defer close(s.doneChan)
		i := 0
		for {
			select {
			case <-s.stopChan:
				// Clear the spinner line
				fmt.Fprintf(os.Stderr, "\r\033[K")
				return
			default:
				frame := spinnerFrames[i%len(spinnerFrames)]
				fmt.Fprintf(os.Stderr, "\r %s %s", frame, s.message)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the spinner animation
func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	close(s.stopChan)
	<-s.doneChan
}

// StopWithMessage stops the spinner and prints a final message
func (s *Spinner) StopWithMessage(message string) {
	s.Stop()
	fmt.Fprintf(os.Stderr, " ✓ %s\n", message)
}

// StopWithError stops the spinner and prints an error message
func (s *Spinner) StopWithError(message string) {
	s.Stop()
	fmt.Fprintf(os.Stderr, " ✗ %s\n", message)
}
