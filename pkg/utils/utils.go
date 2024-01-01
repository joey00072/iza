package utils

import (
	"fmt"
	"os"
)

func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil // The directory does not exist
	}
	if err != nil {
		return false, err // An error other than "not exists", like permission issues
	}
	return info.IsDir(), nil // Return true if the path exists and is a directory
}

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

func TarExits(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func DeleteDir(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		// Handle the error according to your needs
		return err
	}
	return nil
}
