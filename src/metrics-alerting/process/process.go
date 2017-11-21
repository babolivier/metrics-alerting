package process

import (
	"time"

	"metrics-alerting/alert"
	"metrics-alerting/config"
	"metrics-alerting/warp10"
)

func ProcessNumber(
	client warp10.Warp10Client,
	script config.Script,
	alerter alert.Alerter,
) error {
	value, err := client.ReadNumber(script.Script)
	if err != nil {
		return err
	}

	return processFloat(value, script, alerter, nil)
}

func ProcessBool(
	client warp10.Warp10Client,
	script config.Script,
	alerter alert.Alerter,
) error {
	value, err := client.ReadBool(script.Script)
	if err != nil {
		return err
	}

	if value {
		return nil
	}

	return alerter.Alert(script, value, nil)
}

func ProcessSeries(
	client warp10.Warp10Client,
	script config.Script,
	alerter alert.Alerter,
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
			serie.Datapoints[0][1], script, alerter, serie.Labels,
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
) error {
	if value < script.Threshold {
		// Nothing to alert about
		return nil
	}

	return alerter.Alert(script, value, labels)
}

func isRecentEnough(datapoint []float64) bool {
	// Allowed offset between a point and now is 10min
	allowedOffset := int64(600000000)

	now := time.Now().UnixNano() / 1000 // Current timestamp (seconds)

	return now-int64(datapoint[0]) <= allowedOffset
	return false
}
