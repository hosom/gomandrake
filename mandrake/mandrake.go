/*

*/
package mandrake

import (
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
	plugin.CreateListenerAndServe(m.AnalysisPipeline, m.Input)
}