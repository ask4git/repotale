package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const analyzeJSONSchema = `{"type":"object","properties":{"title":{"type":"string"},"excerpt":{"type":"string"},"tags":{"type":"array","items":{"type":"string"}}},"required":["title","excerpt","tags"],"additionalProperties":false}`

// mirrors web/lib/pipeline/generate-post.ts's SYSTEM_PROMPT, adapted for a
// single commit (subject/body/diff) instead of a PR (title/body/diff).
const analyzeSystemPrompt = `당신은 이 repository의 커밋을 비개발자도 읽을 수 있는 개발 블로그 글로 바꾸는 역할입니다.

입력: 커밋 메시지, 전체 diff

먼저 아래 네 가지 관점에서 변경사항을 분석하세요:
- What: 어떤 기능이 생겼나 - 사용자 관점 언어로, 기술 용어 최소화
- Where: 어느 파일/모듈에 들어갔나 - 아키텍처 위치를 비유로 설명
- How: 어떤 방식으로 구현했나 - 기술 선택을 평이한 언어로, "왜 이 방법인지"까지
- Why: 판단 근거 - diff와 커밋 메시지에서 추론한 의사결정 서사. 없는 근거를 지어내지 말고, 원문에 있는 이유만 사용

톤: 기술 문서가 아니라 개발 블로그 글이어야 합니다. "오늘 ~에 ~를 붙였습니다. ~때문인데요" 수준의 문장.

이 분석을 바탕으로 다음 필드를 채워 응답하세요:
- title: 블로그 포스트 제목 (한국어, 한 문장)
- excerpt: What을 중심으로 한 2~4문장 요약 (위 톤을 유지)
- tags: 1~3개의 짧은 영문 소문자 태그 (예: "feature", "bugfix")

hallucination 방지: 커밋 메시지와 diff에 없는 사실을 지어내지 마세요.`

type localPost struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Excerpt     string   `json:"excerpt"`
	Tags        []string `json:"tags"`
	CommitSHA   string   `json:"commitSha"`
	PublishedAt string   `json:"publishedAt"`
	SessionID   string   `json:"sessionId,omitempty"`
}

type generatedFields struct {
	Title   string   `json:"title"`
	Excerpt string   `json:"excerpt"`
	Tags    []string `json:"tags"`
}

func cmdAnalyze() {
	repoPath, err := filepath.Abs(promptRepoPath())
	if err != nil {
		die(err)
	}
	if !isGitRepo(repoPath) {
		die(fmt.Errorf("%s is not a git repository", repoPath))
	}

	// ponytail: only one model adapter for now (Claude Code, via the locally
	// authenticated `claude` CLI) - add a picker prompt when a second one lands.
	if _, err := exec.LookPath("claude"); err != nil {
		die(fmt.Errorf("claude CLI not found in PATH - install Claude Code first (https://claude.com/claude-code); it's the only supported model right now"))
	}
	fmt.Println("model: Claude Code (claude CLI, using your local login)")

	systemPrompt := loadSystemPrompt(repoPath)

	shas, err := recentCommits(repoPath, 10)
	if err != nil {
		die(fmt.Errorf("reading git log: %w", err))
	}
	if len(shas) == 0 {
		fmt.Println("no commits found")
		return
	}

	repoID := repoIdentifier(repoPath)
	outDir := filepath.Join(repotaleDir(), repoID, "posts")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		die(err)
	}

	fmt.Printf("analyzing %d commit(s) in %s...\n", len(shas), repoPath)
	saved := 0
	for _, sha := range shas {
		postPath := filepath.Join(outDir, sha[:7]+".json")
		if _, err := os.Stat(postPath); err == nil {
			fmt.Printf("  %s (already analyzed, skipping)\n", sha[:7])
			continue
		}

		post, err := analyzeCommit(repoPath, sha, systemPrompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  skip %s: %v\n", sha[:7], err)
			continue
		}
		data, err := json.MarshalIndent(post, "", "  ")
		if err != nil {
			die(err)
		}
		if err := os.WriteFile(postPath, data, 0644); err != nil {
			die(err)
		}
		fmt.Printf("  %s %s  (claude session %s)\n", sha[:7], post.Title, post.SessionID)
		saved++
	}

	posts, err := loadLocalPosts(outDir)
	if err != nil {
		die(err)
	}
	if len(posts) == 0 {
		fmt.Println("\nno posts to show")
		return
	}

	fmt.Printf("\n%d new, %d total post(s) in %s\n", saved, len(posts), outDir)
	if err := serveLocal(repoID, posts); err != nil {
		die(err)
	}
}

func promptRepoPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		die(err)
	}
	fmt.Printf("분석할 repo 경로 [%s]: ", cwd)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return cwd
	}
	return line
}

func isGitRepo(path string) bool {
	return exec.Command("git", "-C", path, "rev-parse", "--is-inside-work-tree").Run() == nil
}

func recentCommits(repoPath string, n int) ([]string, error) {
	out, err := exec.Command("git", "-C", repoPath, "log", fmt.Sprintf("-%d", n), "--format=%H").Output()
	if err != nil {
		return nil, err
	}
	var shas []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line != "" {
			shas = append(shas, line)
		}
	}
	return shas, nil
}

func repotaleDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		die(err)
	}
	return filepath.Join(home, ".repotale")
}

func repoIdentifier(repoPath string) string {
	name := strings.ToLower(filepath.Base(repoPath))
	var b strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	id := strings.Trim(b.String(), "-")
	if id == "" {
		id = "repo"
	}
	return id
}

// loadSystemPrompt lets a repo fully override the analysis prompt via
// .repotale/prompt.md (mirrors CLAUDE.md: a project-checked-in file that
// customizes agent behavior). Falls back to the built-in default.
func loadSystemPrompt(repoPath string) string {
	custom, err := os.ReadFile(filepath.Join(repoPath, ".repotale", "prompt.md"))
	if err != nil {
		return analyzeSystemPrompt
	}
	return string(custom)
}

func analyzeCommit(repoPath, sha, systemPrompt string) (localPost, error) {
	message, err := commitMessage(repoPath, sha)
	if err != nil {
		return localPost{}, fmt.Errorf("reading commit message: %w", err)
	}
	diff, err := commitDiff(repoPath, sha)
	if err != nil {
		return localPost{}, fmt.Errorf("reading commit diff: %w", err)
	}
	date, err := commitDate(repoPath, sha)
	if err != nil {
		return localPost{}, fmt.Errorf("reading commit date: %w", err)
	}

	// keep the diff bounded - an oversized commit shouldn't blow the prompt budget
	const maxDiffLen = 20000
	if len(diff) > maxDiffLen {
		diff = diff[:maxDiffLen] + "\n...(truncated)"
	}

	prompt := systemPrompt + "\n\n커밋 메시지:\n" + message + "\n\ndiff:\n" + diff
	fields, sessionID, err := runClaude(prompt)
	if err != nil {
		return localPost{}, err
	}

	publishedAt := date
	if len(date) >= 10 {
		publishedAt = date[:10]
	}

	return localPost{
		Slug:        sha[:7],
		Title:       fields.Title,
		Excerpt:     fields.Excerpt,
		Tags:        fields.Tags,
		CommitSHA:   sha,
		PublishedAt: publishedAt,
		SessionID:   sessionID,
	}, nil
}

func commitMessage(repoPath, sha string) (string, error) {
	out, err := exec.Command("git", "-C", repoPath, "log", "-1", "--format=%B", sha).Output()
	return strings.TrimSpace(string(out)), err
}

func commitDate(repoPath, sha string) (string, error) {
	out, err := exec.Command("git", "-C", repoPath, "log", "-1", "--format=%aI", sha).Output()
	return strings.TrimSpace(string(out)), err
}

func commitDiff(repoPath, sha string) (string, error) {
	out, err := exec.Command("git", "-C", repoPath, "show", sha, "--format=").Output()
	return string(out), err
}

// runClaude shells out to the locally-authenticated Claude Code CLI in
// non-interactive mode. --allowedTools "" keeps it a pure text-in/JSON-out
// call (no filesystem exploration), which is both cheaper and deterministic.
func runClaude(prompt string) (generatedFields, string, error) {
	cmd := exec.Command("claude", "-p", prompt,
		"--output-format", "json",
		"--json-schema", analyzeJSONSchema,
		"--allowedTools", "",
	)
	out, err := cmd.Output()
	if err != nil {
		return generatedFields{}, "", fmt.Errorf("claude CLI: %w", err)
	}

	var wrapper struct {
		IsError          bool            `json:"is_error"`
		SessionID        string          `json:"session_id"`
		StructuredOutput generatedFields `json:"structured_output"`
	}
	if err := json.Unmarshal(out, &wrapper); err != nil {
		return generatedFields{}, "", fmt.Errorf("parsing claude output: %w", err)
	}
	if wrapper.IsError {
		return generatedFields{}, wrapper.SessionID, fmt.Errorf("claude returned an error")
	}
	return wrapper.StructuredOutput, wrapper.SessionID, nil
}
