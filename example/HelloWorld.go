
package main

import (
	"fmt"
	"log"
	"encoding/json"
	"net/rpc/jsonrpc"

	"github.com/hosom/gomandrake/filemeta"
	"github.com/natefinch/pie"
)

type HelloWorld struct {}

func (HelloWorld) Analyze(fmeta string, response *string) error {
	var fm filemeta.FileMeta
	json.Unmarshal([]byte(fmeta), &fm)
	log.Printf("Received call for Hello with name: %q", fm.Filepath)

	*response = fmt.Sprintf("{\"Hello\":%q}", fm.Filepath)

	return nil
}

func main () {
	log.SetPrefix("[helloworld] ")
	p := pie.NewProvider()

	if err := p.RegisterName("HelloWorld", HelloWorld{}); err != nil {
		log.Fatalf("failed to register : %s", err)
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
}
