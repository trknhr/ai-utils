package model

import "time"

// GenerationResult captures a provider execution result in a unified shape.
type GenerationResult struct {
	Provider string        // Provider name (e.g., "claude")
	FilePath string        // Temp file path where the content is stored
	Content  string        // Generated raw text
	Duration time.Duration // Execution duration
	Err      error         // Execution error, if any
}
