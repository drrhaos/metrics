package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"metrics/internal/handlers"
	"metrics/internal/logger"
	"metrics/internal/middlewares/decompress"
	"metrics/internal/server/configure"
	"metrics/internal/signature"
	"metrics/internal/store"
	"metrics/internal/store/pg"
	"metrics/internal/store/ramstorage"

	_ "net/http/pprof"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

const urlGetMetricsConst = "/"
const urlGetPing = "/ping"
const urlUpdateMetricConst = "/update/{typeMetric}/{nameMetric}/{valueMetric}"
const urlUpdateMetricJSONConst = "/update/"
const urlUpdatesMetricJSONConst = "/updates/"
const urlGetMetricConst = "/value/{typeMetric}/{nameMetric}"
const urlGetMetricJSONConst = "/value/"

const flagLogLevel = "info"

func main() {
	cfg := configure.NewConfig()

	if cfg == nil {
		flag.PrintDefaults()
		os.Exit(0)
	}

	err := logger.Initialize(flagLogLevel)
	if err != nil {
		panic(err)
	}

	stMetrics := &store.StorageContext{}

	if cfg.DatabaseDsn != "" {
		stMetrics.SetStorage(pg.NewDatabase(cfg.DatabaseDsn))
	} else {
		stMetrics.SetStorage(ramstorage.NewStorage())
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
	r.Use(decompress.GzipDecompressMiddleware)
	r.Use(signature.CheckSignaturMiddleware(cfg.Key))
	r.Use(signature.AddSignatureMiddleware(cfg.Key))
	r.Mount("/debug", middleware.Profiler())

	logger.Log.Info("Сервер запущен", zap.String("адрес", cfg.Address))

	metricHandler := handlers.NewMetricHandler(cfg)

	r.Get(urlGetMetricsConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetNameMetricsHandler(w, r, stMetrics)
	})
	r.Get(urlGetPing, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetPing(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricHandler(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricJSONHandler(w, r, stMetrics)
	})
	r.Post(urlUpdatesMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdatesMetricJSONHandler(w, r, stMetrics)
	})
	r.Get(urlGetMetricConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricHandler(w, r, stMetrics)
	})
	r.Post(urlGetMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricJSONHandler(w, r, stMetrics)
	})

	if err := http.ListenAndServe(cfg.Address, r); err != nil {
		logger.Log.Fatal(err.Error())
	}
	stMetrics.SaveMetrics(cfg.FileStoragePath)
}
