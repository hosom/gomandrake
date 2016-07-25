/*

Main entry point for the `mandrake` daemon. Loads specified configuration
and listens for Input plugins to pass files to be analyzed by Analyzer 
plugins. Analysis is logged by Logger plugins.

*/

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hosom/gomandrake/config"
)

const (
	// VERSION is the version string for mandrake
	VERSION = "0.0.1"
)

func main() {

	config_path := flag.String("config", filepath.FromSlash("/etc/mandrake.conf"),
								"configuration file")

	version := flag.Bool("version", false, "Output version and exit.")
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if configuration, err := config.ReadConfigFile(*config_path); err != nil {
		fmt.Println("An error occurred")
		os.Exit(1)
	} else {
		fmt.Println(configuration)
	}

	matches, _ := filepath.Glob("./analyzers/*")
	fmt.Println(matches)
}