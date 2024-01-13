package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi"
)

const typeMetricCounter = "counter"
const typeMetricGauge = "gauge"
const typeMetricConst = "typeMetric"
const nameMetricConst = "nameMetric"
const valueMetricConst = "valueMetric"

const urlGetMetricsConst = "/"
const urlUpdateMetricsConst = "/update/{typeMetric}/{nameMetric}/{valueMetric}"
const urlGetMetricConst = "/value/{typeMetric}/{nameMetric}"

func main() {
	cfg, ok := readStartParams()

	if !ok {
		flag.PrintDefaults()
		os.Exit(0)
	}

	storage := &MemStorage{
		counter: make(map[string]int64),
		gauge:   make(map[string]float64),
		mut:     sync.Mutex{},
	}

	r := chi.NewRouter()

	r.Get("urlGetMetricsConst", func(w http.ResponseWriter, r *http.Request) {
		getNameMetricsHandler(w, r, storage)
	})
	r.Post(urlUpdateMetricsConst, func(w http.ResponseWriter, r *http.Request) {
		updateMetricHandler(w, r, storage)
	})
	r.Get(urlGetMetricConst, func(w http.ResponseWriter, r *http.Request) {
		getMetricHandler(w, r, storage)
	})
	log.Fatal(http.ListenAndServe(cfg.Address, r))
}
