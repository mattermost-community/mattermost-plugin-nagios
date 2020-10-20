package watcher

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWatchDirectories(t *testing.T) {
	dummyURL := "http://dummy.restapiexample.com/api/v1/create"

	dir1, err := ioutil.TempDir("", "watcher_test1")
	if err != nil {
		t.Fatalf("ioutil.TempDir: %v", err)
	}
	defer os.RemoveAll(dir1)

	dir2, err := ioutil.TempDir("", "watcher_test2")
	if err != nil {
		t.Fatalf("ioutil.TempDir: %v", err)
	}
	defer os.RemoveAll(dir2)

	directories := []string{dir1, dir2}
	var files []string

	fname := filepath.Join(dir1, "test_file")
	if err := ioutil.WriteFile(fname, []byte(":octopus:"), 0644); err != nil {
		t.Fatalf("ioutil.WriteFile: %v", err)
	}
	files = append(files, fname)

	fname = filepath.Join(dir2, "test_file")
	if err := ioutil.WriteFile(fname, []byte(":octopus:"), 0644); err != nil {
		t.Fatalf("ioutil.WriteFile: %v", err)
	}
	files = append(files, fname)

	fname = filepath.Join(dir2, "test_file.swp")
	if err := ioutil.WriteFile(fname, []byte(":octopus:"), 0644); err != nil {
		t.Fatalf("ioutil.WriteFile: %v", err)
	}
	files = append(files, fname)

	differential, err := NewDifferential([]string{".swp"}, files, http.DefaultClient, dummyURL, "2137")
	if err != nil {
		t.Fatalf("NewDifferential: %v", err)
	}

	t.Run("First directory", func(t *testing.T) {
		done := make(chan struct{})

		go func() {
			time.Sleep(1 * time.Second)
			if err := ioutil.WriteFile(files[0], []byte(":octopus: - :octopus:"), 0644); err != nil {
				t.Fatalf("ioutil.WriteFile: %v", err)
			}
			close(done)
		}()

		if err := WatchDirectories(directories, differential, done); err != nil {
			t.Fatalf("WatchDirectories: %v", err)
		}

		assert.Equal(t, []byte(":octopus: - :octopus:"), differential.previousContents[files[0]])
	})

	t.Run("Second directory", func(t *testing.T) {
		done := make(chan struct{})

		go func() {
			time.Sleep(1 * time.Second)
			if err := ioutil.WriteFile(files[1], []byte(":octopus: - :octopus:"), 0644); err != nil {
				t.Fatalf("ioutil.WriteFile: %v", err)
			}
			close(done)
		}()

		if err := WatchDirectories(directories, differential, done); err != nil {
			t.Fatalf("WatchDirectories: %v", err)
		}

		assert.Equal(t, []byte(":octopus: - :octopus:"), differential.previousContents[files[1]])
	})

	t.Run("Ignored extension", func(t *testing.T) {
		done := make(chan struct{})

		go func() {
			time.Sleep(1 * time.Second)
			if err := ioutil.WriteFile(files[2], []byte(":octopus: - :octopus:"), 0644); err != nil {
				t.Fatalf("ioutil.WriteFile: %v", err)
			}
			close(done)
		}()

		if err := WatchDirectories(directories, differential, done); err != nil {
			t.Fatalf("WatchDirectories: %v", err)
		}

		assert.Equal(t, []byte(":octopus:"), differential.previousContents[files[2]])
	})
}
