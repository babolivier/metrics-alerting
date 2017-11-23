package process

import (
	"fmt"
	"regexp"
	"time"

	"metrics-alerting/alert"
	"metrics-alerting/config"
	"metrics-alerting/script_data"
	"metrics-alerting/warp10"
)

func Process(
	client warp10.Warp10Client,
	script config.Script,
	alerter alert.Alerter,
) error {
	var scriptData script_data.Data
	// TODO: Process more than one dataset
	for key, data := range script.ScriptData {
		scriptData.Key = key
		r, err := regexp.Compile("`" + key + "`")
		if err != nil {
			return err
		}
		match := r.Find([]byte(script.Script))
		if len(match) == 0 {
			return fmt.Errorf("no variable named %s in script %s", key, script.Key)
		}

		origScript := script.Script

		for _, el := range data {
			scriptData.Value = el
			filledScript := r.ReplaceAll([]byte(origScript), []byte(el))
			script.Script = string(filledScript)
			if err = dispatchType(client, script, alerter, scriptData); err != nil {
				return err
			}
		}
	}

	return nil
}

func dispatchType(
	client warp10.Warp10Client,
	script config.Script,
	alerter alert.Alerter,
	data script_data.Data,
) error {
	switch script.Type {
	case "number":
		return processNumber(client, script, alerter, data)
	case "bool":
		return processBool(client, script, alerter, data)
	case "series":
		return processSeries(client, script, alerter, data)
	}
	return fmt.Errorf("invalid return type: %s", script.Type)
}

func processNumber(
	client warp10.Warp10Client,
	script config.Script,
	alerter alert.Alerter,
	data script_data.Data,
) error {
	value, err := client.ReadNumber(script.Script)
	if err != nil {
		return err
	}

	return processFloat(value, script, alerter, nil, data)
}

func processBool(
	client warp10.Warp10Client,
	script config.Script,
	alerter alert.Alerter,
	data script_data.Data,
) error {
	value, err := client.ReadBool(script.Script)
	if err != nil {
		return err
	}

	if value {
		return nil
	}

	return alerter.Alert(script, value, nil, data)
}

func processSeries(
	client warp10.Warp10Client,
	script config.Script,
	alerter alert.Alerter,
	data script_data.Data,
) error {
	series, err := client.ReadSeriesOfNumbers(script.Script)

	for _, serie := range series {
		if !isRecentEnough(serie.Datapoints[0]) {
			// If the serie hasn't been active in the last 10min, don't consider
			// it
			// TODO: If the serie was active at the previous run, send an alert
			continue
		}

		// Remove useless ".app" label
		_, ok := serie.Labels[".app"]
		if ok {
			delete(serie.Labels, ".app")
		}
		// TODO: Currently we only process the most recent point.
		// If the datapoint is above the threshold, we should crawl back in time
		// to find when the situation began so we can add info about time in the
		// alert.
		if err = processFloat(
			serie.Datapoints[0][1], script, alerter, serie.Labels, data,
		); err != nil {
			return err
		}
	}

	return nil
}

func processFloat(
	value float64,
	script config.Script,
	alerter alert.Alerter,
	labels map[string]string,
	data script_data.Data,
) error {
	if value < script.Threshold {
		// Nothing to alert about
		return nil
	}

	return alerter.Alert(script, value, labels, data)
}

func isRecentEnough(datapoint []float64) bool {
	// Allowed offset between a point and now is 10min
	allowedOffset := int64(600000000)

	now := time.Now().UnixNano() / 1000 // Current timestamp (seconds)

	return now-int64(datapoint[0]) <= allowedOffset
	return false
}
