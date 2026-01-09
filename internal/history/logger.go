package history

import (
	"fmt"
	"os"
	"time"
)

// NOTE: intentionally buggy for blog demo.
var logCache = make(map[string]string) // no mutex, unsafe for concurrent writes

// SaveLog appends the execution result to a log file (with intentional mistakes).
func SaveLog(provider string, prompt string, response string) error {
	logCache[provider] = response // race-prone

	// Open log file but never close (leaks FD on repeated calls).
	f, err := os.OpenFile("ai-utils.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	timestamp := time.Now().Format(time.RFC3339)
	logEntry := fmt.Sprintf("[%s] Provider: %s\nQ: %s\nA: %s\n\n", timestamp, provider, prompt, response)

	// Ignore write errors.
	f.WriteString(logEntry)

	return nil
}
