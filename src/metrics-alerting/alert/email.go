package alert

import (
	"encoding/json"
	"fmt"
	"strings"

	"metrics-alerting/config"

	"gopkg.in/gomail.v2"
)

func (a *Alerter) alertEmail(
	script config.Script,
	result interface{},
	labels map[string]string,
) error {
	formatNumber := "Script %s just exceeded its threshold of %.2f and now returns %f"
	formatBool := "Test for script %s and returned false instead of true"

	var body, subject string
	switch script.Type {
	case "number", "series":
		subject = fmt.Sprintf(
			"Threshold exceeded for script %s %s", script.Key,
			getIdentifyingLabels(script, labels),
		)
		body = fmt.Sprintf(
			formatNumber, script.Key, script.Threshold, result.(float64),
		)
	case "bool":
		subject = fmt.Sprintf(
			"Test for script %s failed %s", script.Key,
			getIdentifyingLabels(script, labels),
		)
		body = fmt.Sprintf(formatBool, script.Key)
	}

	if labels != nil {
		jsonLabels, err := json.Marshal(labels)
		if err != nil {
			return err
		}
		body = fmt.Sprintf("%s\n\nLabels: %+v", body, string(jsonLabels))
	}

	body = fmt.Sprintf("%s\n\nScript:\n%s", body, script.Script)

	m := gomail.NewMessage()
	m.SetHeader("From", a.Sender)
	m.SetHeader("To", script.Target)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	return a.Dialer.DialAndSend(m)
}

func getIdentifyingLabels(
	script config.Script,
	labels map[string]string,
) string {
	if len(script.IdentifyingLabels) == 0 {
		return ""
	}

	identifyingLabels := make(map[string]string)
	for _, label := range script.IdentifyingLabels {
		identifyingLabels[label] = labels[label]
	}

	labelsAsStrs := []string{}
	var labelAsStr string
	for key, value := range identifyingLabels {
		labelAsStr = key + ": " + value
		labelsAsStrs = append(labelsAsStrs, labelAsStr)
	}

	return "(" + strings.Join(labelsAsStrs, ", ") + ")"
}
