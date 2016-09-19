/*

 */
package mandrake

import (
	"encoding/json"
	"log"

	"github.com/hosom/gomandrake/config"
	"github.com/hosom/gomandrake/filemeta"
	"github.com/hosom/gomandrake/plugin"
	"golang.org/x/exp/inotify"
)

// Mandrake is a wrapper struct for the bulk of the application logic
type Mandrake struct {
	AnalysisPipeline   chan string
	LoggingPipeline    chan string
	MonitoredDirectory string
	Analyzers          []plugin.AnalyzerCaller
	AnalyzerFilter     map[string][]plugin.AnalyzerCaller
	Loggers            []plugin.LoggerCaller
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

	loggers := []plugin.LoggerCaller{}
	for _, plug := range c.Loggers {
		logger := plugin.NewLoggerCaller(plug)
		loggers = append(loggers, logger)
	}

	return Mandrake{make(chan string), make(chan string), c.MonitoredDirectory, analyzers, filter, loggers}
}

// ListenAndServe starts the goroutines that perform all of the heavy lifting
// including Monitor() and DispatchAnalysis().
func (m Mandrake) ListenAndServe() {
	log.SetPrefix("[mandrake] ")
	go m.DispatchLogging()
	go m.DispatchAnalysis()
	m.Monitor()
}

// DispatchAnalysis intelligently sends a new file to registered plugins so
// that it can be analyzed.
func (m Mandrake) DispatchAnalysis() {
	for fpath := range m.AnalysisPipeline {
		go m.Analysis(fpath)
	}
}

// Analysis is the method that kicks off all of the analysis plugins
// this is utilized so that each file can be analyzed in a goroutine
func (m Mandrake) Analysis(fpath string) {
	fmeta, err := filemeta.NewFileMeta(fpath)
	if err != nil {
		log.Println(err)
	}

	// Create JSON filemeta object to pass to plugins so that plugins
	// receive basic contextual information about the file.
	fs, err := json.Marshal(fmeta)
	// Finalize string form of JSON filemeta to pass to plugins
	fstring := string(fs)

	var analysis []map[string]interface{}

	for _, analyzer := range m.AnalyzerFilter["all"] {
		log.Printf("%s : dispatching call to analyzer %s", fpath, analyzer.Name)
		result, err := analyzer.Analyze(fstring)
		if err != nil {
			log.Print(err)
		}
		analysis = append(analysis, MapFromJSON(result))
	}

	for _, analyzer := range m.AnalyzerFilter[fmeta.Mime] {
		log.Printf("%s : dispatching call to analyzer %s", fpath, analyzer.Name)
		result, err := analyzer.Analyze(fstring)
		if err != nil {
			log.Print(err)
		}
		analysis = append(analysis, MapFromJSON(result))
	}

	report := MapFromJSON(fstring)
	report["analysis"] = analysis

	r, _ := json.Marshal(report)

	log.Printf("Analysis of %s complete", fpath)
	m.LoggingPipeline <- string(r)
	log.Printf("Analysis of %s sent to logging pipeline.", fpath)
}

// DispatchLogging sends the call to the Logger plugins to log the completed
// record of analysis performed by Mandrake
func (m Mandrake) DispatchLogging() {
	for record := range m.LoggingPipeline {
		for _, logger := range m.Loggers {
			_, err := logger.Log(record)
			if err != nil {
				log.Print(err)
			}
		}
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
		case ev := <-watcher.Event:
			m.AnalysisPipeline <- ev.Name
		case err := <-watcher.Error:
			log.Printf("inotify error: %s", err)
		}
	}
}

// MapFromJSON accepts an anonymous JSON object as a string and returns the
// resulting Map
func MapFromJSON(s string) map[string]interface{} {
	if s == "" {
		log.Println("Encountered invalid output from plugin.")
		return make(map[string]interface{})
	}

	var f interface{}
	json.Unmarshal([]byte(s), &f)
	m := f.(map[string]interface{})
	return m
}
