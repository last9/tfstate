package main

import (
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var configPath string

func setupCLI() {
	flag.StringVar(&configPath, "config", "", "Location of the yaml config file")
}

func getConfigPath() string {
	return configPath
}

func parseConfig(cfgpath string, t interface{}) (err error) {
	data, err := ioutil.ReadFile(cfgpath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, t); err != nil {
		return err
	}

	return
}
