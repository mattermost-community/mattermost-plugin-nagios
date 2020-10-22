package watcher

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetAllInDirectory(t *testing.T) {
	const (
		filesMultiplier    = 10
		expectedDirCount   = 2
		expectedFilesCount = 2 * filesMultiplier
	)

	ignoredExtensions := []string{".swp"}
	ignoredExtensionsMap := getIgnoredExtensions(ignoredExtensions)

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
		if err = ioutil.WriteFile(file, []byte(":octopus:"), 0666); err != nil {
			t.Fatalf("ioutil.WriteFile: %v", err)
		}

		file = filepath.Join(subDir, fmt.Sprintf("test_file_%d", i))
		if err = ioutil.WriteFile(file, []byte(":octopus:"), 0666); err != nil {
			t.Fatalf("ioutil.WriteFile: %v", err)
		}

		file = filepath.Join(subDir, fmt.Sprintf("test_file_%d.swp", i))
		if err = ioutil.WriteFile(file, []byte(":octopus:"), 0666); err != nil {
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
		assert.Equal(t, expectedDirCount, len(directories))
	})

	t.Run("Ignored extensions", func(t *testing.T) {
		for _, f := range files {
			if _, ok := ignoredExtensionsMap[filepath.Ext(f)]; ok {
				assert.Fail(t, "Haven't excluded files with ignored extensions.")
				return
			}
		}
	})
}

type mockWatchFuncProvider struct {
	called bool
}

func (m *mockWatchFuncProvider) WatchFn(path string) error {
	m.called = true
	return nil
}

func TestWatchDirectories(t *testing.T) {
	mockWatchFuncProvider := &mockWatchFuncProvider{}

	baseDir, err := ioutil.TempDir("", "watcher_test1")
	if err != nil {
		t.Fatalf("ioutil.TempDir: %v", err)
	}

	defer os.RemoveAll(baseDir)

	file := filepath.Join(baseDir, "test_file")
	if err = ioutil.WriteFile(file, []byte(":octopus:"), 0666); err != nil {
		t.Fatalf("ioutil.WriteFile: %v", err)
	}

	files, directories, err := GetAllInDirectory(baseDir, []string{".swp"})
	if err != nil {
		t.Fatalf("GetAllInDirectory: %v", err)
	}

	done := make(chan struct{})

	go func() {
		for !mockWatchFuncProvider.called {
			time.Sleep(100 * time.Millisecond)
			if err := ioutil.WriteFile(files[0], []byte(":octopus: - :octopus:"), 0666); err != nil {
				t.Logf("ioutil.WriteFile: %v", err)
			}
		}

		close(done)
	}()

	if err := WatchDirectories(directories, mockWatchFuncProvider, done); err != nil {
		t.Fatalf("WatchDirectories: %v", err)
	}

	assert.Equal(t, true, mockWatchFuncProvider.called)
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
