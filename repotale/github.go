package main

import (
	"fmt"
	"net/http"
)

func githubRequest(url, token string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	return http.DefaultClient.Do(req)
}

// validateToken checks that token is accepted by the GitHub API.
func validateToken(token string) error {
	resp, err := githubRequest("https://api.github.com/user", token)
	if err != nil {
		return fmt.Errorf("could not reach GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub rejected the token (status %d)", resp.StatusCode)
	}
	return nil
}

// validateRepo checks that owner/repo exists and is reachable with token.
func validateRepo(token, ownerRepo string) error {
	url := "https://api.github.com/repos/" + ownerRepo
	resp, err := githubRequest(url, token)
	if err != nil {
		return fmt.Errorf("could not reach GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("repo %q not reachable (status %d)", ownerRepo, resp.StatusCode)
	}
	return nil
}
