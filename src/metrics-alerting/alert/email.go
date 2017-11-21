package alert

import (
	"fmt"

	"metrics-alerting/config"

	"gopkg.in/gomail.v2"
)

func alertEmail(
	script config.Script,
	result interface{},
	ms config.MailSettings,
) error {
	formatNumber := `
Script %s just exceeded its threshold of %f and now returns %f

Script:

%s
	`

	formatBool := `
Test for script %s and returned false instead of true

Script:

%s
	`

	var body, subject string
	switch script.Type {
	case "number":
		subject = fmt.Sprintf("Threshold exceeded for script %s", script.Key)
		body = fmt.Sprintf(
			formatNumber, script.Key, script.Threshold, result.(float64),
			script.Script,
		)
	case "bool":
		subject = fmt.Sprintf("Test for script %s failed", script.Key)
		body = fmt.Sprintf(formatBool, script.Key, script.Script)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", ms.Sender)
	m.SetHeader("To", ms.Recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(
		ms.SMTP.Host, ms.SMTP.Port, ms.SMTP.Username, ms.SMTP.Password,
	)
	return d.DialAndSend(m)
}
