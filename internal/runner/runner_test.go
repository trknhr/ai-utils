package runner

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/trknhr/ai-utils/internal/provider"
)

type mockProvider struct {
	name      string
	content   string
	delay     time.Duration
	execErr   error
	available bool
}

func (m mockProvider) Name() string {
	return m.name
}

func (m mockProvider) Available() bool {
	return m.available
}

func (m mockProvider) Execute(ctx context.Context, prompt string, opts provider.ExecuteOptions) (string, error) {
	if m.delay > 0 {
		select {
		case <-time.After(m.delay):
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	if m.execErr != nil {
		return "", m.execErr
	}
	return m.content + ":" + prompt, nil
}

func TestRunSingle(t *testing.T) {
	mp := mockProvider{
		name:      "mock",
		content:   "ok",
		available: true,
	}

	result, err := runSingle(context.Background(), ExecuteOptions{
		Prompt: "hello",
	}, mp)
	if err != nil {
		t.Fatalf("runSingle returned error: %v", err)
	}
	if result != "ok:hello" {
		t.Fatalf("unexpected result: %q", result)
	}
}

func TestRunParallelAllFail(t *testing.T) {
	mp := mockProvider{
		name:      "mock",
		execErr:   errors.New("boom"),
		available: true,
	}

	_, err := runParallel(context.Background(), ExecuteOptions{
		Prompt:    "hello",
		Providers: []provider.Provider{mp},
	})
	if err == nil {
		t.Fatalf("expected error when all providers fail")
	}
}
