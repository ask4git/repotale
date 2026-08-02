package main

import (
	"bufio"
	"fmt"
	"strings"
)

type lang string

const (
	langEN lang = "en"
	langKO lang = "ko"
)

// currentLang is process-wide, set once during startup (initLanguage) before
// any other command runs. It controls the CLI's own UI text - separate from
// the analysis output language (see analysisLanguageDirective in analyze.go).
var currentLang = langEN

const usageEN = `usage: repotale [command]

with no command: analyze the local repo's recent commits (prompts for
repo path and requires the claude CLI, logged in, in PATH)

commands:
  login              log in with a GitHub personal access token
  connect <owner/repo>  connect a GitHub repo
  open               open the connected repo's web dashboard
  update             go install the latest version of repotale
  --language         change the UI language
  -h, --help         show this help
  -v, --version      show version`

const usageKO = `usage: repotale [command]

명령 없이 실행: 로컬 repo의 최근 커밋을 분석 (repo 경로를 물어보고,
claude CLI가 로그인된 상태로 PATH에 있어야 함)

명령:
  login              GitHub personal access token으로 로그인
  connect <owner/repo>  GitHub repo 연결
  open               연결된 repo의 웹 대시보드 열기
  update             repotale 최신 버전으로 업데이트
  --language         UI 언어 변경
  -h, --help         도움말 표시
  -v, --version      버전 표시`

var enMessages = map[string]string{
	"select_language":   "Select language / 언어를 선택하세요:\n  [1] English\n  [2] 한국어\nChoice / 선택 [1]: ",
	"language_saved":    "Language set to %s",
	"error_prefix":      "error:",
	"usage":             usageEN,
	"version_dev":       "repotale (dev build)",
	"update_start":      "updating via go install %s@latest ...",
	"update_done":       "updated",
	"update_go_missing": "go not found in PATH - needed to update (https://go.dev/dl)",
	"update_failed":     "update failed: %w",

	"login_prompt":          "Enter your GitHub personal access token: ",
	"login_reading_token":   "reading token: %w",
	"login_no_token":        "no token provided",
	"login_success":         "logged in successfully",
	"connect_usage":         "usage: repotale connect <owner/repo>",
	"connect_bad_format":    "repo must be in owner/repo form",
	"connect_not_logged_in": "not logged in, run `repotale login` first",
	"connect_success":       "connected to %s",
	"open_no_repo":          "no repo connected, run `repotale connect <owner/repo>` first",

	"analyze_model_line":       "model: Claude Code (claude CLI, using your local login)",
	"analyze_confirm_claude":   "Proceed with analysis using this claude CLI?",
	"cancelled":                "Cancelled",
	"claude_not_found":         "claude CLI not found in PATH - install Claude Code first (https://claude.com/claude-code); it's the only supported model right now",
	"not_git_repo":             "%s is not a git repository",
	"no_commits":               "no commits found",
	"no_new_commits":           "No commits to analyze (already up to date)",
	"token_estimate":           "%d new commit(s), roughly %d input tokens (diff-based estimate, actual may differ)",
	"confirm_proceed":          "Proceed?",
	"analyzing_line":           "analyzing %d commit(s) in %s...",
	"already_analyzed":         "  %s (already analyzed, skipping)",
	"skip_error":               "  skip %s: %v",
	"post_line":                "  %s %s  (claude session %s)",
	"no_posts":                 "no posts to show",
	"summary_line":             "%d new, %d total post(s) in %s",
	"repo_path_prompt":         "Repo path to analyze [%s]: ",
	"tone_prompt":              "Analysis tone (%s) [%s]: ",
	"tone_unknown":             "Unknown tone %q, using default %s",
	"analysis_language_prompt": "Output language for the generated posts (en/ko) [%s]: ",

	"serving_line": "serving at %s - press Ctrl+C to stop",

	"web_heading":       "What's been happening in this repo",
	"web_empty":         "No analysis yet.",
	"web_commit_label":  "commit",
	"web_session_label": "claude session",
}

var koMessages = map[string]string{
	"select_language":   "언어를 선택하세요 / Select language:\n  [1] English\n  [2] 한국어\n선택 / Choice [1]: ",
	"language_saved":    "언어가 %s로 설정되었습니다",
	"error_prefix":      "에러:",
	"usage":             usageKO,
	"version_dev":       "repotale (개발 빌드)",
	"update_start":      "go install %s@latest 로 업데이트합니다...",
	"update_done":       "업데이트했습니다",
	"update_go_missing": "PATH에 go가 없습니다 - 업데이트하려면 필요합니다 (https://go.dev/dl)",
	"update_failed":     "업데이트 실패: %w",

	"login_prompt":          "GitHub personal access token을 입력하세요: ",
	"login_reading_token":   "토큰 읽기 실패: %w",
	"login_no_token":        "토큰이 입력되지 않았습니다",
	"login_success":         "로그인 성공",
	"connect_usage":         "usage: repotale connect <owner/repo>",
	"connect_bad_format":    "repo는 owner/repo 형식이어야 합니다",
	"connect_not_logged_in": "로그인이 안 되어 있습니다, 먼저 `repotale login`을 실행하세요",
	"connect_success":       "%s 에 연결했습니다",
	"open_no_repo":          "연결된 repo가 없습니다, 먼저 `repotale connect <owner/repo>`를 실행하세요",

	"analyze_model_line":       "model: Claude Code (claude CLI, 로컬 로그인 사용)",
	"analyze_confirm_claude":   "이 claude CLI로 분석을 진행할까요?",
	"cancelled":                "취소했습니다",
	"claude_not_found":         "claude CLI를 PATH에서 찾을 수 없습니다 - Claude Code를 먼저 설치하세요 (https://claude.com/claude-code); 지금은 이 모델만 지원합니다",
	"not_git_repo":             "%s 는 git repository가 아닙니다",
	"no_commits":               "커밋을 찾을 수 없습니다",
	"no_new_commits":           "새로 분석할 커밋이 없습니다 (이미 다 분석됨)",
	"token_estimate":           "새로 분석할 커밋 %d개, diff 기준 예상 입력 토큰 약 %d개 (대략치, 실제와 다를 수 있음)",
	"confirm_proceed":          "진행할까요?",
	"analyzing_line":           "%d개 커밋 분석 중 (%s)...",
	"already_analyzed":         "  %s (이미 분석됨, 건너뜀)",
	"skip_error":               "  %s 건너뜀: %v",
	"post_line":                "  %s %s  (claude 세션 %s)",
	"no_posts":                 "표시할 포스트가 없습니다",
	"summary_line":             "%d개 새로 생성, 총 %d개 포스트, 저장 위치: %s",
	"repo_path_prompt":         "분석할 repo 경로 [%s]: ",
	"tone_prompt":              "분석 톤 (%s) [%s]: ",
	"tone_unknown":             "모르는 톤 %q, 기본값 %s로 진행합니다",
	"analysis_language_prompt": "생성할 글의 언어 (en/ko) [%s]: ",

	"serving_line": "%s 에서 서비스 중 - 종료하려면 Ctrl+C",

	"web_heading":       "이 repo에서 무슨 일이 있었나",
	"web_empty":         "아직 분석 결과가 없습니다.",
	"web_commit_label":  "commit",
	"web_session_label": "claude 세션",
}

func t(key string, args ...any) string {
	m := enMessages
	if currentLang == langKO {
		m = koMessages
	}
	tmpl, ok := m[key]
	if !ok {
		tmpl = enMessages[key]
	}
	if len(args) == 0 {
		return tmpl
	}
	return fmt.Sprintf(tmpl, args...)
}

// initLanguage loads the saved language setting from ~/.repotale/settings,
// prompting to choose one (bilingual, since the language isn't known yet)
// the first time there's no saved choice.
func initLanguage(reader *bufio.Reader) {
	settings, err := loadSettings()
	if err == nil {
		if l := settings["language"]; l == string(langEN) || l == string(langKO) {
			currentLang = lang(l)
			return
		}
	}
	currentLang = promptLanguage(reader)
	_ = setSetting("language", string(currentLang))
}

func promptLanguage(reader *bufio.Reader) lang {
	fmt.Print(t("select_language"))
	line, _ := reader.ReadString('\n')
	if strings.TrimSpace(line) == "2" {
		return langKO
	}
	return langEN
}

// cmdLanguage lets the user change the saved language at any time
// (`repotale --language` / `repotale language`).
func cmdLanguage(reader *bufio.Reader) {
	chosen := promptLanguage(reader)
	if err := setSetting("language", string(chosen)); err != nil {
		die(err)
	}
	currentLang = chosen
	fmt.Println(t("language_saved", chosen))
}
