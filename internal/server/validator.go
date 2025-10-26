package server

import (
	"fmt"
	"os"
	"path/filepath"

	"fileserv/internal/models"
)

// ValidateDirectories validates and converts directory paths to Directory structs
func ValidateDirectories(paths []string) ([]models.Directory, error) {
	var dirs []models.Directory
	seen := make(map[string]bool)

	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("invalid path %s: %w", path, err)
		}

		// Check for duplicates
		if seen[absPath] {
			continue
		}
		seen[absPath] = true

		info, err := os.Stat(absPath)
		if err != nil {
			return nil, fmt.Errorf("cannot access %s: %w", path, err)
		}

		if !info.IsDir() {
			return nil, fmt.Errorf("%s is not a directory", path)
		}

		dirs = append(dirs, models.Directory{
			Name: filepath.Base(absPath),
			Path: absPath,
		})
	}

	return dirs, nil
}
