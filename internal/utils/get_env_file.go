package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetEnvFilePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		envPath := filepath.Join(dir, "config", ".env")
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("config/.env not found")
		}

		dir = parent
	}
}
