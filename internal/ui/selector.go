package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/trknhr/ai-utils/internal/domain/model"
)

// SelectResult renders a TUI to pick a generation result.
func SelectResult(results []model.GenerationResult) (string, error) {
	if len(results) == 0 {
		return "", fmt.Errorf("no results available for selection")
	}

	options := make([]huh.Option[string], 0, len(results))
	for _, res := range results {
		label := fmt.Sprintf("%s (%.1fs) - %s", res.Provider, res.Duration.Seconds(), res.FilePath)
		options = append(options, huh.NewOption(label, res.Content))
	}

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("✨ Select the best output").
				Description("Comparison opened in VS Code").
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return "", err
	}

	return selected, nil
}
