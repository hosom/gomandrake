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
)

// AnalyzerCaller is a wrapper specifically intended to be utilized for 
// wrapping Analyzer plugins.
type AnalyzerCaller struct {
	Name		string
	client 		*rpc.Client
}

func NewAnalyzerCaller(fpath string) *AnalyzerCaller {
	name := filepath.Base(fpath)
	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, fpath)
	if err != nil {
		log.Fatalf("Error starting plugin: %s", name)
	}

	defer client.Close()

	return &AnalyzerCaller{name, client}
}

// Analyze sends a filepath location to an analyzer plugin for the plugin
// to perform analysis on.
func  (a AnalyzerCaller) Analyze(fpath string) (result string, err error) {
	log.Printf("Dispatching call to analyzer: %s", a.Name)
	err = a.client.Call(fmt.Sprintf("%s.Analyze", a.Name), fpath, &result)
	
	return result, err
}

// LoggerCaller is a wrapper specifically intended to be utilized for 
// wrapping Logger plugins.
type LoggerCaller struct {
	Name 		string
	client 		*rpc.Client
}

func NewLoggerCaller(fpath string) *LoggerCaller {
	name := filepath.Base(fpath)
	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, fpath)
	if err != nil {
		log.Fatalf("Error starting plugin: %s", name)
	}

	defer client.Close()

	return &LoggerCaller{name, client}
}

// Log sends a json message to a logger plugin describing analysis that 
// has been completed.
func (l LoggerCaller) Log(msg string) (result string, err error) {
	log.Printf("Dispatching call to logger: %s", l.Name)
	err = l.client.Call(fmt.Sprintf("%s.Analyze", l.Name), msg, &result)
	
	return result, err
}

// PluginListener is a struct wrapper for listeners intended to be utilized
// in conjunction with Input plugins
type PluginListener struct {
	pipeline	chan string
}

func NewPluginListener(c chan string) *PluginListener {
	return &PluginListener{c}
}

func CreateListenerAndServe(c chan string, fpath string) {
	l := NewPluginListener(c)
	s, err := pie.StartConsumer(os.Stderr, fpath)
	if err != nil {
		log.Fatalf("Failed to initialize input %s: %s", fpath, err)
	}

	if err := s.RegisterName("mandrake", l); err != nil {
		log.Fatalf("Failed to register mandrake name: %s", err)
	}
	s.ServeCodec(jsonrpc.NewServerCodec)
}

// Analyze sends a file path into the analysis pipeline for it to be analyzed
func (p PluginListener) Analyze(fpath string, response *string) error {
	log.Printf("Request for analysis received: %s", fpath)
	p.pipeline <- fpath
	
	return nil
}