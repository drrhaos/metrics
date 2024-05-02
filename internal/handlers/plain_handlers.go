package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"metrics/internal/store"

	"github.com/go-chi/chi"
)

const form = `<html>
    <head>
    <title></title>
    </head>
    <body>
	<ul>
	%s
	</ul>
    </body>
</html>`

const (
	typeMetricCounter = "counter"
	typeMetricGauge   = "gauge"
	typeMetricConst   = "typeMetric"
	nameMetricConst   = "nameMetric"
	valueMetricConst  = "valueMetric"
)

// UpdateMetricHandler Обновляет значение метрики.
func (mh *MetricsHandler) UpdateMetricHandler(res http.ResponseWriter, req *http.Request, storage *store.StorageContext) {
	if storage == nil {
		panic("Storage nil")
	}

	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
	defer cancel()

	typeMetric := chi.URLParam(req, typeMetricConst)
	nameMetric := chi.URLParam(req, nameMetricConst)
	valueMetric := chi.URLParam(req, valueMetricConst)

	if nameMetric == "" || valueMetric == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	var ok bool
	switch typeMetric {
	case typeMetricCounter:
		valueIntMetric, err := strconv.ParseInt(valueMetric, 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		ok = storage.UpdateCounter(ctx, nameMetric, valueIntMetric)
	case typeMetricGauge:
		valueFloatMetric, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		ok = storage.UpdateGauge(ctx, nameMetric, valueFloatMetric)
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if mh.cfg.StoreInterval == 0 {
		storage.SaveMetrics(mh.cfg.FileStoragePath)
	}
	res.WriteHeader(http.StatusOK)
}

// GetMetricHandler возвращает текущее значение метрики в текстовом виде.
func (mh *MetricsHandler) GetMetricHandler(res http.ResponseWriter, req *http.Request, storage *store.StorageContext) {
	if storage == nil {
		panic("Storage nil")
	}

	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
	defer cancel()

	typeMetric := chi.URLParam(req, typeMetricConst)
	nameMetric := chi.URLParam(req, nameMetricConst)

	if typeMetric == "" || nameMetric == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	var ok bool
	var currentValue string
	switch typeMetric {
	case typeMetricCounter:
		var curDelta int64
		curDelta, ok = storage.GetCounter(ctx, nameMetric)
		currentValue = strconv.FormatInt(curDelta, 10)
	case typeMetricGauge:
		var curValue float64
		curValue, ok = storage.GetGauge(ctx, nameMetric)
		currentValue = strconv.FormatFloat(curValue, 'f', -1, 64)
	default:
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := res.Write([]byte(currentValue))
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.WriteHeader(http.StatusOK)
}

// GetNameMetricsHandler возвращает сохранные метрики.
func (mh *MetricsHandler) GetNameMetricsHandler(res http.ResponseWriter, req *http.Request, storage *store.StorageContext) {
	if storage == nil {
		panic("Storage nil")
	}

	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
	defer cancel()

	var list string
	counters, ok := storage.GetCounters(ctx)
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	for key, val := range counters {
		list += fmt.Sprintf("<li>%s: %d</li>", key, val)
	}

	gauges, ok := storage.GetGauges(ctx)
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	for key, val := range gauges {
		list += fmt.Sprintf("<li>%s: %f</li>", key, val)
	}
	formFull := fmt.Sprintf(form, list)
	res.Header().Set("Content-Type", "text/html")
	_, err := io.WriteString(res, formFull)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.WriteHeader(http.StatusOK)
}

// GetPing проверяет доступность базы данных.
func (mh *MetricsHandler) GetPing(res http.ResponseWriter, req *http.Request, storage *store.StorageContext) {
	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
	defer cancel()

	if !storage.Ping(ctx) {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
