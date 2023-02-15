package watcher

import (
	"bytes"
	"crypto/md5" //nolint:gosec
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/google/go-cmp/cmp"
)

// GetAllInDirectory recursively returns all paths to files and directories in
// dir which have allowed extensions. It returns nil, nil, <err> on
// the first error encountered.
func GetAllInDirectory(dir string, allowedExtensions []string) (
	[]string,
	[]string,
	error) {
	var files, directories []string

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			directories = append(directories, path)
			return nil
		}

		if !isExtensionAllowed(allowedExtensions, filepath.Ext(path)) {
			return nil
		}

		files = append(files, path)

		return nil
	}

	if err := filepath.Walk(dir, walkFn); err != nil {
		return nil, nil, fmt.Errorf("filepath.Walk: %w", err)
	}

	return files, directories, nil
}

func isExtensionAllowed(allowedExtensions []string, ext string) bool {
	for _, v := range allowedExtensions {
		if v == ext {
			return true
		}
	}

	return false
}

// WatchFuncProvider provides functionality to watch a directory.
type WatchFuncProvider interface {
	WatchFn(path string) error
}

// DiffSender sends diff to somewhere.
type DiffSender interface {
	Send(path string, diff string) error
}

// RemoteDiffSender sends diff to remote server.
type RemoteDiffSender struct {
	url    string
	token  string
	client *http.Client
}

// Send implements DiffSender.
func (d *RemoteDiffSender) Send(path string, diff string) error {
	change := Change{
		Name: filepath.Base(path),
		Diff: diff,
	}

	b, err := json.Marshal(change)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, d.url, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Content-Type", http.DetectContentType(b))
	req.Header.Set(TokenHeader, d.token)

	res, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("d.client.Do: %w", err)
	}

	defer res.Body.Close()

	if _, err := io.Copy(io.Discard, res.Body); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	if c := res.StatusCode; !checkStatusCode2xx(c) {
		return fmt.Errorf("server returned non-2xx status code (%d)", c)
	}

	return nil
}

// WatchDirectories watches for changes in directories and calls WatchFn on
// every change. It terminates after done is closed.
func WatchDirectories(
	directories []string,
	provider WatchFuncProvider,
	done <-chan struct{}) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("fsnotify.NewWatcher: %w", err)
	}
	defer w.Close()

	for _, d := range directories {
		if err := w.Add(d); err != nil {
			return fmt.Errorf("w.Add: %w", err)
		}
	}

	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return nil
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				if err := provider.WatchFn(event.Name); err != nil {
					log.Printf("WatchFn(%s): %v", event.Name, err)
				}
			}
		case err, ok := <-w.Errors:
			if !ok {
				return nil
			}

			log.Printf("Received an error from Errors queue: %v", err)
		case _, ok := <-done:
			if !ok {
				return nil
			}
		}
	}
}

// Differential implements WatchFuncProvider. Use NewDifferential to initialize
// Differential.
type Differential struct {
	allowedExtensions []string
	previousChecksum  map[string][16]byte
	previousContents  map[string][]byte
	client            *http.Client
	url, token        string
	diffSender        DiffSender
}

// Change represents a file differences from last modifications.
type Change struct {
	Name string
	Diff string
}

// TokenHeader for setting Header send Diff
const TokenHeader = "X-Nagios-Plugin-Token" //nolint:gosec

func checkStatusCode2xx(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices
}

// WatchFn for reading files in the watched folder and look for difference from the last update
func (d Differential) WatchFn(path string) error {
	if !isExtensionAllowed(d.allowedExtensions, filepath.Ext(path)) {
		return nil
	}

	// This is to allow for changes to propagate on the filesystem. If the file
	// is large, the write won't be atomic. It will happen in, for example, 4096
	// bytes chunks. 1 ms should be enough.
	time.Sleep(time.Millisecond)

	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("os.Stat: %w", err)
	}

	if fileInfo.Size() > 100*1024 {
		log.Printf("Rejected too big file, path: %v, size: %v", path, fileInfo.Size())
		return nil
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("os.ReadFile: %w", err)
	}

	checksum := md5.Sum(contents) //nolint:gosec

	if checksum == d.previousChecksum[path] {
		return nil
	}

	diff := cmp.Diff(string(d.previousContents[path]), string(contents))

	if err := d.diffSender.Send(path, diff); err != nil {
		return fmt.Errorf("d.sendDiff: %w", err)
	}

	log.Printf("Sent the diff (size = %d)", len(diff))

	d.previousChecksum[path] = checksum
	d.previousContents[path] = contents

	return nil
}

// NewDifferential returns initialized Differential.
func NewDifferential(
	allowedExtensions, initialFilePaths []string,
	httpClient *http.Client,
	url, token string) (Differential, error) {
	previousChecksum := make(map[string][16]byte)
	previousContents := make(map[string][]byte)

	for _, p := range initialFilePaths {
		b, err := os.ReadFile(p)
		if err != nil {
			return Differential{}, fmt.Errorf("os.ReadFile: %w", err)
		}

		previousChecksum[p] = md5.Sum(b) //nolint:gosec
		previousContents[p] = b
	}

	return Differential{
		allowedExtensions: allowedExtensions,
		previousChecksum:  previousChecksum,
		previousContents:  previousContents,
		client:            httpClient,
		url:               url,
		token:             token,
		diffSender: &RemoteDiffSender{
			url:    url,
			token:  token,
			client: httpClient,
		},
	}, nil
}
