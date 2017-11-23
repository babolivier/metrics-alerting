package alert

import (
	"fmt"

	"metrics-alerting/config"
	"metrics-alerting/script_data"

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
	data script_data.Data,
) error {
	switch script.Action {
	case "http":
		return a.alertHttp(script, result, labels, data)
	case "email":
		return a.alertEmail(script, result, labels, data)
	default:
		return fmt.Errorf("invalid action type: %s", script.Action)
	}
}
