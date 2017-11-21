package alert

import (
	"fmt"

	"metrics-alerting/config"
	"metrics-alerting/warp10"
)

func ProcessNumber(
	client warp10.Warp10Client,
	script config.Script,
	ms config.MailSettings,
) error {
	value, err := client.ReadNumber(script.Script)
	if err != nil {
		return err
	}

	if value < script.Threshold {
		// Nothing to alert about
		return nil
	}

	return alert(script, value, ms)
}

func ProcessBool(
	client warp10.Warp10Client,
	script config.Script,
	ms config.MailSettings,
) error {
	value, err := client.ReadBool(script.Script)
	if err != nil {
		return err
	}

	if value {
		return nil
	}

	return alert(script, value, ms)
}

func alert(
	script config.Script,
	result interface{},
	ms config.MailSettings,
) error {
	switch script.Action {
	case "http":
		return alertHttp(script, result)
	case "email":
		return alertEmail(script, result, ms)
	default:
		return fmt.Errorf("invalid action type: %s", script.Action)
	}
}
