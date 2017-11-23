package config

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"

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

type ScriptDataSource struct {
	// Type of the data source (either "plain" or "file")
	Source string `yaml:"source"`
	// Key of the data
	Key string `yaml:"key"`
	// Data value (or 1-element slice containing the path to the file containing
	// the values)
	Value []string `yaml:"value"`
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
	// Source/value of the data to use in the script
	ScriptDataSource ScriptDataSource `yaml:"script_data,omitempty"`
	// Loaded data
	ScriptData map[string][]string
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

func (cfg *Config) Load(filePath string) (err error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return
	}

	return cfg.loadData()
}

func (cfg *Config) loadData() error {
	var line string
	var l []byte
	var isPrefix bool
	for i, script := range cfg.Scripts {
		script.ScriptData = make(map[string][]string)
		switch script.ScriptDataSource.Source {
		case "file":
			fp, err := os.Open(script.ScriptDataSource.Value[0])
			if err != nil {
				return err
			}
			reader := bufio.NewReader(fp)

			for true {
				isPrefix = true
				line = ""
				for isPrefix {
					l, isPrefix, err = reader.ReadLine()
					if err != nil && err != io.EOF {
						return err
					}
					line = line + string(l)
				}

				if err == io.EOF {
					break
				}

				// Prevent processing empty line at the end of file
				if len(line) > 0 {
					script.ScriptData[script.ScriptDataSource.Key] = append(
						script.ScriptData[script.ScriptDataSource.Key], line,
					)
				}
			}
			break
		case "plain":
			script.ScriptData[script.ScriptDataSource.Key] = script.ScriptDataSource.Value
			break
		default:
			return fmt.Errorf("invalid data source: %s")
		}

		cfg.Scripts[i] = script
	}

	return nil
}
