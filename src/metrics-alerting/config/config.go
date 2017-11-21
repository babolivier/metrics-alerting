package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type MailSettings struct {
	// Sender of the alert emails
	Sender string `yaml:"sender"`
	// Settings to connect to the mail server
	SMTP SMTPSettings `yaml:"smtp"`
}

type SMTPSettings struct {
	// Host of the mail server
	Host string `yaml:"host"`
	// Port of the mail server
	Port int `yaml:"port"`
	// Username used to authenticate on the mail server
	Username string `yaml:"username"`
	// Password used to authenticate on the mail server
	Password string `yaml:"password"`
}

type Script struct {
	// An identifying key for the script
	Key string `yaml:"key"`
	// The script to run on Warp10
	Script string `yaml:"script"`
	// The type of the value returned by the script
	Type string `yaml:"type"`
	// Value above which an action is required, only required if the type is
	// "number"
	Threshold float64 `yaml:"threshold,omitempty"`
	// The action to take (either "http" or "email")
	Action string `yaml:"action"`
	// The action's target
	Target string `yaml:"target"`
	// The labels that will be mentioned in the email subject, only required if
	// the action is "email"
	IdentifyingLabels []string `yaml:"identifying_labels,omitempty"`
}

type Config struct {
	// Settings to send email alerts, only required if the action of at least 1
	// script is "email"
	Mail MailSettings `yaml:"mail,omitempty"`
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
