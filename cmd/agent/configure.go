package main

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
}

func readStartParams() (Config, bool) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err == nil {
		log.Println("load environment")
	}

	address := flag.String("a", "127.0.0.1:8080", "Net address endpoint host:port")
	reportInterval := flag.Int64("r", 10, "Report interval integer sec > 0")
	pollInterval := flag.Int64("p", 2, "Pool interval integer sec > 0")
	flag.Parse()
	if cfg.Address == "" {
		cfg.Address = *address
	}

	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = *reportInterval
	}

	if cfg.PollInterval == 0 {
		cfg.PollInterval = *pollInterval
	}

	if cfg.PollInterval <= 0 || cfg.ReportInterval <= 0 {
		flag.PrintDefaults()
		return cfg, false
	} else {
		return cfg, true
	}
}
