package utils

import (
	"fmt"
	"os"
)

// EnsurePathExists checks if a path exists, and if not, optionally creates it.
func EnsurePathExists(path string, createIfNotExist bool) error {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		if createIfNotExist {
			return os.MkdirAll(path, os.ModePerm)
		}
		return fmt.Errorf("path does not exist: %s", path)
	}
	// If error occurred (like permission issues)
	if err != nil {
		return fmt.Errorf("error checking path: %w", err)
	}

	// Path exists
	return nil
}
