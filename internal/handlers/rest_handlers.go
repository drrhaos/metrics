package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"metrics/internal/logger"
	"metrics/internal/store"
	_ "net/http/pprof"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (mh *MetricsHandler) UpdateMetricJSONHandler(res http.ResponseWriter, req *http.Request, storage *store.StorageContext) {
	if storage == nil {
		panic("Storage nil")
	}

	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
	defer cancel()

	var metrics Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	typeMetric := metrics.MType
	nameMetric := metrics.ID

	var ok bool
	var respMetrics Metrics
	respMetrics.ID = nameMetric
	respMetrics.MType = typeMetric

	switch typeMetric {
	case typeMetricCounter:
		if metrics.Delta == nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		ok = storage.UpdateCounter(ctx, nameMetric, *metrics.Delta)
		curDelta, exist := storage.GetCounter(ctx, nameMetric)
		if exist {
			curValue := float64(curDelta)
			respMetrics.Value = &curValue
		}

	case typeMetricGauge:
		if metrics.Value == nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		ok = storage.UpdateGauge(ctx, nameMetric, *metrics.Value)
		curValue, exist := storage.GetGauge(ctx, nameMetric)
		if exist {
			respMetrics.Value = &curValue
		}
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(respMetrics)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = res.Write(resp)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if mh.cfg.StoreInterval == 0 {
		storage.SaveMetrics(mh.cfg.FileStoragePath)
	}

	res.WriteHeader(http.StatusOK)
}

func (mh *MetricsHandler) UpdatesMetricJSONHandler(res http.ResponseWriter, req *http.Request, storage *store.StorageContext) {
	if storage == nil {
		panic("Storage nil")
	}

	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
	defer cancel()

	var metrics []Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		logger.Log.Warn("Не удалось прочитать тело запроса")
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	var ok bool
	for _, met := range metrics {
		switch met.MType {
		case typeMetricCounter:
			ok = storage.UpdateCounter(ctx, met.ID, *met.Delta)
		case typeMetricGauge:
			ok = storage.UpdateGauge(ctx, met.ID, *met.Value)
		default:
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if !ok {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if mh.cfg.StoreInterval == 0 {
		storage.SaveMetrics(mh.cfg.FileStoragePath)
	}
	res.WriteHeader(http.StatusOK)
}

func (mh *MetricsHandler) GetMetricJSONHandler(res http.ResponseWriter, req *http.Request, storage *store.StorageContext) {
	if storage == nil {
		panic("Storage nil")
	}

	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
	defer cancel()

	var metrics Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	typeMetric := metrics.MType
	nameMetric := metrics.ID

	if typeMetric == "" || nameMetric == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	switch typeMetric {
	case typeMetricCounter:
		curDelta, exist := storage.GetCounter(ctx, nameMetric)
		if !exist {
			res.WriteHeader(http.StatusNotFound)
			return

		}
		metrics.Delta = &curDelta
	case typeMetricGauge:
		curGauge, exist := storage.GetGauge(ctx, nameMetric)
		if !exist {
			res.WriteHeader(http.StatusNotFound)
			return

		}
		metrics.Value = &curGauge
	}

	resp, err := json.Marshal(metrics)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(resp)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.WriteHeader(http.StatusOK)
}
