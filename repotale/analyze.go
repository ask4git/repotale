package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const analyzeJSONSchema = `{"type":"object","properties":{"title":{"type":"string"},"excerpt":{"type":"string"},"tags":{"type":"array","items":{"type":"string"}}},"required":["title","excerpt","tags"],"additionalProperties":false}`

//go:embed presets/*.md
var presetsFS embed.FS

// spectrum from narrative/warm to dry/formal. ponytail: add "xsoft"/"xhard"
// here + presets/xsoft.md, presets/xhard.md when that's actually needed.
var presetNames = []string{"soft", "medium", "hard"}

const defaultPreset = "medium"

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

func cmdAnalyze(stdin *bufio.Reader) {
	repoPath, err := filepath.Abs(promptRepoPath(stdin))
	if err != nil {
		die(err)
	}
	if !isGitRepo(repoPath) {
		die(fmt.Errorf(t("not_git_repo"), repoPath))
	}

	// ponytail: only one model adapter for now (Claude Code, via the locally
	// authenticated `claude` CLI) - add a picker prompt when a second one lands.
	if _, err := exec.LookPath("claude"); err != nil {
		die(fmt.Errorf("%s", t("claude_not_found")))
	}
	fmt.Println(t("analyze_model_line"))
	if !confirm(stdin, t("analyze_confirm_claude"), true) {
		fmt.Println(t("cancelled"))
		return
	}

	shas, err := recentCommits(repoPath, 10)
	if err != nil {
		die(fmt.Errorf("reading git log: %w", err))
	}
	if len(shas) == 0 {
		fmt.Println(t("no_commits"))
		return
	}

	repoID := repoIdentifier(repoPath)
	outDir := filepath.Join(repotaleDir(), repoID, "posts")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		die(err)
	}

	newShas, estTokens, err := planAnalysis(repoPath, outDir, shas)
	if err != nil {
		die(err)
	}

	var systemPrompt string
	if len(newShas) == 0 {
		fmt.Println(t("no_new_commits"))
	} else {
		fmt.Println(t("token_estimate", len(newShas), estTokens))
		if !confirm(stdin, t("confirm_proceed"), true) {
			fmt.Println(t("cancelled"))
			return
		}

		systemPrompt, err = loadSystemPrompt(repoPath, stdin)
		if err != nil {
			die(err)
		}
	}

	fmt.Println(t("analyzing_line", len(shas), repoPath))
	saved := 0
	for _, sha := range shas {
		postPath := filepath.Join(outDir, sha[:7]+".json")
		if _, err := os.Stat(postPath); err == nil {
			fmt.Println(t("already_analyzed", sha[:7]))
			continue
		}

		post, err := analyzeCommit(repoPath, sha, systemPrompt)
		if err != nil {
			fmt.Fprintln(os.Stderr, t("skip_error", sha[:7], err))
			continue
		}
		data, err := json.MarshalIndent(post, "", "  ")
		if err != nil {
			die(err)
		}
		if err := os.WriteFile(postPath, data, 0644); err != nil {
			die(err)
		}
		fmt.Println(t("post_line", sha[:7], post.Title, post.SessionID))
		saved++
	}

	posts, err := loadLocalPosts(outDir)
	if err != nil {
		die(err)
	}
	if len(posts) == 0 {
		fmt.Println()
		fmt.Println(t("no_posts"))
		return
	}

	fmt.Println()
	fmt.Println(t("summary_line", saved, len(posts), outDir))
	if err := serveLocal(repoID, posts); err != nil {
		die(err)
	}
}

func confirm(reader *bufio.Reader, question string, defaultYes bool) bool {
	suffix := "[y/N]"
	if defaultYes {
		suffix = "[Y/n]"
	}
	fmt.Printf("%s (%s): ", question, suffix)
	line, _ := reader.ReadString('\n')
	answer := strings.ToLower(strings.TrimSpace(line))
	if answer == "" {
		return defaultYes
	}
	return answer == "y" || answer == "yes"
}

// planAnalysis figures out which commits still need analysis and gives a
// rough token estimate (~4 chars/token) from their diff size, without
// calling claude - so the cost is visible before it's spent.
func planAnalysis(repoPath, outDir string, shas []string) (newShas []string, estTokens int, err error) {
	for _, sha := range shas {
		postPath := filepath.Join(outDir, sha[:7]+".json")
		if _, err := os.Stat(postPath); err == nil {
			continue
		}
		diff, err := commitDiff(repoPath, sha)
		if err != nil {
			return nil, 0, fmt.Errorf("reading commit diff: %w", err)
		}
		newShas = append(newShas, sha)
		estTokens += len(diff) / 4
	}
	return newShas, estTokens, nil
}

func promptRepoPath(reader *bufio.Reader) string {
	cwd, err := os.Getwd()
	if err != nil {
		die(err)
	}
	fmt.Print(t("repo_path_prompt", cwd))
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

// loadSystemPrompt resolves the analysis prompt with this precedence:
//  1. .repotale/prompt.md in the target repo (full override, mirrors
//     CLAUDE.md - a project-checked-in file that customizes agent behavior).
//     When present, the tone/output-language pickers are skipped entirely -
//     full override means full control, no prompts, no injected directives.
//  2. a tone preset the user picks interactively, backed by an editable file
//     under ~/.repotale/presets/, plus an explicit output-language directive
//     (independent of the tone - and of the CLI's own UI language).
func loadSystemPrompt(repoPath string, reader *bufio.Reader) (string, error) {
	if custom, err := os.ReadFile(filepath.Join(repoPath, ".repotale", "prompt.md")); err == nil {
		return string(custom), nil
	}

	if err := ensurePresetsOnDisk(); err != nil {
		return "", fmt.Errorf("setting up presets: %w", err)
	}
	preset, err := loadPreset(promptTone(reader))
	if err != nil {
		return "", err
	}

	outputLang := promptAnalysisLanguage(reader)
	return preset + analysisLanguageDirective(outputLang), nil
}

func presetsDir() string {
	return filepath.Join(repotaleDir(), "presets")
}

// ensurePresetsOnDisk writes the built-in presets to ~/.repotale/presets/
// the first time, so they're a plain editable file - same pattern as
// CLAUDE.md. Never overwrites a file that's already there, in case the user
// edited their own copy.
func ensurePresetsOnDisk() error {
	dir := presetsDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	for _, name := range presetNames {
		path := filepath.Join(dir, name+".md")
		if _, err := os.Stat(path); err == nil {
			continue
		}
		content, err := presetsFS.ReadFile("presets/" + name + ".md")
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, content, 0644); err != nil {
			return err
		}
	}
	return nil
}

func loadPreset(name string) (string, error) {
	data, err := os.ReadFile(filepath.Join(presetsDir(), name+".md"))
	if err != nil {
		return "", fmt.Errorf("reading preset %q: %w", name, err)
	}
	return string(data), nil
}

func promptTone(reader *bufio.Reader) string {
	fmt.Print(t("tone_prompt", strings.Join(presetNames, "/"), defaultPreset))
	line, _ := reader.ReadString('\n')
	tone := strings.TrimSpace(line)
	if tone == "" {
		return defaultPreset
	}
	for _, name := range presetNames {
		if tone == name {
			return tone
		}
	}
	fmt.Println(t("tone_unknown", tone, defaultPreset))
	return defaultPreset
}

// promptAnalysisLanguage asks what language the generated posts should be
// written in - independent of currentLang (the CLI's own UI language), but
// defaulting to it since that's usually what's wanted.
func promptAnalysisLanguage(reader *bufio.Reader) lang {
	fmt.Print(t("analysis_language_prompt", currentLang))
	line, _ := reader.ReadString('\n')
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "en":
		return langEN
	case "ko":
		return langKO
	default:
		return currentLang
	}
}

func analysisLanguageDirective(l lang) string {
	if l == langEN {
		return "\n\nOutput language: write the title and excerpt entirely in English."
	}
	return "\n\n출력 언어: title과 excerpt를 반드시 한국어로 작성하세요."
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

	prompt := systemPrompt + "\n\ncommit message:\n" + message + "\n\ndiff:\n" + diff
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
