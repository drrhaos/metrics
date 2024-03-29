package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/drrhaos/metrics/internal/database"
	"github.com/drrhaos/metrics/internal/logger"
	"github.com/drrhaos/metrics/internal/ramstorage"
	"github.com/drrhaos/metrics/internal/signature"
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
const urlUpdatesMetricJSONConst = "/updates/"
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

	if err := logger.Initialize(flagLogLevel); err != nil {
		panic(err)
	}

	stMetrics := &StorageContext{}

	if cfg.DatabaseDsn != "" {
		stMetrics.setStorage(database.NewDatabase(cfg.DatabaseDsn))
	} else {
		stMetrics.setStorage(ramstorage.NewStorage())
	}

	if cfg.Restore {
		stMetrics.LoadMetrics(cfg.FileStoragePath)
	}

	if cfg.StoreInterval != 0 {
		go func() {
			for {
				time.Sleep(time.Duration(cfg.StoreInterval) * time.Second)
				stMetrics.SaveMetrics(cfg.FileStoragePath)
			}
		}()
	}

	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(gzipDecompressMiddleware)
	r.Use(signature.CheckSignaturMiddleware(cfg.Key))
	r.Use(signature.AddSignatureMiddleware(cfg.Key))

	logger.Log.Info("Сервер запущен", zap.String("адрес", cfg.Address))

	r.Get(urlGetMetricsConst, func(w http.ResponseWriter, r *http.Request) {
		getNameMetricsHandler(w, r, stMetrics)
	})
	r.Get(urlGetPing, func(w http.ResponseWriter, r *http.Request) {
		getPing(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricConst, func(w http.ResponseWriter, r *http.Request) {
		updateMetricHandler(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		updateMetricJSONHandler(w, r, stMetrics)
	})
	r.Post(urlUpdatesMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		updatesMetricJSONHandler(w, r, stMetrics)
	})
	r.Get(urlGetMetricConst, func(w http.ResponseWriter, r *http.Request) {
		getMetricHandler(w, r, stMetrics)
	})
	r.Post(urlGetMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		getMetricJSONHandler(w, r, stMetrics)
	})

	if err := http.ListenAndServe(cfg.Address, r); err != nil {
		logger.Log.Fatal(err.Error())
	}
	stMetrics.SaveMetrics(cfg.FileStoragePath)
}
