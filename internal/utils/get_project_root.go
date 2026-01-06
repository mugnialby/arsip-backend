package utils

import (
	"errors"
	"os"
	"path/filepath"
)

func GetProjectRoot() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(exe)

	for {
		// Marker 1: go.mod
		if exists(filepath.Join(dir, "go.mod")) {
			return dir, nil
		}

		// Marker 2: config/.env
		if exists(filepath.Join(dir, "config", ".env")) {
			return dir, nil
		}

		// Stop at filesystem root
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("project root not found")
		}

		dir = parent
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
