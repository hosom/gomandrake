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
	Path 		string
	Args 		[]string
	MimeFilter	[]string
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