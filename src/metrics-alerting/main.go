package main

import (
	"flag"
	"fmt"

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

	cfg, _ := config.Load(*configPath)
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
		var err error
		switch script.Type {
		case "number":
			err = process.ProcessNumber(client, script, alerter)
			break
		case "bool":
			err = process.ProcessBool(client, script, alerter)
			break
		case "series":
			err = process.ProcessSeries(client, script, alerter)
			break
		default:
			err = fmt.Errorf("invalid return type: %s", script.Type)
		}

		if err != nil {
			logrus.Error(err)
		}
	}
}
