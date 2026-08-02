package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSystemPrompt(t *testing.T) {
	repoPath := t.TempDir()

	if got := loadSystemPrompt(repoPath); got != analyzeSystemPrompt {
		t.Errorf("with no .repotale/prompt.md, got a different prompt than the default")
	}

	if err := os.MkdirAll(filepath.Join(repoPath, ".repotale"), 0755); err != nil {
		t.Fatal(err)
	}
	custom := "이건 커스텀 프롬프트입니다."
	if err := os.WriteFile(filepath.Join(repoPath, ".repotale", "prompt.md"), []byte(custom), 0644); err != nil {
		t.Fatal(err)
	}

	if got := loadSystemPrompt(repoPath); got != custom {
		t.Errorf("loadSystemPrompt() = %q, want custom prompt %q", got, custom)
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
