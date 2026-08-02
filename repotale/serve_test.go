package main

import (
	"bytes"
	"strings"
	"testing"
)

type localPageData struct {
	Lang         string
	RepoID       string
	Posts        []localPost
	Heading      string
	Empty        string
	CommitLabel  string
	SessionLabel string
}

func renderLocalPage(t *testing.T, posts []localPost) string {
	t.Helper()
	data := localPageData{
		Lang:         "en",
		RepoID:       "my-project",
		Posts:        posts,
		Heading:      "What's been happening in this repo",
		Empty:        "No analysis yet.",
		CommitLabel:  "commit",
		SessionLabel: "claude session",
	}
	var buf bytes.Buffer
	if err := localTmpl.Execute(&buf, data); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

func TestLocalTmplRendersSessionID(t *testing.T) {
	out := renderLocalPage(t, []localPost{
		{
			Slug:        "abc1234",
			Title:       "테스트 포스트",
			Excerpt:     "세션ID 렌더 확인용",
			Tags:        []string{"test"},
			CommitSHA:   "abc1234567890",
			PublishedAt: "2026-08-02",
			SessionID:   "6c33497a-9e62-43c1-9620-15a44f7ed573",
		},
	})

	if !strings.Contains(out, "6c33497a-9e62-43c1-9620-15a44f7ed573") {
		t.Errorf("rendered page missing session id, got:\n%s", out)
	}
	if !strings.Contains(out, "commit abc1234") {
		t.Errorf("rendered page missing truncated commit sha, got:\n%s", out)
	}
}

func TestLocalTmplOmitsSessionIDWhenEmpty(t *testing.T) {
	out := renderLocalPage(t, []localPost{
		{Slug: "abc1234", Title: "t", CommitSHA: "abc1234567890", PublishedAt: "2026-08-02"},
	})
	if strings.Contains(out, "claude session") {
		t.Errorf("expected no session id text for a post with an empty SessionID")
	}
}
