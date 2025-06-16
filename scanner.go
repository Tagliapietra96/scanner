// Package scanner provides utilities for recursive directory traversal with filtering options.
package scanner

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	// ErrDoNotInclude is returned by callback functions to indicate that the current
	// file or directory should not be included in the scan results, but traversal
	// should continue normally if it's a directory.
	ErrDoNotInclude = errors.New("the current element must be not included in the result")
	// ErrSkip is returned by callback functions to indicate that the current
	// file or directory should be completely skipped. If it's a directory,
	// its contents will not be traversed.
	ErrSkip = errors.New("the current element must be skipped, if it is a directory it will not used to traverse further")
)

// Scan asynchronously traverses the directory structure starting at root path.
// It uses a callback function to filter entries and sends results through channels.
// The callback function receives the full path and os.DirEntry for each file/directory.
// It can return ErrSkip to skip an entry entirely, ErrDoNotInclude to exclude from
// results but continue traversal, or any other error to abort scanning.
// Matching paths are sent to rc channel and errors to ec channel.
// Both channels are closed when the scan completes.
// The scan uses goroutines with a semaphore to limit concurrent directory reads
// based on available CPU cores.
func Scan(root string, callBack func(string, os.DirEntry) error, rc chan<- string, ec chan<- error) {
	go func() {
		defer close(rc)
		defer close(ec)

		var wg sync.WaitGroup
		s := make(chan int8, max(1, runtime.NumCPU()-2))

		var do func(string)
		do = func(p string) {
			defer wg.Done()

			des, err := os.ReadDir(p)
			if err != nil {
				ec <- err
				return
			}

			for _, de := range des {
				var e error
				np := filepath.Join(p, de.Name())

				if callBack != nil {
					e = callBack(np, de)
					if e != nil && e != ErrSkip && e != ErrDoNotInclude {
						ec <- e
					}
				}

				if e == ErrSkip {
					continue
				}

				if e != ErrDoNotInclude {
					rc <- np
				}

				if de.IsDir() {
					wg.Add(1)
					go func() {
						s <- 0
						defer func() { <-s }()
						do(np)
					}()
				}
			}
		}

		wg.Add(1)
		s <- 0
		go func() {
			defer func() { <-s }()
			do(root)
		}()

		wg.Wait()
	}()
}

// ScanSync synchronously scans the directory structure starting at root path.
// It applies the callback function to each entry and returns a slice of matching paths.
// This is a convenience wrapper around Scan that handles channel management internally.
// The callback function follows the same contract as Scan: return ErrSkip to skip
// entries, ErrDoNotInclude to exclude from results, or other errors to abort.
// Returns the collected paths and the first error encountered, if any.
func ScanSync(root string, callBack func(string, os.DirEntry) error) ([]string, error) {
	rc := make(chan string)
	ec := make(chan error)

	Scan(root, callBack, rc, ec)
	r := make([]string, 0, 1024)

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
