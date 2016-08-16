package main

import (
	"log"
	"fmt"
	"net/rpc/jsonrpc"

	"github.com/natefinch/pie"
)

type JSONLogger struct {}

func (JSONLogger) Log(record string, response *string) error {
	fmt.Println(record)

	*response = "true"
	return nil
}

func main() {
	log.SetPrefix("[jsonlogger] ")
	p := pie.NewProvider()

	if err := p.RegisterName("JSONLogger", JSONLogger{}); err != nil {
		log.Fatalf("failed to register : %s", err)
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
}