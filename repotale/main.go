package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strings"
)

func die(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		cmdAnalyze()
		return
	}

	switch os.Args[1] {
	case "login":
		cmdLogin()
	case "connect":
		cmdConnect()
	case "open":
		cmdOpen()
	case "--version", "-v", "version":
		cmdVersion()
	case "update":
		cmdUpdate()
	case "-h", "--help", "help":
		usage()
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`usage: repotale [command]

with no command: analyze the local repo's recent commits (prompts for
repo path and requires the claude CLI, logged in, in PATH)

commands:
  login              log in with a GitHub personal access token
  connect <owner/repo>  connect a GitHub repo
  open               open the connected repo's web dashboard
  update             go install the latest version of repotale
  -h, --help         show this help
  -v, --version      show version`)
}

const modulePath = "github.com/ask4git/repotale/repotale"

func cmdUpdate() {
	if _, err := exec.LookPath("go"); err != nil {
		die(fmt.Errorf("go not found in PATH - needed to update (https://go.dev/dl)"))
	}

	fmt.Println("updating via go install", modulePath+"@latest ...")
	cmd := exec.Command("go", "install", modulePath+"@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		die(fmt.Errorf("update failed: %w", err))
	}
	fmt.Println("updated")
}

func cmdVersion() {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		fmt.Println("repotale (dev build)")
		return
	}
	fmt.Println("repotale", info.Main.Version)
}

func cmdLogin() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Print("Enter your GitHub personal access token: ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		// io.EOF with data already read just means no trailing newline
		// (e.g. `echo -n "$TOKEN" | repotale login`) — the token is still valid.
		if err != nil && err != io.EOF {
			die(fmt.Errorf("reading token: %w", err))
		}
		token = strings.TrimSpace(line)
	}

	if token == "" {
		die(fmt.Errorf("no token provided"))
	}

	if err := validateToken(token); err != nil {
		die(err)
	}

	cfg, err := loadConfig()
	if err != nil {
		die(err)
	}
	cfg.Token = token

	if err := saveConfig(cfg); err != nil {
		die(err)
	}

	fmt.Println("logged in successfully")
}

func cmdConnect() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: repotale connect <owner/repo>")
		os.Exit(1)
	}
	ownerRepo := os.Args[2]

	parts := strings.Split(ownerRepo, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		die(fmt.Errorf("repo must be in owner/repo form"))
	}

	cfg, err := loadConfig()
	if err != nil {
		die(err)
	}
	if cfg.Token == "" {
		die(fmt.Errorf("not logged in, run `repotale login` first"))
	}

	if err := validateRepo(cfg.Token, ownerRepo); err != nil {
		die(err)
	}

	cfg.Repo = ownerRepo
	if err := saveConfig(cfg); err != nil {
		die(err)
	}

	fmt.Println("connected to", ownerRepo)
}

func cmdOpen() {
	cfg, err := loadConfig()
	if err != nil {
		die(err)
	}
	if cfg.Repo == "" {
		die(fmt.Errorf("no repo connected, run `repotale connect <owner/repo>` first"))
	}

	url := "http://localhost:3000/repo/" + cfg.Repo
	if err := openBrowser(url); err != nil {
		die(err)
	}
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
