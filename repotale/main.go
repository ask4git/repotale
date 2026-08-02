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
	fmt.Fprintln(os.Stderr, t("error_prefix"), err)
	os.Exit(1)
}

func main() {
	stdin := bufio.NewReader(os.Stdin)

	if len(os.Args) >= 2 && (os.Args[1] == "--language" || os.Args[1] == "language") {
		cmdLanguage(stdin)
		return
	}
	initLanguage(stdin)

	if len(os.Args) < 2 {
		cmdAnalyze(stdin)
		return
	}

	switch os.Args[1] {
	case "login":
		cmdLogin(stdin)
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
	fmt.Println(t("usage"))
}

const modulePath = "github.com/ask4git/repotale/repotale"

func cmdUpdate() {
	if _, err := exec.LookPath("go"); err != nil {
		die(fmt.Errorf("%s", t("update_go_missing")))
	}

	fmt.Println(t("update_start", modulePath))
	cmd := exec.Command("go", "install", modulePath+"@latest")
	// GOPROXY=direct: proxy.golang.org caches @latest for a while, which would
	// otherwise make `update` reinstall a stale version right after a fresh push.
	// GOSUMDB=off: sum.golang.org has its own separate propagation delay and
	// 404s on a version it hasn't indexed yet, even with GOPROXY=direct.
	cmd.Env = append(os.Environ(), "GOPROXY=direct", "GOSUMDB=off")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		die(fmt.Errorf(t("update_failed"), err))
	}
	fmt.Println(t("update_done"))
}

func cmdVersion() {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		fmt.Println(t("version_dev"))
		return
	}
	fmt.Println("repotale", info.Main.Version)
}

func cmdLogin(reader *bufio.Reader) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Print(t("login_prompt"))
		line, err := reader.ReadString('\n')
		// io.EOF with data already read just means no trailing newline
		// (e.g. `echo -n "$TOKEN" | repotale login`) — the token is still valid.
		if err != nil && err != io.EOF {
			die(fmt.Errorf(t("login_reading_token"), err))
		}
		token = strings.TrimSpace(line)
	}

	if token == "" {
		die(fmt.Errorf("%s", t("login_no_token")))
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

	fmt.Println(t("login_success"))
}

func cmdConnect() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, t("connect_usage"))
		os.Exit(1)
	}
	ownerRepo := os.Args[2]

	parts := strings.Split(ownerRepo, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		die(fmt.Errorf("%s", t("connect_bad_format")))
	}

	cfg, err := loadConfig()
	if err != nil {
		die(err)
	}
	if cfg.Token == "" {
		die(fmt.Errorf("%s", t("connect_not_logged_in")))
	}

	if err := validateRepo(cfg.Token, ownerRepo); err != nil {
		die(err)
	}

	cfg.Repo = ownerRepo
	if err := saveConfig(cfg); err != nil {
		die(err)
	}

	fmt.Println(t("connect_success", ownerRepo))
}

func cmdOpen() {
	cfg, err := loadConfig()
	if err != nil {
		die(err)
	}
	if cfg.Repo == "" {
		die(fmt.Errorf("%s", t("open_no_repo")))
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
