/*

*/

package analyzer

import (
	"net/rpc/jsonrpc"
	"path/filepath"

	"github.com/natefinch/pie"
)

type Analyzer struct {
	client		*rpc.Client
	method		string
}

func NewAnalyzer(path string) (*Analyzer, error) {

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, path)
	if err != nil {
		log.Fatalf("Error executing plugin: %s", err)
	}
	defer client.Close()

	method := fmt.Sprintf("%s.Analyze", filepath.Base(path))

	a := Analyzer{client, method}

}

func (a *Analyzer) Analyze(fpath string) (result string, err error) {
	err = p.client.Call(a.method, fpath, &result)
	return result, err
}