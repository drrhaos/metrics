package store

import (
	"context"
	_ "net/http/pprof"
)

type StorageInterface interface {
	SaveMetrics(filePath string) bool
	LoadMetrics(filePath string) bool
	UpdateCounter(ctx context.Context, nameMetric string, valueMetric int64) bool
	UpdateGauge(ctx context.Context, nameMetric string, valueMetric float64) bool
	GetGauges(ctx context.Context) (map[string]float64, bool)
	GetCounters(ctx context.Context) (map[string]int64, bool)
	GetCounter(ctx context.Context, nameMetric string) (currentValue int64, exists bool)
	GetGauge(ctx context.Context, nameMetric string) (currentValue float64, exists bool)
	Ping(ctx context.Context) bool
}

type StorageContext struct {
	storage StorageInterface
}

func (sc *StorageContext) SetStorage(storage StorageInterface) {
	sc.storage = storage
}

func (sc *StorageContext) SaveMetrics(filePath string) bool {
	return sc.storage.SaveMetrics(filePath)
}

func (sc *StorageContext) LoadMetrics(filePath string) bool {
	return sc.storage.LoadMetrics(filePath)
}

func (sc *StorageContext) UpdateCounter(ctx context.Context, nameMetric string, valueMetric int64) bool {
	return sc.storage.UpdateCounter(ctx, nameMetric, valueMetric)
}

func (sc *StorageContext) UpdateGauge(ctx context.Context, nameMetric string, valueMetric float64) bool {
	return sc.storage.UpdateGauge(ctx, nameMetric, valueMetric)
}

func (sc *StorageContext) GetGauges(ctx context.Context) (map[string]float64, bool) {
	return sc.storage.GetGauges(ctx)
}

func (sc *StorageContext) GetCounters(ctx context.Context) (map[string]int64, bool) {
	return sc.storage.GetCounters(ctx)
}

func (sc *StorageContext) GetCounter(ctx context.Context, nameMetric string) (currentValue int64, exists bool) {
	return sc.storage.GetCounter(ctx, nameMetric)
}

func (sc *StorageContext) GetGauge(ctx context.Context, nameMetric string) (currentValue float64, exists bool) {
	return sc.storage.GetGauge(ctx, nameMetric)
}

func (sc *StorageContext) Ping(ctx context.Context) (exists bool) {
	return sc.storage.Ping(ctx)
}
