package config

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// MailSettings represent the settings used to send email alerts
type MailSettings struct {
	// Sender of the email alerts
	Sender string `yaml:"sender"`
	// SMTP represent the settings needed to connect to the mail server
	SMTP SMTPSettings `yaml:"smtp"`
}

// SMTPSettings represent the settings needed to connect to the mail server
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

// ScriptDataSource represent the configuration structure for providing external
// data to iterate over. Optional.
type ScriptDataSource struct {
	// Type of the data source (either "plain" or "file")
	Source string `yaml:"source"`
	// Key of the data
	Key string `yaml:"key"`
	// Data value (or 1-element slice containing the path to the file containing
	// the values)
	Value []string `yaml:"value"`
}

// Script represents an instance of a script.
type Script struct {
	// Key is an identifying key for the script
	Key string `yaml:"key"`
	// Script is the script to run on Warp10
	Script string `yaml:"script"`
	// Type of the value returned by the script
	Type string `yaml:"type"`
	// Threshold is the value above which an action is required, only required
	// if the type is "number"
	Threshold float64 `yaml:"threshold,omitempty"`
	// Action identifies the action to take (either "http" or "email")
	Action string `yaml:"action"`
	// Target is the action's target
	Target string `yaml:"target"`
	// IdentifyingLabels represents a list of labels that will be mentioned in
	// the email subject. Optional.
	IdentifyingLabels []string `yaml:"identifying_labels,omitempty"`
	// ScriptDataSource represents the source/value of the data to use in the script
	ScriptDataSource ScriptDataSource `yaml:"script_data,omitempty"`
	// ScriptData represents loaded data, which isn't directly filled by parsing
	// the configuration file, but rather by reading it from ScriptDataSource
	ScriptData map[string][]string
}

// Config represents the global configuration of the app.
type Config struct {
	// Mail represents the settings needed to send email alerts. Only required
	// if the action of at least 1 script is "email"
	Mail *MailSettings `yaml:"mail,omitempty"`
	// Warp10Exec represents the full URL to Warp10's /exec
	Warp10Exec string `yaml:"warp10_exec"`
	// ReadToken represents Warp10's read token
	ReadToken string `yaml:"token"`
	// Script represents the WarpScripts to run
	Scripts []Script `yaml:"scripts"`
}

// Load parses the configuration file and load external data if needed.
// Returns an error if something went wrong when reading or parsing the
// configuration file, or reading the external data file if any.
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
