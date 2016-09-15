/*
JSONLogger is a plugin for mandrake that provides basic logging services to
log records as JSON objects. One record is logged on each line.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc/jsonrpc"
	"os"
	"path/filepath"

	"github.com/natefinch/pie"
)

// JSONLogger is the API struct that describes the exposed API.
type JSONLogger struct {
	// LogFile is the open file that is being written to.
	LogFile *os.File
}

// Log exposes the API for the main mandrake daemon to send requests for
// analysis records to be logged.
func (j JSONLogger) Log(record string, response *string) error {
	_, err := j.LogFile.WriteString(fmt.Sprintf("%s\n", record))
	if err != nil {
		log.Println(err)
	}

	*response = "true"
	return nil
}

// main is the primary body of code for the plugin. Here command line
// arguments are parsed, the API is established and served.
func main() {
	log_path := flag.String("output", filepath.FromSlash("/var/log/mandrake.log"),
		"location for the plugin to send the log")
	flag.Parse()

	log_path_value := *log_path

	log.SetPrefix("[jsonlogger] ")
	log.Printf("Sending logs to %s", log_path_value)

	f, err := os.Create(log_path_value)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	p := pie.NewProvider()

	if err := p.RegisterName("JSONLogger", JSONLogger{f}); err != nil {
		log.Fatalf("failed to register : %s", err)
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
}
