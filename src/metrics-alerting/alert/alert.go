package alert

import (
	"fmt"

	"metrics-alerting/config"

	"gopkg.in/gomail.v2"
)

type Alerter struct {
	Dialer *gomail.Dialer
	Sender string
}

func (a *Alerter) Alert(
	script config.Script,
	result interface{},
	labels map[string]string,
) error {
	switch script.Action {
	case "http":
		return a.alertHttp(script, result, labels)
	case "email":
		return a.alertEmail(script, result, labels)
	default:
		return fmt.Errorf("invalid action type: %s", script.Action)
	}
}
