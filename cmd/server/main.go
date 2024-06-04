package main

import (
	"fmt"
	"time"

	"metrics/internal/logger"
	"metrics/internal/server/configure"
	"metrics/internal/server/grpcmode"
	"metrics/internal/server/restmode"
	"metrics/internal/store"
	"metrics/internal/store/pg"
	"metrics/internal/store/ramstorage"

	_ "net/http/pprof"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

const flagLogLevel = "info"

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
		logger.Log.Panic("Error read config")
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

	if cfg.GRPC {
		grpcmode.Run(cfg, stMetrics)
	} else {
		restmode.Run(cfg, stMetrics)
	}
}
