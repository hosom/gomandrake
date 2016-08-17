package main

import (
	"flag"
	"path/filepath"
	"fmt"
	"log"
	"os"
	"net/rpc/jsonrpc"

	"github.com/natefinch/pie"
)

type JSONLogger struct {
	LogFile			*os.File
}

func (j JSONLogger) Log(record string, response *string) error {
	_, err := j.LogFile.WriteString(fmt.Sprintf("%s\n", record))
	if err != nil {
		log.Println(err)
	}

	*response = "true"
	return nil
}

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

	p := pie.NewProvider()

	if err := p.RegisterName("JSONLogger", JSONLogger{f}); err != nil {
		log.Fatalf("failed to register : %s", err)
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
}