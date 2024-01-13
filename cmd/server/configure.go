package main

import (
	"flag"
	"net/url"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func readStartParams() (Config, bool) {
	cfg := Config{}
	env.Parse(&cfg)
	address := flag.String("a", "127.0.0.1:8080", "Net address endpoint host:port")
	flag.Parse()
	if cfg.Address == "" {
		cfg.Address = *address
	}
	_, errURL := url.ParseRequestURI("http://" + cfg.Address)
	if errURL != nil {
		flag.PrintDefaults()
		return cfg, false
	} else {
		return cfg, true
	}
}
