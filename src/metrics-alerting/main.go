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
)

func main() {
	flag.Parse()

	cfg := config.Config{}
	if err := cfg.Load(*configPath); err != nil {
		logrus.Panic(err)
	}
	client := warp10.Warp10Client{
		ExecEndpoint: cfg.Warp10Exec,
		ReadToken:    cfg.ReadToken,
	}
	dialer := gomail.NewDialer(
		cfg.Mail.SMTP.Host, cfg.Mail.SMTP.Port, cfg.Mail.SMTP.Username,
		cfg.Mail.SMTP.Password,
	)
	alerter := alert.Alerter{
		Dialer: dialer,
		Sender: cfg.Mail.Sender,
	}

	for _, script := range cfg.Scripts {
		if err := process.Process(client, script, alerter); err != nil {
			logrus.Error(err)
		}
	}
}
