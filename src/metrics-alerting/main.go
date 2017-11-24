package main

import (
	"flag"

	"metrics-alerting/alert"
	"metrics-alerting/config"
	"metrics-alerting/process"
	"metrics-alerting/warp10"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

var (
	configPath = flag.String("config", "config.yaml", "The path to the config file. For more information, see the config file in this repository.")
	logFile    = flag.String("log-file", "", "The path to the file logs should be directed to. If empty or not provided, logs will be output through standard output.")
)

func main() {
	flag.Parse()

	if err := logConfig(); err != nil {
		logrus.Warn(err)
	}

	cfg := config.Config{}
	if err := cfg.Load(*configPath); err != nil {
		logrus.Panic(err)
	}
	client := warp10.Warp10Client{
		ExecEndpoint: cfg.Warp10Exec,
		ReadToken:    cfg.ReadToken,
	}
	alerter := alert.Alerter{
		Dialer: nil,
		Sender: cfg.Mail.Sender,
	}

	for _, script := range cfg.Scripts {
		if script.Action == "email" && alerter.Dialer == nil {
			if cfg.Mail == nil {
				logrus.Errorf(
					"no configuration set for sending emails, ignoring script %s",
					script.Key,
				)
				continue
			}

			alerter.Dialer = gomail.NewDialer(
				cfg.Mail.SMTP.Host, cfg.Mail.SMTP.Port, cfg.Mail.SMTP.Username,
				cfg.Mail.SMTP.Password,
			)
		}

		if err := process.Process(client, script, alerter); err != nil {
			logrus.Error(err)
		}
	}
}
