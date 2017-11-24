package alert

import (
	"fmt"

	"metrics-alerting/config"
	"metrics-alerting/script_data"

	"github.com/sirupsen/logrus"
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
	logFailure(script, data)

	switch script.Action {
	case "http":
		return a.alertHttp(script, result, labels, data)
	case "email":
		return a.alertEmail(script, result, labels, data)
	default:
		return fmt.Errorf("invalid action type: %s", script.Action)
	}
}

func logFailure(script config.Script, data script_data.Data) {
	var entry *logrus.Entry
	if len(data.Key) > 0 {
		entry = logrus.WithField(data.Key, data.Value)
	} else {
		entry = logrus.NewEntry(logrus.New())
	}

	entry.Infof("Test for script \"%s\" failed", script.Key)
}
