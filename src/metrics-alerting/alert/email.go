package alert

import (
	"encoding/json"
	"fmt"
	"strings"

	"metrics-alerting/config"
	"metrics-alerting/script_data"

	"gopkg.in/gomail.v2"
)

func (a *Alerter) alertEmail(
	script config.Script,
	result interface{},
	labels map[string]string,
	data script_data.Data,
) error {
	formatNumber := "Script \"%s\" just exceeded its threshold of %.2f and now returns %f"
	formatBool := "Test for script \"%s\" failed and returned false instead of true"

	var body, subject string
	switch script.Type {
	case "number", "series":
		subject = fmt.Sprintf(
			"Threshold exceeded for script \"%s\" %s%s", script.Key,
			getIdentifyingLabels(script, labels), getScriptData(data),
		)
		body = fmt.Sprintf(
			formatNumber, script.Key, script.Threshold, result.(float64),
		)
	case "bool":
		subject = fmt.Sprintf(
			"Test for script \"%s\" failed %s%s", script.Key,
			getIdentifyingLabels(script, labels), getScriptData(data),
		)
		body = fmt.Sprintf(formatBool, script.Key)
	}

	if labels != nil {
		jsonLabels, err := json.Marshal(labels)
		if err != nil {
			return err
		}
		body = fmt.Sprintf("%s\n\nLabels: %s", body, string(jsonLabels))
	}

	if len(data.Key) > 0 {
		dataMap := make(map[string]string)
		dataMap[data.Key] = data.Value
		jsonData, err := json.Marshal(dataMap)
		if err != nil {
			return err
		}
		body = fmt.Sprintf("%s\n\nData: %s", body, string(jsonData))
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
		if len(labels[label]) > 0 {
			identifyingLabels[label] = labels[label]
		}
	}

	labelsAsStrs := []string{}
	var labelAsStr string
	for key, value := range identifyingLabels {
		labelAsStr = key + ": " + value
		labelsAsStrs = append(labelsAsStrs, labelAsStr)
	}

	return "(" + strings.Join(labelsAsStrs, ", ") + ")"
}

func getScriptData(data script_data.Data) string {
	if len(data.Key) == 0 {
		return ""
	}

	return "(" + data.Key + "=" + data.Value + ")"
}
