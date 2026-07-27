package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "login":
		cmdLogin()
	case "connect":
		cmdConnect()
	case "open":
		cmdOpen()
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("usage: repotale <login|connect <owner/repo>|open>")
}

func cmdLogin() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Print("Enter your GitHub personal access token: ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading token:", err)
			os.Exit(1)
		}
		token = strings.TrimSpace(line)
	}

	if token == "" {
		fmt.Fprintln(os.Stderr, "error: no token provided")
		os.Exit(1)
	}

	if err := validateToken(token); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	cfg.Token = token

	if err := saveConfig(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
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
		fmt.Fprintln(os.Stderr, "error: repo must be in owner/repo form")
		os.Exit(1)
	}

	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if cfg.Token == "" {
		fmt.Fprintln(os.Stderr, "error: not logged in, run `repotale login` first")
		os.Exit(1)
	}

	if err := validateRepo(cfg.Token, ownerRepo); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	cfg.Repo = ownerRepo
	if err := saveConfig(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fmt.Println("connected to", ownerRepo)
}

func cmdOpen() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if cfg.Repo == "" {
		fmt.Fprintln(os.Stderr, "error: no repo connected, run `repotale connect <owner/repo>` first")
		os.Exit(1)
	}

	url := "http://localhost:3000/repo/" + cfg.Repo
	if err := openBrowser(url); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
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
