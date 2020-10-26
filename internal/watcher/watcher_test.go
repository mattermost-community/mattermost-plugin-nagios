package watcher

import (
	"crypto/md5" //nolint:gosec
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetAllInDirectory(t *testing.T) {
	const (
		filesMultiplier    = 10
		expectedDirsCount  = 2
		expectedFilesCount = 2 * filesMultiplier
	)

	ignoredExtensions := []string{".swp"}
	ignoredExtensionsLookup := getIgnoredExtensions(ignoredExtensions)

	baseDir, err := ioutil.TempDir("", "watcher_test1")
	if err != nil {
		t.Fatalf("ioutil.TempDir: %v", err)
	}

	defer os.RemoveAll(baseDir)

	subDir, err := ioutil.TempDir(baseDir, "watcher_test2")
	if err != nil {
		t.Fatalf("ioutil.TempDir: %v", err)
	}

	for i := 0; i < filesMultiplier; i++ {
		file := filepath.Join(baseDir, fmt.Sprintf("test_file_%d", i))
		if err = ioutil.WriteFile(file, []byte(":octopus:"), 0755); err != nil {
			t.Fatalf("ioutil.WriteFile: %v", err)
		}

		file = filepath.Join(baseDir, fmt.Sprintf("test_file_%d.swp", i))
		if err = ioutil.WriteFile(file, []byte(":octopus:"), 0755); err != nil {
			t.Fatalf("ioutil.WriteFile: %v", err)
		}

		file = filepath.Join(subDir, fmt.Sprintf("test_file_%d", i))
		if err = ioutil.WriteFile(file, []byte(":octopus:"), 0755); err != nil {
			t.Fatalf("ioutil.WriteFile: %v", err)
		}

		file = filepath.Join(subDir, fmt.Sprintf("test_file_%d.swp", i))
		if err = ioutil.WriteFile(file, []byte(":octopus:"), 0755); err != nil {
			t.Fatalf("ioutil.WriteFile: %v", err)
		}
	}

	files, directories, err := GetAllInDirectory(baseDir, ignoredExtensions)
	if err != nil {
		t.Fatalf("GetAllInDirectory: %v", err)
	}

	t.Run("Count files", func(t *testing.T) {
		assert.Equal(t, expectedFilesCount, len(files))
	})

	t.Run("Count directories", func(t *testing.T) {
		assert.Equal(t, expectedDirsCount, len(directories))
	})

	t.Run("Ignored extensions", func(t *testing.T) {
		for _, f := range files {
			if _, ok := ignoredExtensionsLookup[filepath.Ext(f)]; ok {
				assert.Fail(t, "Haven't excluded files with ignored extensions.")
				return
			}
		}
	})
}

type mockWatchFuncProvider struct {
	called      bool
	calledMutex sync.Mutex
}

func (m *mockWatchFuncProvider) WatchFn(path string) error {
	m.calledMutex.Lock()
	m.called = true
	m.calledMutex.Unlock()

	return nil
}

func TestWatchDirectories(t *testing.T) {
	mock := &mockWatchFuncProvider{}

	baseDir, err := ioutil.TempDir("", "watcher_test1")
	if err != nil {
		t.Fatalf("ioutil.TempDir: %v", err)
	}

	defer os.RemoveAll(baseDir)

	file := filepath.Join(baseDir, "test_file")
	if err = ioutil.WriteFile(file, []byte(":octopus:"), 0755); err != nil {
		t.Fatalf("ioutil.WriteFile: %v", err)
	}

	_, directories, err := GetAllInDirectory(baseDir, []string{".swp"})
	if err != nil {
		t.Fatalf("GetAllInDirectory: %v", err)
	}

	done := make(chan struct{})

	go func() {
		for {
			mock.calledMutex.Lock()
			if mock.called {
				mock.calledMutex.Unlock()
				break
			}
			mock.calledMutex.Unlock()

			if err := ioutil.WriteFile(file, []byte(":octopus: - :octopus:"), 0755); err != nil {
				t.Errorf("ioutil.WriteFile: %v", err)
			}

			<-time.After(100 * time.Millisecond)
		}

		close(done)
	}()

	if err := WatchDirectories(directories, mock, done); err != nil {
		t.Fatalf("WatchDirectories: %v", err)
	}

	assert.Equal(t, true, mock.called)
}

func TestNewDifferential(t *testing.T) {
	t.Run("Empty struct", func(t *testing.T) {
		expected := Differential{
			ignoredExtensions: getIgnoredExtensions(nil),
			previousChecksum:  make(map[string][16]byte),
			previousContents:  make(map[string][]byte),
			client:            nil,
			url:               "",
			token:             "",
		}

		actual, err := NewDifferential(nil, nil, nil, "", "")
		if err != nil {
			t.Fatalf("NewDifferential: %v", err)
		}

		assert.Equal(t, expected, actual)
	})

	t.Run("Checksum & contents", func(t *testing.T) {
		previousChecksum := make(map[string][16]byte)
		previousContents := make(map[string][]byte)

		baseDir, err := ioutil.TempDir("", "watcher_test1")
		if err != nil {
			t.Fatalf("ioutil.TempDir: %v", err)
		}

		defer os.RemoveAll(baseDir)

		for i := 0; i < 10; i++ {
			file := filepath.Join(baseDir, fmt.Sprintf("test_file_%d", i))
			if err = ioutil.WriteFile(file, []byte(fmt.Sprintf(":octopus:%d", i)), 0755); err != nil {
				t.Fatalf("ioutil.WriteFile: %v", err)
			}

			var b []byte
			b, err = ioutil.ReadFile(file)
			if err != nil {
				t.Fatalf("ioutil.ReadFile: %v", err)
			}

			previousChecksum[file] = md5.Sum(b) //nolint:gosec
			previousContents[file] = b
		}

		expected := Differential{
			ignoredExtensions: getIgnoredExtensions([]string{".swp"}),
			previousChecksum:  previousChecksum,
			previousContents:  previousContents,
			client:            http.DefaultClient,
			url:               "dummy",
			token:             "2137",
		}

		files, _, err := GetAllInDirectory(baseDir, []string{})
		if err != nil {
			t.Fatalf("GetAllInDirectory: %v", err)
		}

		actual, err := NewDifferential([]string{".swp"}, files, http.DefaultClient, "dummy", "2137")
		if err != nil {
			t.Fatalf("NewDifferential: %v", err)
		}

		assert.Equal(t, expected, actual)
	})
}
