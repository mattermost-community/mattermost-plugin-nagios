package watcher

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/google/go-cmp/cmp"
)

// GetAllInDirectory recursively returns all paths to files and directories in
// dir. It returns nil, nil, <err> on the first error encountered.
func GetAllInDirectory(dir string) ([]string, []string, error) {
	var files, directories []string

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			directories = append(directories, path)
		} else {
			files = append(files, path)
		}

		return nil
	}

	if err := filepath.Walk(dir, walkFn); err != nil {
		return nil, nil, fmt.Errorf("Walk: %w", err)
	}

	return files, directories, nil
}

type WatchFuncProvider interface {
	WatchFn(path string) error
}

// WatchDirectories watches for changes in directories and calls WatchFn on
// every change. It terminates after done is closed.
func WatchDirectories(
	directories []string,
	provider WatchFuncProvider,
	done <-chan struct{}) error {

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("NewWatcher: %w", err)
	}
	defer w.Close()

	for _, d := range directories {
		if err := w.Add(d); err != nil {
			return fmt.Errorf("Add: %w", err)
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
	ignoredExtensions map[string]struct{}
	previousChecksum  map[string][16]byte
	previousContents  map[string][]byte
	client            *http.Client
	url, token        string
}

type Change struct {
	Name string
	Diff string
}

func checkStatusCode2xx(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices
}

func (d Differential) sendDiff(path string, diff string) error {
	change := Change{
		Name: filepath.Base(path),
		Diff: diff,
	}

	b, err := json.Marshal(change)
	if err != nil {
		return fmt.Errorf("Encode: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, d.url, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("NewRequest: %w", err)
	}

	res, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("Do: %w", err)
	}

	if c := res.StatusCode; !checkStatusCode2xx(c) {
		return fmt.Errorf("server returned non-2xx status code (%d)", c)
	}

	return nil
}

func (d Differential) WatchFn(path string) error {
	if _, ok := d.ignoredExtensions[filepath.Ext(path)]; ok {
		return nil
	}

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("ReadFile: %w", err)
	}

	checksum := md5.Sum(contents)

	if checksum == d.previousChecksum[path] {
		return nil
	}

	diff := cmp.Diff(string(d.previousContents[path]), string(contents))

	if err := d.sendDiff(path, diff); err != nil {
		return fmt.Errorf("sendDiff: %w", err)
	}

	log.Printf("Sent the diff (length = %d)", len(diff))

	d.previousChecksum[path] = checksum
	d.previousContents[path] = contents

	return nil
}

func getIgnoredExtensions(extensions []string) map[string]struct{} {
	lookup := make(map[string]struct{})

	for _, e := range extensions {
		lookup[e] = struct{}{}
	}

	return lookup
}

// NewDifferential returns initialized Differential.
func NewDifferential(
	ignoredExtensions, initialFilePaths []string,
	httpClient *http.Client,
	url, token string) (Differential, error) {

	previousChecksum := make(map[string][16]byte)
	previousContents := make(map[string][]byte)

	for _, p := range initialFilePaths {
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return Differential{}, fmt.Errorf("ReadFile: %w", err)
		}
		previousChecksum[p] = md5.Sum(b)
		previousContents[p] = b
	}

	return Differential{
		ignoredExtensions: getIgnoredExtensions(ignoredExtensions),
		previousChecksum:  previousChecksum,
		previousContents:  previousContents,
		client:            httpClient,
		url:               url,
		token:             token,
	}, nil
}
