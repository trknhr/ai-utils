package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/trknhr/ai-utils/internal/domain/model"
	"github.com/trknhr/ai-utils/internal/provider"
	"github.com/trknhr/ai-utils/internal/ui"
)

// ExecuteOptions bundles the knobs for executing prompts.
type ExecuteOptions struct {
	Parallel   bool
	Prompt     string
	Providers  []provider.Provider
	Timeout    time.Duration
	WorkingDir string
	Verbose    bool
	Model      string
}

// Run executes a prompt either with a single provider or in parallel.
func Run(ctx context.Context, opts ExecuteOptions) (string, error) {
	if len(opts.Providers) == 0 {
		return "", errors.New("no providers supplied")
	}

	if opts.Parallel && len(opts.Providers) > 1 {
		return runParallel(ctx, opts)
	}

	return runSingle(ctx, opts, opts.Providers[0])
}

func runSingle(ctx context.Context, opts ExecuteOptions, prov provider.Provider) (string, error) {
	spinner := ui.NewSpinner(fmt.Sprintf("Executing with %s...", prov.Name()))
	spinner.Start()

	content, err := prov.Execute(ctx, opts.Prompt, provider.ExecuteOptions{
		Timeout:    opts.Timeout,
		WorkingDir: opts.WorkingDir,
		Verbose:    opts.Verbose,
		Model:      opts.Model,
	})
	if err != nil {
		spinner.StopWithError("Execution failed")
		return "", err
	}

	spinner.StopWithMessage("Done!")
	return content, nil
}

func runParallel(ctx context.Context, opts ExecuteOptions) (string, error) {
	fmt.Fprintln(os.Stderr, ui.Logo)
	fmt.Fprintln(os.Stderr, "⚡ Running providers in parallel...")

	type indexedResult struct {
		idx int
		res model.GenerationResult
	}

	results := make([]model.GenerationResult, len(opts.Providers))
	resultCh := make(chan indexedResult, len(opts.Providers))
	progressLinesPrinted := 0

	var wg sync.WaitGroup
	for i, prov := range opts.Providers {
		wg.Add(1)
		go func(idx int, p provider.Provider) {
			defer wg.Done()

			start := time.Now()
			content, err := p.Execute(ctx, opts.Prompt, provider.ExecuteOptions{
				Timeout:    opts.Timeout,
				WorkingDir: opts.WorkingDir,
				Verbose:    opts.Verbose,
				Model:      opts.Model,
			})
			duration := time.Since(start)

			var filePath string
			if err == nil {
				filePath, err = writeResultFile(p.Name(), content)
			}

			resultCh <- indexedResult{idx: idx, res: model.GenerationResult{
				Provider: p.Name(),
				FilePath: filePath,
				Content:  content,
				Duration: duration,
				Err:      err,
			}}
		}(i, prov)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for item := range resultCh {
		results[item.idx] = item.res

		progressLinesPrinted++
		if item.res.Err != nil {
			fmt.Fprintf(os.Stderr, "✗ %s failed (%.1fs): %v\n", item.res.Provider, item.res.Duration.Seconds(), item.res.Err)
			continue
		}
		fmt.Fprintf(os.Stderr, "✓ %s done (%.1fs): %s\n", item.res.Provider, item.res.Duration.Seconds(), item.res.FilePath)
	}

	successes := make([]model.GenerationResult, 0, len(results))
	for _, res := range results {
		if res.Err != nil {
			// Failure already printed progressively above.
			continue
		}
		successes = append(successes, res)
	}

	if len(successes) == 0 {
		return "", fmt.Errorf("all providers failed")
	}

	filePaths := make([]string, 0, len(successes))
	for _, res := range successes {
		filePaths = append(filePaths, res.FilePath)
	}
	openInVSCode(filePaths)

	clearLastTerminalLines(os.Stderr, progressLinesPrinted)

	return ui.SelectResult(successes)
}

func clearLastTerminalLines(f *os.File, lines int) {
	if lines <= 0 || !isCharDevice(f) {
		return
	}

	// Move up and clear each line. This keeps earlier logs (e.g. header/logo) intact.
	for i := 0; i < lines; i++ {
		fmt.Fprint(f, "\x1b[1A\x1b[2K\r")
	}
}

func isCharDevice(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
