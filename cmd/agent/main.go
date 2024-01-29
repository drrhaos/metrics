package main

import (
	"flag"
	"os"
)

const typeMetricCounter = "counter"
const typeMetricGauge = "gauge"
const randomValueName = "RandomValue"
const pollCountName = "PollCount"

const urlUpdateJSONConst = "http://%s/update/"

var nameGauges = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

func main() {

	cfg, ok := readStartParams()
	if !ok {
		flag.PrintDefaults()
		os.Exit(0)
	}
	collectMetrics(cfg)
}
