package watcher

import (
	"crypto/md5" //nolint:gosec
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

	allowedExtensions := []string{".cfg"}

	baseDir, err := ioutil.TempDir("", "watcher_test_*")
	if err != nil {
		t.Fatalf("ioutil.TempDir: %v", err)
	}

	defer os.RemoveAll(baseDir)

	subDir, err := ioutil.TempDir(baseDir, "watcher_test_*")
	if err != nil {
		t.Fatalf("ioutil.TempDir: %v", err)
	}

	for i := 0; i < filesMultiplier; i++ {
		if _, err = ioutil.TempFile(baseDir, "*.cfg"); err != nil {
			t.Fatalf("TempFile: %v", err)
		}

		if _, err = ioutil.TempFile(baseDir, "*.swp"); err != nil {
			t.Fatalf("TempFile: %v", err)
		}

		if _, err = ioutil.TempFile(baseDir, "*"); err != nil {
			t.Fatalf("TempFile: %v", err)
		}

		if _, err = ioutil.TempFile(subDir, "*"); err != nil {
			t.Fatalf("TempFile: %v", err)
		}

		if _, err = ioutil.TempFile(subDir, "*.cfg"); err != nil {
			t.Fatalf("TempFile: %v", err)
		}

		if _, err = ioutil.TempFile(subDir, "*.swp"); err != nil {
			t.Fatalf("TempFile: %v", err)
		}
	}

	files, directories, err := GetAllInDirectory(baseDir, allowedExtensions)
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
			if !isExtensionAllowed(allowedExtensions, filepath.Ext(f)) {
				assert.Fail(t, "ignored file included")
				return
			}
		}
	})
}

type mockWatchFuncProvider struct {
	called    bool
	calledMtx sync.Mutex
}

func (m *mockWatchFuncProvider) WatchFn(string) error {
	m.calledMtx.Lock()
	defer m.calledMtx.Unlock()

	m.called = true

	return nil
}

func TestWatchDirectories(t *testing.T) {
	mock := &mockWatchFuncProvider{}

	baseDir, err := ioutil.TempDir("", "watcher_test_*")
	if err != nil {
		t.Fatalf("ioutil.TempDir: %v", err)
	}

	defer os.RemoveAll(baseDir)

	f, err := ioutil.TempFile(baseDir, "*")
	if err != nil {
		t.Fatalf("TempFile: %v", err)
	}

	_, directories, err := GetAllInDirectory(baseDir, []string{".cfg"})
	if err != nil {
		t.Fatalf("GetAllInDirectory: %v", err)
	}

	done := make(chan struct{})

	go func() {
		for {
			mock.calledMtx.Lock()
			if mock.called {
				mock.calledMtx.Unlock()
				break
			}
			mock.calledMtx.Unlock()

			if _, err := f.WriteString("test"); err != nil {
				t.Errorf("WriteString: %v", err)
			}

			<-time.After(100 * time.Millisecond) // constant backoff
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
			allowedExtensions: make([]string, 0),
			previousChecksum:  make(map[string][16]byte),
			previousContents:  make(map[string][]byte),
			client:            nil,
			url:               "",
			token:             "",
			diffSender: &RemoteDiffSender{
				url:    "",
				token:  "",
				client: nil,
			},
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

		baseDir, err := ioutil.TempDir("", "watcher_test_*")
		if err != nil {
			t.Fatalf("ioutil.TempDir: %v", err)
		}

		defer os.RemoveAll(baseDir)

		for i := 0; i < 10; i++ {
			var f *os.File

			f, err = ioutil.TempFile(baseDir, "*.cfg")
			if err != nil {
				t.Fatalf("TempFile: %v", err)
			}

			if _, err = f.Write([]byte{byte(i)}); err != nil {
				t.Fatalf("Write: %v", err)
			}

			var b []byte

			b, err = ioutil.ReadFile(f.Name())
			if err != nil {
				t.Fatalf("ioutil.ReadFile: %v", err)
			}

			previousChecksum[f.Name()] = md5.Sum(b) //nolint:gosec
			previousContents[f.Name()] = b
		}

		expected := Differential{
			allowedExtensions: []string{".cfg"},
			previousChecksum:  previousChecksum,
			previousContents:  previousContents,
			client:            http.DefaultClient,
			url:               "dummy",
			token:             "2137",
			diffSender: &RemoteDiffSender{
				url:    "dummy",
				token:  "2137",
				client: http.DefaultClient,
			},
		}

		files, _, err := GetAllInDirectory(baseDir, []string{".cfg"})
		if err != nil {
			t.Fatalf("GetAllInDirectory: %v", err)
		}

		actual, err := NewDifferential([]string{".cfg"}, files, http.DefaultClient, "dummy", "2137")
		if err != nil {
			t.Fatalf("NewDifferential: %v", err)
		}

		assert.Equal(t, expected, actual)
	})
}

func TestDifferentialIgnore(t *testing.T) {
	baseDir, err := ioutil.TempDir("", "watcher_test_*")
	if err != nil {
		t.Fatalf("ioutil.TempDir: %v", err)
	}

	defer os.RemoveAll(baseDir)

	file, err := ioutil.TempFile(baseDir, "*")
	if err != nil {
		t.Fatalf("ioutil.TempFile: %v", err)
	}
	defer file.Close()

	d, err := NewDifferential([]string{".cfg"}, []string{}, http.DefaultClient, "dummy", "2137")
	if err != nil {
		t.Fatalf("NewDifferential: %v", err)
	}

	if err = d.WatchFn(filepath.Join(baseDir, file.Name())); err != nil {
		t.Fatalf("WatchFn: %v", err)
	}

	assert.Equal(t, 0, len(d.previousChecksum))
}

type mockDiffSender struct {
	called    bool
	calledMtx sync.Mutex
}

func (d *mockDiffSender) Send(path string, diff string) error {
	d.calledMtx.Lock()
	defer d.calledMtx.Unlock()

	d.called = true

	return nil
}

func TestDifferentialFiltered(t *testing.T) {
	baseDir, err := ioutil.TempDir("", "watcher_test_*")
	if err != nil {
		t.Fatalf("ioutil.TempDir: %v", err)
	}

	defer os.RemoveAll(baseDir)

	file, err := ioutil.TempFile(baseDir, "*.cfg")
	if err != nil {
		t.Fatalf("ioutil.TempFile: %v", err)
	}
	defer file.Close()

	d, err := NewDifferential([]string{".cfg"}, []string{}, http.DefaultClient, "dummy", "2137")
	if err != nil {
		t.Fatalf("NewDifferential: %v", err)
	}

	mockedDiffSender := &mockDiffSender{}
	d.diffSender = mockedDiffSender

	if err = d.WatchFn(file.Name()); err != nil {
		t.Fatalf("WatchFn: %v", err)
	}

	assert.True(t, mockedDiffSender.called)
	assert.Equal(t, 1, len(d.previousChecksum))
}

func TestDifferentialLargeFile(t *testing.T) {
	baseDir, err := ioutil.TempDir("", "watcher_test_*")
	if err != nil {
		t.Fatalf("ioutil.TempDir: %v", err)
	}

	defer os.RemoveAll(baseDir)

	testContent, err := ioutil.ReadFile("testdata/test-large-file.cfg")
	if err != nil {
		t.Fatalf("ioutil.ReadFile: %v", err)
	}

	testFile := filepath.Join(baseDir, "config.cfg")

	err = ioutil.WriteFile(testFile, testContent, 0600)
	if err != nil {
		t.Fatalf("ioutil.WriteFile: %v", err)
	}

	d, err := NewDifferential([]string{".cfg"}, []string{}, http.DefaultClient, "dummy", "2137")
	if err != nil {
		t.Fatalf("NewDifferential: %v", err)
	}

	if err = d.WatchFn(testFile); err != nil {
		t.Fatalf("WatchFn: %v", err)
	}

	assert.Equal(t, 0, len(d.previousChecksum))
}
