package watcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetAllInDirectory(t *testing.T) {
	expectedDirCount := 2

	dir1, err := ioutil.TempDir("", "watcher_test1")
	if err != nil {
		t.Fatalf("ioutil.TempDir: %v", err)
	}
	defer os.RemoveAll(dir1)

	dir2, err := ioutil.TempDir(dir1, "watcher_test2")
	if err != nil {
		t.Fatalf("ioutil.TempDir: %v", err)
	}

	expectedFilesCount := 0
	for i := 0; i < 10; i++ {
		file := filepath.Join(dir1, fmt.Sprintf("test_file_%d", i))
		if err := ioutil.WriteFile(file, []byte(":octopus:"), 0644); err != nil {
			t.Fatalf("ioutil.WriteFile: %v", err)
		}

		file = filepath.Join(dir2, fmt.Sprintf("test_file_%d", i))
		if err := ioutil.WriteFile(file, []byte(":octopus:"), 0644); err != nil {
			t.Fatalf("ioutil.WriteFile: %v", err)
		}

		file = filepath.Join(dir2, fmt.Sprintf("test_file_%d.swp", i))
		if err := ioutil.WriteFile(file, []byte(":octopus:"), 0644); err != nil {
			t.Fatalf("ioutil.WriteFile: %v", err)
		}

		expectedFilesCount += 2
	}

	ignoredExtensions := []string{".swp"}
	files, directories, err := GetAllInDirectory(dir1, ignoredExtensions)
	if err != nil {
		t.Fatalf("GetAllInDirectory: %v", err)
	}

	t.Run("Count files", func(t *testing.T) {
		assert.Equal(t, expectedFilesCount, len(files))
	})

	t.Run("Count directories", func(t *testing.T) {
		assert.Equal(t, expectedDirCount, len(directories))
	})

	t.Run("Ignored extensions", func(t *testing.T) {
		for _, f := range files {
			for _, e := range ignoredExtensions {
				if filepath.Ext(f) == e {
					assert.Fail(t, "Haven't excluded files with ignored extensions.")
					return
				}
			}
		}
	})
}

func TestWatchDirectories(t *testing.T) {
	dir1, err := ioutil.TempDir("", "watcher_test1")
	if err != nil {
		t.Fatalf("ioutil.TempDir: %v", err)
	}
	defer os.RemoveAll(dir1)

	dir2, err := ioutil.TempDir(dir1, "watcher_test2")
	if err != nil {
		t.Fatalf("ioutil.TempDir: %v", err)
	}

	file := filepath.Join(dir1, "test_file")
	if err := ioutil.WriteFile(file, []byte(":octopus:"), 0644); err != nil {
		t.Fatalf("ioutil.WriteFile: %v", err)
	}

	file = filepath.Join(dir2, "test_file")
	if err := ioutil.WriteFile(file, []byte(":octopus:"), 0644); err != nil {
		t.Fatalf("ioutil.WriteFile: %v", err)
	}

	ignoredExtensions := []string{".swp"}
	files, directories, err := GetAllInDirectory(dir1, ignoredExtensions)
	if err != nil {
		t.Fatalf("GetAllInDirectory: %v", err)
	}

	dummyURL := "http://dummy.restapiexample.com/api/v1/create"
	differential, err := NewDifferential(files, http.DefaultClient, dummyURL, "2137")
	if err != nil {
		t.Fatalf("NewDifferential: %v", err)
	}

	t.Run("Files in base directory", func(t *testing.T) {
		done := make(chan struct{})
		expected := ":octopus: - :octopus:"

		go func() {
			time.Sleep(1 * time.Second)
			if err := ioutil.WriteFile(files[0], []byte(expected), 0644); err != nil {
				t.Fatalf("ioutil.WriteFile: %v", err)
			}
			close(done)
		}()

		if err := WatchDirectories(directories, differential, done); err != nil {
			t.Fatalf("WatchDirectories: %v", err)
		}

		assert.Equal(t, []byte(expected), differential.previousContents[files[0]])
	})

	t.Run("Files in sub-directory", func(t *testing.T) {
		done := make(chan struct{})
		expected := ":octopus: - :octopus:"

		go func() {
			time.Sleep(1 * time.Second)
			if err := ioutil.WriteFile(files[1], []byte(expected), 0644); err != nil {
				t.Fatalf("ioutil.WriteFile: %v", err)
			}
			close(done)
		}()

		if err := WatchDirectories(directories, differential, done); err != nil {
			t.Fatalf("WatchDirectories: %v", err)
		}

		assert.Equal(t, []byte(expected), differential.previousContents[files[1]])
	})
}

func TestNewDifferential(t *testing.T) {
	t.Run("Empty struct", func(t *testing.T) {
		expected := Differential{
			previousChecksum: make(map[string][16]byte),
			previousContents: make(map[string][]byte),
			client:           nil,
			url:              "",
			token:            "",
		}

		actual, err := NewDifferential([]string{}, nil, "", "")
		if err != nil {
			t.Fatalf("NewDifferential: %v", err)
		}

		assert.Equal(t, expected, actual)
	})
}
