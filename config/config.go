/*
Reading of configuration files.
*/
package config

import (
	"io/ioutil"
	"bytes"
	"encoding/json"
	"fmt"
)

// Config is a json-decoded configuration for running mandrake
type Config struct {
	InputPaths 		[]string
	AnalyzerPaths 	[]string
	LoggerPaths		[]string
}

// ReadConfigFile reads in the given JSON encoded configuration file and
// returns the Config object associated with the decoded configuration data.
func ReadConfigFile(filename string) (*Config, error) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read config file %q: %v", filename, err)
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	
	var out Config
	if err := dec.Decode(&out); err != nil {
		return nil, fmt.Errorf("could not decode config file %q: %v", filename, err)
	}

	return &out, nil
}