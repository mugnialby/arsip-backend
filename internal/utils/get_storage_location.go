package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetStorageLocation() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		envPath := filepath.Join(dir, "storage")
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("Storage directory not found")
		}

		dir = parent
	}
}
