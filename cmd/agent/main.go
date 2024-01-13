package main

import (
	"flag"
	"os"
)

func main() {
	cfg, ok := readStartParams()
	if !ok {
		flag.PrintDefaults()
		os.Exit(0)
	}
	collectMetrics(cfg)
}
