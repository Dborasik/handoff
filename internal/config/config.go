package config

import (
	"os"
	"path/filepath"
)

// DBPath returns the path to the SQLite database file.
// Respects HANDOFF_DB env var, otherwise defaults to ~/.handoff/handoff.db
func DBPath() (string, error) {
	if p := os.Getenv("HANDOFF_DB"); p != "" {
		return p, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, ".handoff")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(dir, "handoff.db"), nil
}
