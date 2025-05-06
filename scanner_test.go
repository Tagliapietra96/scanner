package scanner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tagliapietra96/scanner"
)

func compareSlices(a []string, b []string) []string {
	var diff []string
	m := make(map[string]bool)
	for _, item := range a {
		_, ok := m[item]
		if ok {
			diff = append(diff, item)
		} else {
			m[item] = true
		}
	}
	for _, item := range b {
		if _, found := m[item]; !found {
			diff = append(diff, item)
		}
	}
	return diff
}

func BenchmarkFilepathWalk(b *testing.B) {
	for b.Loop() {
		var fileCount int

		err := filepath.Walk("/Users", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fileCount++
			return nil
		})

		if err != nil {
			b.Fatalf("filepath.Walk failed: %v", err)
		}
		if fileCount == 0 {
			b.Fatalf("filepath.Walk found no files")
		}
	}
}

func BenchmarkWalkDir(b *testing.B) {
	for b.Loop() {
		var fileCount int

		err := filepath.WalkDir("/Users", func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			fileCount++
			return nil
		})

		if err != nil {
			b.Fatalf("filepath.WalkDir failed: %v", err)
		}
		if fileCount == 0 {
			b.Fatalf("filepath.WalkDir found no files")
		}
	}
}

func BenchmarkScanSync(b *testing.B) {
	for b.Loop() {
		var fileCount int

		r, err := scanner.ScanSync("/Users", -1, nil)
		fileCount = len(r)

		if err != nil {
			b.Fatalf("Scanner failed: %v", err)
		}
		if fileCount == 0 {
			b.Fatalf("Scanner found no files")
		}
	}
}

func TestScan(t *testing.T) {
	root := "/Users"
	r, err := scanner.ScanSync(root, -1, nil)
	if err != nil {
		t.Fatalf("Scanner failed: %v", err)
	}
	if len(r) == 0 {
		t.Fatalf("Scanner found no files")
	}

	lr := len(r)

	res := []string{}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		res = append(res, path)
		return nil
	})

	if err != nil {
		t.Fatalf("filepath.Walk failed: %v", err)
	}

	if len(res) != lr {
		diff := compareSlices(r, res)

		for _, d := range diff {
			t.Logf("Difference: %s", d)
		}
		t.Fatalf("Scanner found %d files, but filepath.Walk found %d files", lr, len(res))
	}

	res = []string{}
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		res = append(res, path)
		return nil
	})

	if err != nil {
		t.Fatalf("filepath.WalkDir failed: %v", err)
	}

	if len(res) != lr {
		diff := compareSlices(r, res)

		for _, d := range diff {
			t.Logf("Difference: %s", d)
		}
		t.Fatalf("Scanner found %d files, but filepath.WalkDir found %d files", lr, len(res))
	}
}
