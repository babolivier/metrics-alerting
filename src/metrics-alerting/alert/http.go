package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"metrics-alerting/config"
	"metrics-alerting/script_data"
)

type alertBody struct {
	Key    string            `json:"scriptKey"`
	Value  string            `json:"value"`
	Labels map[string]string `json:"labels"`
	Data   map[string]string `json:"data"`
}

func (a *Alerter) alertHttp(
	script config.Script,
	result interface{},
	labels map[string]string,
	data script_data.Data,
) error {
	var value string
	switch script.Type {
	case "number":
		value = strconv.FormatFloat(result.(float64), 'e', -1, 64)
	case "bool":
		value = strconv.FormatBool(result.(bool))
	}

	returnData := make(map[string]string)
	returnData[data.Key] = data.Value

	alert := alertBody{
		Key:    script.Key,
		Value:  value,
		Labels: labels,
		Data:   returnData,
	}

	body, err := json.Marshal(alert)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", script.Target, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"target %s returned non-200 status code %d", script.Target, resp.StatusCode,
		)
	}

	return err
}
