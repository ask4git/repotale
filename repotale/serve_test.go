package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestLocalTmplRendersSessionID(t *testing.T) {
	data := struct {
		RepoID string
		Posts  []localPost
	}{
		RepoID: "my-project",
		Posts: []localPost{
			{
				Slug:        "abc1234",
				Title:       "테스트 포스트",
				Excerpt:     "세션ID 렌더 확인용",
				Tags:        []string{"test"},
				CommitSHA:   "abc1234567890",
				PublishedAt: "2026-08-02",
				SessionID:   "6c33497a-9e62-43c1-9620-15a44f7ed573",
			},
		},
	}

	var buf bytes.Buffer
	if err := localTmpl.Execute(&buf, data); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "6c33497a-9e62-43c1-9620-15a44f7ed573") {
		t.Errorf("rendered page missing session id, got:\n%s", out)
	}
	if !strings.Contains(out, "commit abc1234") {
		t.Errorf("rendered page missing truncated commit sha, got:\n%s", out)
	}
}

func TestLocalTmplOmitsSessionIDWhenEmpty(t *testing.T) {
	data := struct {
		RepoID string
		Posts  []localPost
	}{
		RepoID: "my-project",
		Posts: []localPost{
			{Slug: "abc1234", Title: "t", CommitSHA: "abc1234567890", PublishedAt: "2026-08-02"},
		},
	}

	var buf bytes.Buffer
	if err := localTmpl.Execute(&buf, data); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), "claude session") {
		t.Errorf("expected no session id text for a post with an empty SessionID")
	}
}
