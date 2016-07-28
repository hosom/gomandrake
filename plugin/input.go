/*

*/

package plugin

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/natefinch/pie"
)

type Input struct {
	Name		string
	client		*rpc.Client
}

func NewInput(name string) *Input {
	log.SetPrefix(fmt.Sprintf("[%s] "))

	return &Input{name, pie.NewConsumer()}
}

func (i Input) Analyze(fname string) {
	result := ""
	err := i.client.Call("mandrake.Analyze", fname, &result)
	log.Printf("File sent to analysis pipeline: %s", fname)
	if err != nil {
		log.Printf("Error occurred: %s", err)
	}
}