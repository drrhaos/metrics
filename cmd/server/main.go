package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"metrics/internal/handlers"
	"metrics/internal/logger"
	"metrics/internal/middlewares/cryptodata"
	"metrics/internal/middlewares/decompress"
	"metrics/internal/middlewares/signature"
	"metrics/internal/server/configure"
	"metrics/internal/store"
	"metrics/internal/store/pg"
	"metrics/internal/store/ramstorage"

	_ "net/http/pprof"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

const (
	urlGetMetricsConst        = "/"
	urlGetPing                = "/ping"
	urlUpdateMetricConst      = "/update/{typeMetric}/{nameMetric}/{valueMetric}"
	urlUpdateMetricJSONConst  = "/update/"
	urlUpdatesMetricJSONConst = "/updates/"
	urlGetMetricConst         = "/value/{typeMetric}/{nameMetric}"
	urlGetMetricJSONConst     = "/value/"
	flagLogLevel              = "info"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)

	err := logger.Initialize(flagLogLevel)
	if err != nil {
		panic(err)
	}

	var cfg configure.Config
	ok := cfg.ReadConfig()

	if !ok {
		os.Exit(0)
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

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	r.Use(logger.RequestLogger)
	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(decompress.GzipDecompressMiddleware)
	r.Use(cryptodata.DecryptMiddleware(cfg.CryptoKeyPath))
	r.Use(signature.CheckSignaturMiddleware(cfg.Key))
	r.Use(signature.AddSignatureMiddleware(cfg.Key))
	r.Mount("/debug", middleware.Profiler())

	logger.Log.Info("Сервер запущен", zap.String("адрес", cfg.Address))

	metricHandler := handlers.NewMetricHandler(&cfg)

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
	go func() {
		if err := http.ListenAndServe(cfg.Address, r); err != nil {
			logger.Log.Fatal(err.Error())
		}
		stMetrics.SaveMetrics(cfg.FileStoragePath)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	logger.Log.Info("Получен сигнал прерывания, начинается грейсфул шатдаун")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Error("Ошибка при выполнении грейсфул шатдауна", zap.Error(err))
	}

	logger.Log.Info("Сервер успешно остановлен")
}
