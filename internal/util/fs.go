package util

import (
	"os"
	"path/filepath"
)

func StateDir(app string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, "."+app)
	return dir, os.MkdirAll(dir, 0700)
}

