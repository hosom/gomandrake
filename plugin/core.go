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

func NewAnalyzerCaller(c config.PluginConfig) *AnalyzerCaller {
	a := AnalyzerCaller{}
	a.Path = c.Path
	a.Name = filepath.Base(a.Path)
	a.Args = c.Args
	a.MimeFilter = c.MimeFilter

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, a.Path, a.Args...)
	if err != nil {
		log.Fatalf("Error starting plugin: %s", a.Name)
	}

	a.client = client
	return &a
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
	client 		*rpc.Client
}

func NewLoggerCaller(fpath string) *LoggerCaller {
	name := filepath.Base(fpath)
	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, fpath)
	if err != nil {
		log.Fatalf("Error starting plugin: %s", name)
	}

	return &LoggerCaller{name, client}
}

// Log sends a json message to a logger plugin describing analysis that 
// has been completed.
func (l LoggerCaller) Log(msg string) (result string, err error) {
	log.Printf("Dispatching call to logger: %s", l.Name)
	err = l.client.Call(fmt.Sprintf("%s.Analyze", l.Name), msg, &result)
	
	return result, err
}