/*

*/
package mandrake

import (
	"log"
	"encoding/json"

	"golang.org/x/exp/inotify"
	"github.com/hosom/gomandrake/config"
	"github.com/hosom/gomandrake/filemeta"
	"github.com/hosom/gomandrake/plugin"
)

// Mandrake is a wrapper struct for the bulk of the application logic
type Mandrake struct {
	AnalysisPipeline	chan string
	MonitoredDirectory	string
	Analyzers			[]plugin.AnalyzerCaller
	AnalyzerFilter		map[string][]plugin.AnalyzerCaller
}


// NewMandrake creates and returns a Mandrake struct utilizing a passed 
// parsed configuration file to create the correct fields.
func NewMandrake(c config.Config) Mandrake {
	analyzers := []plugin.AnalyzerCaller{}
	filter := make(map[string][]plugin.AnalyzerCaller)
	for _, plug := range c.Analyzers {
		analyzer := plugin.NewAnalyzerCaller(plug)
		// Build a slice of all AnalyzerCaller structs
		analyzers = append(analyzers, analyzer)

		// Create a map to function as a mime_type filter for analyzers
		for _, mime := range analyzer.MimeFilter {
			filter[mime] = append(filter[mime], analyzer)
		}
	}

	return Mandrake{make(chan string), c.MonitoredDirectory, analyzers, filter}
}

// ListenAndServe starts the goroutines that perform all of the heavy lifting
// including Monitor() and DispatchAnalysis(). 
func (m Mandrake) ListenAndServe() {
	log.SetPrefix("[mandrake] ")
	log.Println(m.Analyzers[0])
	go m.DispatchAnalysis()
	m.Monitor()
}

// DispatchAnalysis intelligently sends a new file to registered plugins so
// that it can be analyzed.
func (m Mandrake) DispatchAnalysis() {	
	for fpath := range m.AnalysisPipeline {
		fmeta, err := filemeta.NewFileMeta(fpath)
		if err != nil {
			log.Println(err)
		}

		// Create JSON filemeta object to pass to plugins so that plugins
		// receive basic contextual information about the file.
		fs, err := json.Marshal(fmeta)

		var analysis []map[string]interface{}

		for _, analyzer := range m.AnalyzerFilter["all"] {
			result, err := analyzer.Analyze(string(fs))
			if err != nil {
				log.Print(err)
			}
			analysis = append(analysis, MapFromJSON(result))
		}

		for _, analyzer := range m.AnalyzerFilter[fmeta.Mime] {
			result, err := analyzer.Analyze(string(fs))
			if err != nil {
				log.Print(err)
			}
			analysis = append(analysis, MapFromJSON(result))
		}

		l, _ json.Marshal(analysis)
		log.Println(string(l))
		log.Println(string(fs))
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

// MapFromJSON accepts an anonymous JSON object as a string and returns the
// resulting Map
func MapFromJSON(s string) map[string]interface{} {
	log.Printf("Performing mapping with string: %s", s)
	var f interface{}
	json.Unmarshal([]byte(s), &f)
	m := f.(map[string]interface{})
	log.Printf("Mapping complete: %s", m)
	return m
}