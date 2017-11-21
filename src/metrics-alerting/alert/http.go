package alert

import (
	"bytes"
	"encoding/json"
	"net/http"

	"metrics-alerting/config"
)

type alertBody struct {
	Key   string `json:"scriptKey"`
	Value string `json:"value"`
}

func alertHttp(script config.Script, result interface{}) error {
	alert := alertBody{
		Key:   script.Key,
		Value: result,
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
			"target %s returned non-200 status code %d", script.Target, resp.StatusCode
		)
	}

	return err
}
