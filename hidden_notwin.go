//go:build !windows

package scanner

import (
	"path/filepath"
	"strings"
)

// IsHidden checks if the given path is a hidden file or directory.
func IsHidden(path string) bool {
	filename := filepath.Base(path)

	// Check if the filename starts with a dot (.)
	if strings.HasPrefix(filename, ".") {
		return true
	}
	// Check if the filename starts with a tilde (~)
	if strings.HasPrefix(filename, "~") {
		return true
	}
	// Check if the filename starts with a hash (#)
	if strings.HasPrefix(filename, "#") {
		return true
	}
	return false
}
