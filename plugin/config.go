/*

*/
package plugin

import (
	"io/ioutil"
	"bytes"
	"encoding/json"
	"fmt"
)

// Generic function for reading configuration files for plugins
// intended to make plugin development easier.
func ReadConfig(fpath string, config interface{}) (interface{}, error) {
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file %q: %v", fpath, err)
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	
	if err := dec.Decode(&config); err != nil {
		return nil, fmt.Errorf("could not decode config file %q: %v", fpath, err)
	}

	return config, nil
}