package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Script struct {
	// An identifying key for the script
	Key string `yaml:"key"`
	// The script to run on Warp10
	Script string `yaml:"script"`
	// The type of the value returned by the script
	Type string `yaml:"type"`
	// Value above which an action is required
	Threshold float64 `yaml:"threshold,omitempty"`
	// The action to take (either "http" or "email")
	Action string `yaml:"action"`
	// The action's target
	Target string `yaml:"target"`
}

type Config struct {
	// Full URL to Warp10's /exec
	Warp10Exec string `yaml:"warp10_exec"`
	// Warp10 read token
	ReadToken string `yaml:"token"`
	// WarpScripts to run, with an identifier as its key
	Scripts []Script `yaml:"scripts"`
}

func Load(filePath string) (cfg Config, err error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(content, &cfg)
	return
}
