/*

*/

package plugin


import (
	"fmt"
	"log"
	"os"
	"net/rpc"
	"net/rpc/jsonrpc"
	"path/filepath"

	"github.com/natefinch/pie"
	"github.com/hosom/gomandrake/config"
)

// AnalyzerCaller is a wrapper specifically intended to be utilized for 
// wrapping Analyzer plugins.
type AnalyzerCaller struct {
	Name 		string
	Path 		string
	Args 		[]string
	MimeFilter	[]string
	client 		*rpc.Client
}

// allTypeAnalyzer is a function to check if a plugin should be ran against
// all mimetypes. If this is the case, then we will make some behind the 
// scenes configuration optimizations to prevent multiple executions of the
// same plugin.
func allTypeAnalyzer(m []string) bool {
	for _, filter := range m {
		if filter == "all" {
			return true
		}
	}
	return false
}

func NewAnalyzerCaller(c config.AnalyzerConfig) AnalyzerCaller {
	a := AnalyzerCaller{}
	a.Path = c.Path
	a.Name = filepath.Base(a.Path)
	a.Args = c.Args
	if allTypeAnalyzer(c.MimeFilter) == true {
		a.MimeFilter = []string{"all"}
	} else {
		a.MimeFilter = c.MimeFilter
	}

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, a.Path, a.Args...)
	if err != nil {
		log.Fatalf("Error starting plugin: %s", a.Name)
	}

	a.client = client
	return a
}

// Analyze sends a filepath location to an analyzer plugin for the plugin
// to perform analysis on.
func  (a AnalyzerCaller) Analyze(fmeta string) (result string, err error) {
	log.Printf("Dispatching call to analyzer: %s", a.Name)
	err = a.client.Call(fmt.Sprintf("%s.Analyze", a.Name), fmeta, &result)
	
	return result, err
}

// LoggerCaller is a wrapper specifically intended to be utilized for 
// wrapping Logger plugins.
type LoggerCaller struct {
	Name 		string
	Path 		string
	Args		[]string
	client 		*rpc.Client
}

func NewLoggerCaller(c config.LoggerConfig) LoggerCaller {
	l := LoggerCaller{}
	l.Path = c.Path
	l.Name = filepath.Base(l.Path)
	l.Args = c.Args

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, l.Path, l.Args...)
	if err != nil {
		log.Fatalf("Error starting plugin: %s", l.Name)
	}

	l.client = client 
	return l
}

// Log sends a json message to a logger plugin describing analysis that 
// has been completed.
func (l LoggerCaller) Log(msg string) (result string, err error) {
	log.Printf("Dispatching call to logger: %s", l.Name)
	err = l.client.Call(fmt.Sprintf("%s.Analyze", l.Name), msg, &result)
	
	return result, err
}