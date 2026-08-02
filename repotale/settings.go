package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// settingsPath returns ~/.repotale/settings - a plain `key=value` file
// (same spirit as .zshrc: human-editable, not JSON) for user preferences
// like language. Session state (token, connected repo) stays in
// config.json instead - different lifecycle, different concerns.
func settingsPath() string {
	return filepath.Join(repotaleDir(), "settings")
}

func loadSettings() (map[string]string, error) {
	settings := map[string]string{}

	f, err := os.Open(settingsPath())
	if os.IsNotExist(err) {
		return settings, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		settings[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return settings, scanner.Err()
}

func saveSettings(settings map[string]string) error {
	if err := os.MkdirAll(repotaleDir(), 0755); err != nil {
		return err
	}

	keys := make([]string, 0, len(settings))
	for k := range settings {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&b, "%s=%s\n", k, settings[k])
	}
	return os.WriteFile(settingsPath(), []byte(b.String()), 0644)
}

// setSetting reads-modifies-writes a single key, leaving any others intact.
func setSetting(key, value string) error {
	settings, err := loadSettings()
	if err != nil {
		return err
	}
	settings[key] = value
	return saveSettings(settings)
}
