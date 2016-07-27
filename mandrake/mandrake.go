/*

*/
package mandrake

import (
	"log"
	"net/rpc/jsonrpc"
	"github.com/natefinch/pie"
)

type Mandrake struct {
	consumers		chan string
}

func NewMandrake() (*Mandrake, err) {
	return Mandrake{make(chan string)}
}

func (m *Mandrake) ListenAndServe() {
	p := pie.NewProvider()
	if err := p.RegisterName("mandrake", m); err != nil {
		log.Fatalf("Failed to register plugin: %s", err)
	}

	go p.ServeCodec(jsonrpc.NewServerCodec)
}

func (m *Mandrake) Analyze(fname string, response *string) error {
	log.Printf("Beginning analysis: %s", fname)

	m.consumers <- fname

	return nil
}