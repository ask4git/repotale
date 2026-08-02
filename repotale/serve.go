package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed templates/local.html.tmpl
var templatesFS embed.FS

var localTmpl = template.Must(template.ParseFS(templatesFS, "templates/local.html.tmpl"))

func loadLocalPosts(postsDir string) ([]localPost, error) {
	entries, err := os.ReadDir(postsDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var posts []localPost
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(postsDir, entry.Name()))
		if err != nil {
			return nil, err
		}
		var post localPost
		if err := json.Unmarshal(data, &post); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", entry.Name(), err)
		}
		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PublishedAt > posts[j].PublishedAt
	})
	return posts, nil
}

// serveLocal runs a self-contained web view of posts (no Node/Next.js/DB
// needed - just this binary) and blocks until the process is killed
// (Ctrl+C), the same way `next dev` or `python -m http.server` would.
func serveLocal(repoID string, posts []localPost) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Lang         string
			RepoID       string
			Posts        []localPost
			Heading      string
			Empty        string
			CommitLabel  string
			SessionLabel string
		}{
			Lang:         string(currentLang),
			RepoID:       repoID,
			Posts:        posts,
			Heading:      t("web_heading"),
			Empty:        t("web_empty"),
			CommitLabel:  t("web_commit_label"),
			SessionLabel: t("web_session_label"),
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := localTmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	ln, err := net.Listen("tcp", "127.0.0.1:4321")
	if err != nil {
		// 4321 taken (e.g. another repotale instance) - fall back to any free port
		ln, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return fmt.Errorf("starting local server: %w", err)
		}
	}

	url := "http://" + ln.Addr().String() + "/"
	fmt.Println(t("serving_line", url))
	_ = openBrowser(url)

	return http.Serve(ln, mux)
}
