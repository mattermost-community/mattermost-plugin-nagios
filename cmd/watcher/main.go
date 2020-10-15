package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/mattermost/mattermost-plugin-starter-template/internal/watcher"
)

func main() {
	var (
		nagiosCfgDir = flag.String("dir", "/usr/local/nagios/etc/", "Nagios configuration files directory")
		address      = flag.String("address", "", "Mattermost Server address")
		token        = flag.String("token", "", "Nagios plugin token")
	)
	flag.Parse()

	baseDir := *nagiosCfgDir

	if !filepath.IsAbs(baseDir) {
		log.Fatal("dir argument must be an absolute path, like /usr/local/nagios/etc/")
	}

	files, directories, err := watcher.GetAllInDirectory(baseDir)
	if err != nil {
		log.Fatalf("GetAllInDirectory: %v", err)
	}

	differential, err := watcher.NewDifferential([]string{".swp"}, files, http.DefaultClient, *address, *token)
	if err != nil {
		log.Fatalf("NewDifferential: %v", err)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		<-interrupt

		// We received an interrupt signal, shut down.
		log.Println("Bye")
	}()

	if err := watcher.WatchDirectories(directories, differential, done); err != nil {
		log.Panic(err)
	}
}
