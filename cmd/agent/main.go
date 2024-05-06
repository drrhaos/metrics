package main

import (
	"flag"
	"fmt"
	"os"

	"metrics/internal/agent/configure"
	"metrics/internal/agent/services"
	"metrics/internal/logger"
)

const flagLogLevel = "info"

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)

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
