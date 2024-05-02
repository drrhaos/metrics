package main

import (
	"flag"
	"os"

	"metrics/internal/agent/configure"
	"metrics/internal/agent/services"
	"metrics/internal/logger"
)

const flagLogLevel = "info"

func main() {
	if err := logger.Initialize(flagLogLevel); err != nil {
		panic(err)
	}

	var cfg configure.Config
	ok := cfg.ReadConfig()

	if !ok {
		flag.PrintDefaults()
		os.Exit(0)
	}

	services.CollectMetrics(cfg)
}
