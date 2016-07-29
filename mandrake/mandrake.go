/*

*/
package mandrake

import (
	"log"

	"golang.org/x/exp/inotify"
	"github.com/hosom/gomandrake/config"
	//"github.com/hosom/gomandrake/plugin"
)

// Mandrake is a wrapper struct for the bulk of the application logic
type Mandrake struct {
	AnalysisPipeline	chan string
	MonitoredDirectory	string
}

// NewMandrake creates and returns a Mandrake struct utilizing a passed 
// parsed configuration file to create the correct fields.
func NewMandrake(c config.Config) Mandrake {
	return Mandrake{make(chan string), c.MonitoredDirectory}
}

// ListenAndServe starts the goroutines that perform all of the heavy lifting
// including Monitor() and DispatchAnalysis(). 
func (m Mandrake) ListenAndServe() {
	log.SetPrefix("[mandrake] ")
	go m.DispatchAnalysis()
	m.Monitor()
}

// DispatchAnalysis intelligently sends a new file to registered plugins so
// that it can be analyzed.
func (m Mandrake) DispatchAnalysis() {	
	for fpath := range m.AnalysisPipeline {
		log.Printf("%s", fpath)
	}
}

// Monitor uses inotify to monitor the MonitoredDirectory for IN_CLOSE_WRITE
// events. Files written to the MonitoredDirectory will be sent to the 
// analysis pipeline to be analyzed.
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