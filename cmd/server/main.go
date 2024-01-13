package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

const typeMetricCounter = "counter"
const typeMetricGauge = "gauge"

func main() {
	cfg, ok := readStartParams()

	if !ok {
		flag.PrintDefaults()
		os.Exit(0)
	}

	storage := MemStorage{
		counter: make(map[string]int64),
		gauge:   make(map[string]float64),
	}

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		getNameMetricsHandler(w, r, storage)
	})
	r.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", func(w http.ResponseWriter, r *http.Request) {
		updateMetricHandler(w, r, storage)
	})
	r.Get("/value/{typeMetric}/{nameMetric}", func(w http.ResponseWriter, r *http.Request) {
		getMetricHandler(w, r, storage)
	})
	log.Fatal(http.ListenAndServe(cfg.Address, r))
}
