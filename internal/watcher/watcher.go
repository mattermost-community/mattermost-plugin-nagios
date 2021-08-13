package watcher

import (
	"bufio"
	"bytes"
	"crypto/md5" //nolint:gosec
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/google/go-cmp/cmp"
)

func getIgnoredExtensions(extensions []string) map[string]struct{} {
	lookup := make(map[string]struct{})

	for _, e := range extensions {
		lookup[e] = struct{}{}
	}

	return lookup
}

var (
	// MaxFileSize to set the maximum size of the file to be read and stored in memory,
	// if more than that, the file will be saved to a temporary folder
	MaxFileSize int64 = 100 * 1024 * 1024

	// MaxReadSize to set the buffer used to read the contents of the file each loop
	MaxReadSize int64 = 5 * 1024 * 1024

	// TemporaryDirectory the location of the folder to store temporary files
	TemporaryDirectory = os.Getenv("HOME") + "/temp/watcher"
)

// GetAllInDirectory recursively returns all paths to files and directories in
// dir (excluding files with ignored extensions). It returns nil, nil, <err> on
// the first error encountered.
func GetAllInDirectory(dir string, ignoredExtensions []string) (
	[]string,
	[]string,
	error) {
	var files, directories []string

	ignoredExtensionsLookup := getIgnoredExtensions(ignoredExtensions)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if _, ok := ignoredExtensionsLookup[filepath.Ext(path)]; ok {
			return nil
		}

		if info.IsDir() {
			directories = append(directories, path)
		} else {
			files = append(files, path)
		}

		return nil
	}

	if err := filepath.Walk(dir, walkFn); err != nil {
		return nil, nil, fmt.Errorf("filepath.Walk: %w", err)
	}

	return files, directories, nil
}

// WatchFuncProvider for interfacing WatchFn for WatchDirectories
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
	ignoredExtensions map[string]struct{}
	previousChecksum  map[string][16]byte
	previousContents  map[string][]byte
	client            *http.Client
	url, token        string
}

// Change struct for send Difference
type Change struct {
	Name string
	Diff string
}

// TokenHeader for setting Header send Diff
const TokenHeader = "X-Nagios-Plugin-Token" //nolint:gosec

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

	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	if c := res.StatusCode; !checkStatusCode2xx(c) {
		return fmt.Errorf("server returned non-2xx status code (%d)", c)
	}

	return nil
}

// WatchFn for reading files in the watched folder and look for difference from the last update
func (d Differential) WatchFn(path string) error {
	if _, ok := d.ignoredExtensions[filepath.Ext(path)]; ok {
		return nil
	}

	time.Sleep(10 * time.Millisecond)

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error read file status %w", err)
	}

	if info.Size() >= MaxFileSize {
		log.Printf("File equal or higher than %v MB", MaxFileSize/1024/1024)

		srcFile, ok := d.previousContents[path]
		filePath := after(string(srcFile), "FilePath##")
		containFilePath := strings.Contains(string(srcFile), "FilePath##")

		if ok && !containFilePath {
			// After watched file updated, the size higher than threshold

			err = CopyFile(path, TemporaryDirectory+"/"+info.Name())
			if err != nil {
				return fmt.Errorf("error copy file  %w", err)
			}

			// find diff
			contents, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("ioutil.ReadFile: %w", err)
			}

			checksum := md5.Sum(contents) //nolint:gosec
			if checksum == d.previousChecksum[path] {
				return nil
			}

			diff := cmp.Diff(string(d.previousContents[path]), string(contents))

			// set new previous content
			contentName := "FilePath##" + TemporaryDirectory + "/" + info.Name()
			d.previousContents[path] = []byte(contentName)

			if err := d.sendDiff(path, diff); err != nil {
				return fmt.Errorf("d.sendDiff: %w", err)
			}

			log.Printf("Sent the diff (size = %d)", len(diff))

			return nil
		} else if ok && containFilePath {
			// Before and after update watched file is higher than threshold

			err := d.readFileAndFindDiff(path, filePath)
			if err != nil {
				return fmt.Errorf("error read file and make diff %w", err)
			}

			err = CopyFile(path, filePath)
			if err != nil {
				return fmt.Errorf("error copy file  %w", err)
			}

			return nil

		} else {
			// New files in watched directory

			fileTemp := TemporaryDirectory + "/" + info.Name()
			err = CopyFile(path, fileTemp)
			if err != nil {
				return fmt.Errorf("error copy file  %w", err)
			}

			err = d.readFileAndFindDiff(path, filePath)
			if err != nil {
				return fmt.Errorf("error read file and make diff %w", err)
			}
			return nil
		}

	} else {
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("ioutil.ReadFile: %w", err)
		}

		checksum := md5.Sum(contents) //nolint:gosec

		if checksum == d.previousChecksum[path] {
			return nil
		}

		diff := cmp.Diff(string(d.previousContents[path]), string(contents))

		if err := d.sendDiff(path, diff); err != nil {
			return fmt.Errorf("d.sendDiff: %w", err)
		}

		log.Printf("Sent the diff (size = %d)", len(diff))

		d.previousChecksum[path] = checksum
		d.previousContents[path] = contents
	}
	return nil
}

// NewDifferential returns initialized Differential.
func NewDifferential(
	ignoredExtensions, initialFilePaths []string,
	httpClient *http.Client,
	url, token, tempDir string) (Differential, error) {
	previousChecksum := make(map[string][16]byte)
	previousContents := make(map[string][]byte)

	if tempDir != "" {
		TemporaryDirectory = tempDir
	}
	for _, p := range initialFilePaths {
		info, err := os.Stat(p)
		if err != nil {
			return Differential{}, fmt.Errorf("error read file status %w", err)
		}

		if info.Size() >= MaxFileSize {
			log.Printf("File equal or higher than %v MB", MaxFileSize/1024/1024)

			_, err = os.Stat(TemporaryDirectory)
			if os.IsNotExist(err) {
				err := os.MkdirAll(TemporaryDirectory, 0777)
				if err != nil {
					return Differential{}, fmt.Errorf("error create temporary folder %w", err)
				}
			}

			filePath := TemporaryDirectory + "/" + info.Name()
			err = CopyFile(p, filePath)
			if err != nil {
				return Differential{}, fmt.Errorf("error copy file  %w", err)
			}

			contentName := "FilePath##" + filePath
			// Only write temporary folder directory
			previousContents[p] = []byte(contentName)

		} else {
			b, err := ioutil.ReadFile(p)
			if err != nil {
				return Differential{}, fmt.Errorf("ioutil.ReadFile: %w", err)
			}

			previousChecksum[p] = md5.Sum(b) //nolint:gosec
			previousContents[p] = b
		}
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

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}

	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func after(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:]
}

func (d Differential) readFileAndFindDiff(path, srcFile string) (err error) {

	f1, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error openning watched file : %w", err)
	}
	f2, err := os.Open(srcFile)
	if err != nil {
		return fmt.Errorf("error openning temp file: %w", err)
	}
	defer f1.Close()
	defer f2.Close()

	r1 := bufio.NewReader(f1)
	r2 := bufio.NewReader(f2)

	eof1 := false
	eof2 := false

	for {
		buf1 := make([]byte, MaxReadSize)
		buf2 := make([]byte, MaxReadSize)
		// read first file
		n, err := r1.Read(buf1)
		if err == io.EOF {
			eof1 = true
		} else if err != nil {
			return fmt.Errorf("error read chunkz watched file: %w", err)
		}
		content1 := buf1[:n]

		// read second file
		n, err = r2.Read(buf2)
		if err == io.EOF {
			eof2 = true
		} else if err != nil {
			return fmt.Errorf("error read chunkz watched file: %w", err)
		}
		content2 := buf2[:n]

		if eof1 && eof2 {
			break
		}

		checksum1 := md5.Sum(content1) //nolint:gosec
		checksum2 := md5.Sum(content2) //nolint:gosec

		if checksum1 == checksum2 {
			continue
		}

		diff := cmp.Diff(string(content2), string(content1))

		if err := d.sendDiff(path, diff); err != nil {
			return fmt.Errorf("d.sendDiff: %w", err)
		}
	}
	return nil
}
