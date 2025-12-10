package provider

import "testing"

func TestFilterGeminiNoise(t *testing.T) {
	out := `[INFO] ok
[ERROR] [IDEClient] Directory mismatch.
real payload
[ERROR] [IDEClient] Another warning`

	got := filterGeminiNoise(out)
	want := "[INFO] ok\nreal payload"
	if got != want {
		t.Fatalf("unexpected filtered output.\n got:\n%q\nwant:\n%q", got, want)
	}
}
