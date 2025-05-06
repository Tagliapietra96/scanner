//go:build !windows

package scanner

import (
	"os"
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

// ConfigDir returns the full config directory for the given application name
// on Unix-like systems, using XDG_CONFIG_HOME or defaulting to $HOME/.config.
func ConfigDir(dir string) (string, error) {
	d, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, dir), nil
}

// DataDir returns the full data directory for the given application name
// on Unix-like systems, using XDG_DATA_HOME or defaulting to $HOME/.local/share.
func DataDir(dir string) (string, error) {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataHome = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dataHome, dir), nil
}

// CacheDir returns the full cache directory for the given application name
// using XDG_CACHE_HOME or defaulting to $HOME/.cache.
func CacheDir(dir string) (string, error) {
	d, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, dir), nil
}

// TempDir returns the OS temporary directory for the application
func TempDir(dir string) string {
	return filepath.Join(os.TempDir(), dir)
}
