/*

*/
package mandrake

import (
	"log"
	"golang.org/x/exp/inotify"
	"github.com/hosom/gomandrake/config"
	//"github.com/hosom/gomandrake/plugin"
)

type Mandrake struct {
	AnalysisPipeline	chan string
	MonitoredDirectory	string
}

func NewMandrake(c config.Config) (*Mandrake, error) {
	return &Mandrake{make(chan string), c.MonitoredDirectory}, nil
}

func (m Mandrake) ListenAndServe() {
	go m.DispatchAnalysis()
	m.Monitor()
}

func (m Mandrake) DispatchAnalysis() {	
	for fpath := range m.AnalysisPipeline {
		log.Printf("%s", fpath)
	}
}

func (m Mandrake) Monitor() {
	log.Println("starting inotify watcher")
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("adding watcher to %s directory", m.MonitoredDirectory)
	err = watcher.AddWatch(m.MonitoredDirectory, inotify.IN_CLOSE_WRITE)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case ev := <- watcher.Event:
			m.AnalysisPipeline <- ev.Name
		case err := <- watcher.Error:
			log.Printf("inotify error: %s", err)
		}
	}
}