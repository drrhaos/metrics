// Package handlers ручки для добавления и получения метрик.
package handlers

import (
	"metrics/internal/server/configure"
)

// MetricsHandler хранит информацию о доступных ручках.
type MetricsHandler struct {
	cfg *configure.Config
}

// NewMetricHandler инициализирует новый объект типа MetricsHandler.
func NewMetricHandler(cfg *configure.Config) *MetricsHandler {
	return &MetricsHandler{
		cfg: cfg,
	}
}
