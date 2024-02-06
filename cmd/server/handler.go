package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

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

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func updateMetricJSONHandler(res http.ResponseWriter, req *http.Request, storage *SwitchStorage) {
	if storage == nil {
		panic("Storage nil")
	}
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
		ok = storage.UpdateCounter(nameMetric, *metrics.Delta)
		curDelta, exist := storage.GetCounter(nameMetric)
		if exist {
			curValue := float64(curDelta)
			respMetrics.Value = &curValue
		}

	case typeMetricGauge:
		if metrics.Value == nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		ok = storage.UpdateGauge(nameMetric, *metrics.Value)
		curValue, exist := storage.GetGauge(nameMetric)
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

	if cfg.StoreInterval == 0 {
		storage.SaveMetrics(cfg.FileStoragePath)
	}

	res.WriteHeader(http.StatusOK)
}

func updateMetricHandler(res http.ResponseWriter, req *http.Request, storage *SwitchStorage) {
	if storage == nil {
		panic("Storage nil")
	}
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
		ok = storage.UpdateCounter(nameMetric, valueIntMetric)
	case typeMetricGauge:
		valueFloatMetric, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		ok = storage.UpdateGauge(nameMetric, valueFloatMetric)
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if cfg.StoreInterval == 0 {
		storage.SaveMetrics(cfg.FileStoragePath)
	}
	res.WriteHeader(http.StatusOK)
}

func getMetricJSONHandler(res http.ResponseWriter, req *http.Request, storage *SwitchStorage) {
	if storage == nil {
		panic("Storage nil")
	}
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
		curDelta, exist := storage.GetCounter(nameMetric)
		if !exist {
			res.WriteHeader(http.StatusNotFound)
			return

		}
		metrics.Delta = &curDelta
	case typeMetricGauge:
		curGauge, exist := storage.GetGauge(nameMetric)
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

func getMetricHandler(res http.ResponseWriter, req *http.Request, storage *SwitchStorage) {
	if storage == nil {
		panic("Storage nil")
	}
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
		curDelta, ok = storage.GetCounter(nameMetric)
		currentValue = strconv.FormatInt(curDelta, 10)
	case typeMetricGauge:
		var curValue float64
		curValue, ok = storage.GetGauge(nameMetric)
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

func getNameMetricsHandler(res http.ResponseWriter, req *http.Request, storage *SwitchStorage) {
	if storage == nil {
		panic("Storage nil")
	}
	var list string
	counters, ok := storage.GetCounters()
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	for key, val := range counters {
		list += fmt.Sprintf("<li>%s: %d</li>", key, val)
	}

	gauges, ok := storage.GetGauges()
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

func getPing(res http.ResponseWriter, req *http.Request, storage *SwitchStorage) {
	if storage.DB.IsClosed() {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !storage.DB.Ping() {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
