package main

import (
	"flag"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/drrhaos/metrics/internal/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

const typeMetricCounter = "counter"
const typeMetricGauge = "gauge"
const typeMetricConst = "typeMetric"
const nameMetricConst = "nameMetric"
const valueMetricConst = "valueMetric"

const urlGetMetricsConst = "/"
const urlUpdateMetricConst = "/update/{typeMetric}/{nameMetric}/{valueMetric}"
const urlUpdateMetricJSONConst = "/update/"
const urlGetMetricConst = "/value/{typeMetric}/{nameMetric}"
const urlGetMetricJSONConst = "/value/"

const flagLogLevel = "info"

var cfg Config

func main() {
	ok := cfg.readStartParams()

	if !ok {
		flag.PrintDefaults()
		os.Exit(0)
	}

	storage := &MemStorage{
		Counter: make(map[string]int64),
		Gauge:   make(map[string]float64),
		mut:     sync.Mutex{},
	}

	if cfg.Restore {
		storage.loadMetrics(cfg.FileStoragePath)
	}

	if err := logger.Initialize(flagLogLevel); err != nil {
		panic(err)
	}

	if cfg.StoreInterval != 0 {
		go func() {
			for {
				time.Sleep(time.Duration(cfg.StoreInterval) * time.Second)
				storage.saveMetrics(cfg.FileStoragePath)
			}
		}()
	}

	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(gzipDecompressMiddleware)

	logger.Log.Info("Сервер запущен", zap.String("адрес", cfg.Address))

	r.Get(urlGetMetricsConst, func(w http.ResponseWriter, r *http.Request) {
		getNameMetricsHandler(w, r, storage)
	})
	r.Post(urlUpdateMetricConst, func(w http.ResponseWriter, r *http.Request) {
		updateMetricHandler(w, r, storage)
	})
	r.Post(urlUpdateMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		updateMetricJSONHandler(w, r, storage)
	})
	r.Get(urlGetMetricConst, func(w http.ResponseWriter, r *http.Request) {
		getMetricHandler(w, r, storage)
	})
	r.Post(urlGetMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		getMetricJSONHandler(w, r, storage)
	})

	err := http.ListenAndServe(cfg.Address, r)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
	storage.saveMetrics(cfg.FileStoragePath)
}
