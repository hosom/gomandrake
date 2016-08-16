
package main

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"

	"github.com/natefinch/pie"
)

type HelloWorld struct {}

func (HelloWorld) Analyze(fname string, response *string) error {
	log.Printf("Received call for Hello with name: %q", fname)

	*response = fmt.Sprintf("{\"Hello\":%q}", fname)

	return nil
}

func main () {
	log.SetPrefix("[helloworld ] ")
	p := pie.NewProvider()

	if err := p.RegisterName("HelloWorld", HelloWorld{}); err != nil {
		log.Fatalf("failed to register : %s", err)
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
}
