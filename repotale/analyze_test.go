package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSystemPrompt_RepoOverrideWinsAndSkipsPrompt(t *testing.T) {
	repoPath := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repoPath, ".repotale"), 0755); err != nil {
		t.Fatal(err)
	}
	custom := "이건 커스텀 프롬프트입니다."
	if err := os.WriteFile(filepath.Join(repoPath, ".repotale", "prompt.md"), []byte(custom), 0644); err != nil {
		t.Fatal(err)
	}

	// empty reader: if the override didn't short-circuit, promptTone would
	// block reading from it and this test would hang instead of failing fast.
	reader := bufio.NewReader(strings.NewReader(""))
	got, err := loadSystemPrompt(repoPath, reader)
	if err != nil {
		t.Fatal(err)
	}
	if got != custom {
		t.Errorf("loadSystemPrompt() = %q, want custom prompt %q", got, custom)
	}
}

func TestLoadSystemPrompt_NoOverrideUsesPresetChoice(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	repoPath := t.TempDir()

	reader := bufio.NewReader(strings.NewReader("hard\nen\n"))
	got, err := loadSystemPrompt(repoPath, reader)
	if err != nil {
		t.Fatal(err)
	}

	preset, err := loadPreset("hard")
	if err != nil {
		t.Fatal(err)
	}
	want := preset + analysisLanguageDirective(langEN)
	if got != want {
		t.Errorf("loadSystemPrompt() did not return the hard preset + the chosen output-language directive")
	}
}

func TestPromptTone(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if err := ensurePresetsOnDisk(); err != nil {
		t.Fatal(err)
	}

	cases := map[string]string{
		"soft\n":      "soft",
		"hard\n":      "hard",
		"\n":          defaultPreset, // empty input -> default
		"bogus\n":     defaultPreset, // unknown name -> falls back to default
		"  medium \n": "medium",      // surrounding whitespace trimmed
	}
	for input, want := range cases {
		got := promptTone(bufio.NewReader(strings.NewReader(input)))
		if got != want {
			t.Errorf("promptTone(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestEnsurePresetsOnDiskDoesNotClobberEdits(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if err := ensurePresetsOnDisk(); err != nil {
		t.Fatal(err)
	}
	edited := "내가 고친 프리셋"
	path := filepath.Join(presetsDir(), "soft.md")
	if err := os.WriteFile(path, []byte(edited), 0644); err != nil {
		t.Fatal(err)
	}

	if err := ensurePresetsOnDisk(); err != nil {
		t.Fatal(err)
	}
	got, err := loadPreset("soft")
	if err != nil {
		t.Fatal(err)
	}
	if got != edited {
		t.Errorf("ensurePresetsOnDisk() overwrote a user-edited preset")
	}
}

func TestNormalizeGitHubURL(t *testing.T) {
	cases := map[string]string{
		"git@github.com:ask4git/repotale.git":     "https://github.com/ask4git/repotale",
		"git@github.com:ask4git/repotale":         "https://github.com/ask4git/repotale",
		"https://github.com/ask4git/repotale.git": "https://github.com/ask4git/repotale",
		"https://github.com/ask4git/repotale":     "https://github.com/ask4git/repotale",
		"git@gitlab.com:someone/somerepo.git":     "",
		"":                                        "",
	}
	for input, want := range cases {
		if got := normalizeGitHubURL(input); got != want {
			t.Errorf("normalizeGitHubURL(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestRepoIdentifier(t *testing.T) {
	cases := map[string]string{
		"/Users/me/My Project":  "my-project",
		"/Users/me/repotale":    "repotale",
		"/Users/me/foo.bar_baz": "foo-bar-baz",
		"/":                     "repo",
	}
	for path, want := range cases {
		if got := repoIdentifier(path); got != want {
			t.Errorf("repoIdentifier(%q) = %q, want %q", path, got, want)
		}
	}
}
