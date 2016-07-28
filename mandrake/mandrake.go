/*

*/
package mandrake

import (
	"log"
	
	"github.com/hosom/gomandrake/config"
	"github.com/hosom/gomandrake/plugin"
)

type Mandrake struct {
	AnalysisPipeline	chan string
	Input				string
}

func NewMandrake(c config.Config) (*Mandrake, error) {
	return &Mandrake{make(chan string), c.Input}, nil
}

func (m *Mandrake) ListenAndServe() {
	go plugin.CreateListenerAndServe(m.AnalysisPipelinem, m.Input)
}