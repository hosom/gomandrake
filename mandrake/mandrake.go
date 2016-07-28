/*

*/
package mandrake

import (
	"log"
	//"golang.org/x/exp/inotify"
	"github.com/hosom/gomandrake/config"
	"github.com/hosom/gomandrake/inputs"
	//"github.com/hosom/gomandrake/plugin"
)

type Mandrake struct {
	AnalysisPipeline	chan string
	Input				string
}

func NewMandrake(c config.Config) (*Mandrake, error) {
	return &Mandrake{make(chan string), c.Input}, nil
}

func (m *Mandrake) ListenAndServe() {
	go m.DispatchAnalysis()
	i := inputs.INotify{"/tmp/", m.AnalysisPipeline}
	i.Monitor()
	//plugin.CreateListenerAndServe(m.AnalysisPipeline, m.Input)
}

func (m *Mandrake) DispatchAnalysis() {	
	for fpath := range m.AnalysisPipeline {
		log.Printf("%s", fpath)
	}
}