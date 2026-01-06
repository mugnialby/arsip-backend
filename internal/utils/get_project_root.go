package utils

import (
	"errors"
	"os"
	"path/filepath"
)

func GetProjectRoot() (string, error) {

	// 1️⃣ Explicit override (BEST PRACTICE)
	// Example:
	//   APP_ROOT=/app
	if root := os.Getenv("APP_ROOT"); root != "" {
		return filepath.Clean(root), nil
	}

	// 2️⃣ Working directory (Docker WORKDIR)
	if wd, err := os.Getwd(); err == nil {
		if isValidRoot(wd) {
			return wd, nil
		}
	}

	// 3️⃣ Executable directory (fallback)
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		if isValidRoot(exeDir) {
			return exeDir, nil
		}
	}

	return "", errors.New("project root not found: set APP_ROOT environment variable")
}

func isValidRoot(dir string) bool {
	// Relaxed checks — container-safe
	if exists(filepath.Join(dir, "storage")) {
		return true
	}

	if exists(filepath.Join(dir, "config")) {
		return true
	}

	return false
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
