package main

import "testing"

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
