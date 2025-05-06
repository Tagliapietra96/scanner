// Package scanner provides utilities for recursive directory traversal with filtering options.
package scanner

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// scan recursively traverses the directory structure starting at path p.
// It respects the maximum depth m, applies the filter function fn to each entry,
// and sends matching paths to rc and errors to ec. It manages concurrency internally.
// If maxDepth is a negative value, it will traverse all levels of the directory tree.
func scan(p string, m int, fn func(string, os.DirEntry) bool, rc chan<- string, ec chan<- error) {
	var wg sync.WaitGroup
	s := make(chan string, max(1, runtime.NumCPU()/2))

	var do func(string, int)
	do = func(pp string, mm int) {
		defer wg.Done()
		des, err := os.ReadDir(pp)
		if err != nil {
			ec <- err
			return
		}

		for _, de := range des {
			if fn != nil && !fn(filepath.Join(pp, de.Name()), de) {
				continue
			}

			rc <- filepath.Join(pp, de.Name())
			if mm != 0 && de.IsDir() {
				wg.Add(1)
				go func() {
					s <- ""
					defer func() { <-s }()
					do(filepath.Join(pp, de.Name()), mm-1)
				}()
			}
		}
	}

	wg.Add(1)
	s <- ""
	go func() {
		defer func() { <-s }()
		do(p, m)
	}()

	wg.Wait()
}

// Scan asynchronously traverses the directory structure starting at root path.
// It respects the maximum depth, applies the filter function to each entry,
// and sends matching paths to rc and errors to ec. Both channels are closed when done.
// If maxDepth is a negative value, it will traverse all levels of the directory tree.
func Scan(root string, maxDepth int, filter func(string, os.DirEntry) bool, rc chan<- string, ec chan<- error) {
	go func() {
		defer close(rc)
		defer close(ec)
		scan(root, maxDepth, filter, rc, ec)
	}()
}

// ScanSync synchronously scans the directory structure starting at root path.
// It applies the filter function to each entry and returns a slice of matching paths.
// It provides a shorthand to scan the directory tree without needing to manage channels.
// it directly returns the results and errors.
// If maxDepth is a negative value, it will traverse all levels of the directory tree.
func ScanSync(root string, maxDepth int, filter func(string, os.DirEntry) bool) ([]string, error) {
	rc := make(chan string)
	ec := make(chan error)

	go Scan(root, maxDepth, filter, rc, ec)
	r := make([]string, 0)

	for {
		select {
		case rs, ok := <-rc:
			if !ok {
				return r, nil
			}
			r = append(r, rs)
		case err, ok := <-ec:
			if !ok || err == nil {
				continue
			}
			return r, err
		}
	}
}

// FilterDir returns true only for directory entries.
func FilterDir(_ string, de os.DirEntry) bool {
	return de.IsDir()
}

// FilterFile returns true only for non-directory entries.
func FilterFile(_ string, de os.DirEntry) bool {
	return !de.IsDir()
}

// FilterHidden returns true for entries that are hidden.
// Uses IsHidden to check if the path is a hidden file or directory.
func FilterHidden(p string, _ os.DirEntry) bool {
	return IsHidden(p)
}

// FilterRegular returns true only for regular file entries.
func FilterRegular(_ string, de os.DirEntry) bool {
	i, e := de.Info()
	if e != nil {
		return false
	}
	return i.Mode().IsRegular()
}

// FilterSymlink returns true only for symbolic link entries.
func FilterSymlink(_ string, de os.DirEntry) bool {
	i, e := de.Info()
	if e != nil {
		return false
	}
	return i.Mode()&os.ModeSymlink != 0
}

// FilterDevice returns true only for device file entries.
func FilterDevice(_ string, de os.DirEntry) bool {
	i, e := de.Info()
	if e != nil {
		return false
	}
	return i.Mode()&os.ModeDevice != 0
}

// FilterNamedPipe returns true only for named pipe entries.
func FilterNamedPipe(_ string, de os.DirEntry) bool {
	i, e := de.Info()
	if e != nil {
		return false
	}
	return i.Mode()&os.ModeNamedPipe != 0
}

// FilterSocket returns true only for socket file entries.
func FilterSocket(_ string, de os.DirEntry) bool {
	i, e := de.Info()
	if e != nil {
		return false
	}
	return i.Mode()&os.ModeSocket != 0
}

// FilterCharDev returns true only for character device entries.
func FilterCharDev(_ string, de os.DirEntry) bool {
	i, e := de.Info()
	if e != nil {
		return false
	}
	return i.Mode()&os.ModeCharDevice != 0
}

// FilterByExtension returns a filter function that matches files with the specified extension.
// The extension can be provided with or without the leading dot.
func FilterByExtension(e string) func(string, os.DirEntry) bool {
	return func(p string, de os.DirEntry) bool {
		if de.IsDir() {
			return false
		}
		ext := filepath.Ext(de.Name())
		return ext == e || ext == "."+e
	}
}

// FilterBySize returns a filter function that matches files based on their size.
// The op parameter specifies the comparison operator ("<", "<=", ">", ">=", "=", "==", "!=").
func FilterBySize(size int64, op string) func(string, os.DirEntry) bool {
	return func(_ string, de os.DirEntry) bool {
		i, e := de.Info()
		if e != nil {
			return false
		}

		s := i.Size()

		switch op {
		case "<":
			return s < size
		case "<=":
			return s <= size
		case ">":
			return s > size
		case ">=":
			return s >= size
		case "=", "==":
			return s == size
		case "!=":
			return s != size
		}
		return false
	}
}
