package cli

import "testing"

func TestAppendLanguageHint_PrefersFlag(t *testing.T) {
	prompt := "Hello"
	got := appendLanguageHint(prompt, "ja", "en")
	want := "Hello\n\nPlease respond in ja."
	if got != want {
		t.Fatalf("expected flag lang to win.\n got: %q\nwant: %q", got, want)
	}
}

func TestAppendLanguageHint_ConfigWhenNoFlag(t *testing.T) {
	prompt := "Hello"
	got := appendLanguageHint(prompt, "", "fr")
	want := "Hello\n\nPlease respond in fr."
	if got != want {
		t.Fatalf("expected config lang.\n got: %q\nwant: %q", got, want)
	}
}

func TestAppendLanguageHint_NoLang(t *testing.T) {
	prompt := "Hello"
	got := appendLanguageHint(prompt, "", "")
	if got != prompt {
		t.Fatalf("expected prompt unchanged when no lang, got: %q", got)
	}
}
