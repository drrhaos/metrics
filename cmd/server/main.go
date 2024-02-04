package main

import (
	"flag"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/drrhaos/metrics/internal/logger"
	"github.com/drrhaos/metrics/internal/storage"
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
const urlGetPing = "/ping"
const urlUpdateMetricConst = "/update/{typeMetric}/{nameMetric}/{valueMetric}"
const urlUpdateMetricJSONConst = "/update/"
const urlGetMetricConst = "/value/{typeMetric}/{nameMetric}"
const urlGetMetricJSONConst = "/value/"

const flagLogLevel = "info"

var cfg Config
var database Database

func main() {
	ok := cfg.readStartParams()

	if !ok {
		flag.PrintDefaults()
		os.Exit(0)
	}

	storage := &storage.MemStorage{
		Counter: make(map[string]int64),
		Gauge:   make(map[string]float64),
		Mut:     sync.Mutex{},
	}

	if cfg.Restore {
		storage.LoadMetrics(cfg.FileStoragePath)
	}

	if err := logger.Initialize(flagLogLevel); err != nil {
		panic(err)
	}

	if err := database.Connect(cfg.DatabaseDsn); err != nil {
		logger.Log.Info("Не удалось подключиться к базе данных", zap.Error(err))
	}
	defer database.Close()

	if cfg.StoreInterval != 0 {
		go func() {
			for {
				time.Sleep(time.Duration(cfg.StoreInterval) * time.Second)
				storage.SaveMetrics(cfg.FileStoragePath)
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
	r.Get(urlGetPing, func(w http.ResponseWriter, r *http.Request) {
		getPing(w, r)
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

	if err := http.ListenAndServe(cfg.Address, r); err != nil {
		logger.Log.Fatal(err.Error())
	}
	storage.SaveMetrics(cfg.FileStoragePath)
}
