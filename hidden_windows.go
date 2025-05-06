//go:build windows

package scanner

import (
	"os"
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
