package main

import (
	"flag"
	"fmt"

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
	}
}
