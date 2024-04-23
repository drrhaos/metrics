package handlers

import (
	"metrics/internal/server/configure"
)

type MetricsHandler struct {
	cfg configure.Config
}

func NewMetricHandler(cfg configure.Config) *MetricsHandler {
	return &MetricsHandler{
		cfg: cfg,
	}
}
