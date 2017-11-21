package main

import (
	"flag"
	"fmt"
	"log"

	"metrics-alerting/alert"
	"metrics-alerting/config"
	"metrics-alerting/warp10"
)

var (
	configPath = flag.String("config", "config.yaml", "The path to the config file. For more information, see the config file in this repository.")
)

func main() {
	flag.Parse()

	cfg, _ := config.Load(*configPath)
	client := warp10.Warp10Client{
		ExecEndpoint: cfg.Warp10Exec,
		ReadToken:    cfg.ReadToken,
	}

	for _, script := range cfg.Scripts {
		var err error
		switch script.Type {
		case "number":
			err = alert.ProcessNumber(client, script)
			break
		case "bool":
			err = alert.ProcessBool(client, script)
			break
		default:
			err = fmt.Errorf("invalid return type: %s", script.Type)
		}

		if err != nil {
			log.Fatal(err)
		}
	}
}
