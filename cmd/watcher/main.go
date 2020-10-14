package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/mattermost/mattermost-plugin-starter-template/internal/watcher"
)

func main() {
	nagiosCfgDir := flag.String("dir", "/usr/local/nagios/etc/", "Nagios configuration files directory")
	flag.Parse()

	baseDir := *nagiosCfgDir

	if !filepath.IsAbs(baseDir) {
		log.Fatal("dir argument must be an absolute path, like /usr/local/nagios/etc/")
	}

	files, directories, err := watcher.GetAllInDirectory(baseDir)
	if err != nil {
		log.Panic(err)
	}

	diffWatch := watcher.NewDifferential()

	done := make(chan struct{})

	go func() {
		defer close(done)

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		<-interrupt

		// We received an interrupt signal, shut down.
		log.Println("Bye")
	}()

	if err := watcher.WatchDirectories(directories, diffWatch, done); err != nil {
		log.Panic(err)
	}
}
