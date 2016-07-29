/*
Reading of configuration files.
*/
package config

import (
	"io/ioutil"
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
)

// Config is a json-decoded configuration for running mandrake
type Config struct {
	MonitoredDirectory 	string
	InputPaths 			[]string
	AnalyzerPaths 		[]string
	LoggerPaths			[]string
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

// GetInputs looks at the directories in the InputPaths and returns a list of
// strings containing the available input plugins.
func (c *Config) GetInputs() []string {

	var plugins []string

	for _, path := range c.InputPaths {
		files, _ := filepath.Glob(path + "/*")
		plugins = append(plugins, files...)
	}

	return plugins
}