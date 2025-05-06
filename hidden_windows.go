//go:build windows

package scanner

import (
	"os"
	"path/filepath"
	"syscall"
)

// IsHidden checks if the given path is a hidden file or directory.
func IsHidden(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	data, ok := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	if ok {
		return data.FileAttributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0
	}
	return false
}

// ConfigDir returns the full config directory for the given application name
// on Windows, using %AppData% (roaming).
func ConfigDir(dir string) (string, error) {
	d, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, dir), nil
}

// DataDir returns the full data directory for the given application name
// on Windows, using %LocalAppData%.
func DataDir(dir string) (string, error) {
	d := os.Getenv("LOCALAPPDATA")
	if d == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		d = filepath.Join(home, "AppData", "Local")
	}
	return filepath.Join(d, dir), nil
}

// CacheDir returns the full cache directory for the given application name
// on Windows, using %LocalAppData%\Cache.
func CacheDir(dir string) (string, error) {
	d, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, dir), nil
}

// TempDir returns the OS temporary directory for the application
func TempDir(dir string) string {
	tmp := os.TempDir()
	return filepath.Join(tmp, dir)
}
