/*

*/

package main

import (
	"flag"
	"path/filepath"
	"os"
	"fmt"
	"log"

	"golang.org/x/exp/inotify"
	"github.com/hosom/gomandrake/plugin"
)

const (
	// VERSION is the version of the plugin
	VERSION = "0.0.1"
	// NAME is the name of the plugin
	NAME = "inotify"
)

type config struct {
	MonitorPath		string
}

func main() {
	config_path := flag.String("config", 
		filepath.FromSlash("/etc/mandrake/inotify.conf"),
		"configuration file")
	version := flag.Bool("version",
		false,
		"Output version and exit.")
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	var c *config
	configuration, err := plugin.ReadConfig(*config_path, c)
	if err != nil {
		log.Fatal(err)
	}

	c = configuration.(*config)

	i := plugin.NewInput(NAME)

	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.AddWatch(c.MonitorPath, inotify.IN_CLOSE_WRITE)

	for {
		select {
		case ev := <-watcher.Event:
			i.Analyze(ev.Name)
		case err := <-watcher.Error:
			log.Printf("Error: %s", err)
		}
	}
}