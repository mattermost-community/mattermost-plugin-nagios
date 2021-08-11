package watcher

import (
	"fmt"
	"net/http"
	"runtime"
	"testing"
)

func TestWatchFn(t *testing.T) {

	// MaxFileSize = 17*10 ^ 3
	// MaxReadSize = 5 * 1024

	ignoredExtension := []string{".swp"}
	// path := "/home/hasbi/file_watcher/100MB.bin"
	path := "/home/hasbi/file_watcher/BTOP.txt"
	// path := "/home/hasbi/file_watcher/file.txt"

	files := []string{path}

	differential, err := NewDifferential(ignoredExtension, files, http.DefaultClient, "http://localhost:8065", "2137", "/home/hasbi/")
	if err != nil {
		fmt.Printf("%+v \n", err)
		return
	}

	err = differential.WatchFn(path)
	if err != nil {
		fmt.Printf("%+v \n", err)
		return
	}

}

func CheckMem() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	fmt.Printf("\n")
	fmt.Printf("Alloc: %d MB, TotalAlloc: %d MB, Sys: %d MB\n",
		ms.Alloc/1024/1024, ms.TotalAlloc/1024/1024, ms.Sys/1024/1024)
	fmt.Printf("Mallocs: %d, Frees: %d\n",
		ms.Mallocs, ms.Frees)
	fmt.Printf("HeapAlloc: %d MB, HeapSys: %d MB, HeapIdle: %d MB\n",
		ms.HeapAlloc/1024/1024, ms.HeapSys/1024/1024, ms.HeapIdle/1024/1024)
	fmt.Printf("HeapObjects: %d\n", ms.HeapObjects)
	fmt.Printf("\n")
}
