// Модуль store предназначен для хранения метрик
package store

import (
	"context"
	_ "net/http/pprof"
)

// StorageInterface описывает набор методов которые должны реализовывать хранилища
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

// StorageContext содержит текущее хранилище
type StorageContext struct {
	storage StorageInterface
}

// SetStorage устанавливает хранилище
func (sc *StorageContext) SetStorage(storage StorageInterface) {
	sc.storage = storage
}

// SaveMetrics сохраняет метрики
func (sc *StorageContext) SaveMetrics(filePath string) bool {
	return sc.storage.SaveMetrics(filePath)
}

// LoadMetrics загружает метрики
func (sc *StorageContext) LoadMetrics(filePath string) bool {
	return sc.storage.LoadMetrics(filePath)
}

// UpdateCounter обновляет метрику counter
func (sc *StorageContext) UpdateCounter(ctx context.Context, nameMetric string, valueMetric int64) bool {
	return sc.storage.UpdateCounter(ctx, nameMetric, valueMetric)
}

// UpdateGauge обновляет метрику gauge
func (sc *StorageContext) UpdateGauge(ctx context.Context, nameMetric string, valueMetric float64) bool {
	return sc.storage.UpdateGauge(ctx, nameMetric, valueMetric)
}

// GetGauges возвращает метрики gauges
func (sc *StorageContext) GetGauges(ctx context.Context) (map[string]float64, bool) {
	return sc.storage.GetGauges(ctx)
}

// GetCounters возвращает метрики counters
func (sc *StorageContext) GetCounters(ctx context.Context) (map[string]int64, bool) {
	return sc.storage.GetCounters(ctx)
}

// GetCounter возвращает метрику counters
func (sc *StorageContext) GetCounter(ctx context.Context, nameMetric string) (currentValue int64, exists bool) {
	return sc.storage.GetCounter(ctx, nameMetric)
}

// GetGauge возвращает метрику gauge
func (sc *StorageContext) GetGauge(ctx context.Context, nameMetric string) (currentValue float64, exists bool) {
	return sc.storage.GetGauge(ctx, nameMetric)
}

// Ping проверяет доступность хранилища
func (sc *StorageContext) Ping(ctx context.Context) (exists bool) {
	return sc.storage.Ping(ctx)
}
