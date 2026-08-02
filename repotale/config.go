package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config is the persisted CLI state stored at ~/.repotale/config.json.
type Config struct {
	Token string `json:"token"`
	Repo  string `json:"repo"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".repotale", "config.json"), nil
}

// loadConfig returns the stored config, or a zero-value Config if no
// config file exists yet.
func loadConfig() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// saveConfig writes cfg to disk, creating ~/.repotale if needed.
// The file is created with 0600 permissions since it holds a secret token.
func saveConfig(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
