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

func updateMetricJSONHandler(res http.ResponseWriter, req *http.Request, storage *MemStorage) {
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
	switch typeMetric {
	case typeMetricCounter:
		if metrics.Delta == nil {
			break
		}
		ok = storage.updateCounter(nameMetric, *metrics.Delta)
	case typeMetricGauge:
		if metrics.Value == nil {
			break
		}
		ok = storage.updateGauge(nameMetric, *metrics.Value)
	default:
		res.WriteHeader(http.StatusNotFound)
		return
	}
	if ok {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}

	var respMetrics Metrics
	curStrValue, exist := storage.getMetric(typeMetric, nameMetric)
	if exist {
		respMetrics.ID = nameMetric
		respMetrics.MType = typeMetric
		curValue, err := strconv.ParseFloat(curStrValue, 64)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		respMetrics.Value = &curValue
	}

	resp, err := json.Marshal(respMetrics)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = res.Write(resp)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
	}

	if cfg.StoreInterval == 0 {
		storage.saveMetrics(cfg.FileStoragePath)
	}
}

func updateMetricHandler(res http.ResponseWriter, req *http.Request, storage *MemStorage) {
	if storage == nil {
		panic("Storage nil")
	}
	typeMetric := chi.URLParam(req, typeMetricConst)
	nameMetric := chi.URLParam(req, nameMetricConst)
	valueMetric := chi.URLParam(req, valueMetricConst)

	if typeMetric != typeMetricCounter && typeMetric != typeMetricGauge {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if nameMetric == "" || valueMetric == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	var ok bool
	switch typeMetric {
	case typeMetricCounter:
		valueIntMetric, err := strconv.ParseInt(valueMetric, 10, 64)
		if err != nil {
			break
		}
		ok = storage.updateCounter(nameMetric, valueIntMetric)
	case typeMetricGauge:
		valueFloatMetric, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			break
		}
		ok = storage.updateGauge(nameMetric, valueFloatMetric)
	default:
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if ok {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}

	if cfg.StoreInterval == 0 {
		storage.saveMetrics(cfg.FileStoragePath)
	}
}

func getMetricJSONHandler(res http.ResponseWriter, req *http.Request, storage *MemStorage) {
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

	curStrValue, exist := storage.getMetric(typeMetric, nameMetric)
	if !exist {
		res.WriteHeader(http.StatusNotFound)
		return

	}
	switch typeMetric {
	case typeMetricCounter:
		curValue, err := strconv.ParseInt(curStrValue, 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		metrics.Delta = &curValue
	case typeMetricGauge:
		curValue, err := strconv.ParseFloat(curStrValue, 64)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		metrics.Value = &curValue
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
	}
}

func getMetricHandler(res http.ResponseWriter, req *http.Request, storage *MemStorage) {
	if storage == nil {
		panic("Storage nil")
	}
	typeMetric := chi.URLParam(req, typeMetricConst)
	nameMetric := chi.URLParam(req, nameMetricConst)

	if typeMetric == "" || nameMetric == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	currentValue, ok := storage.getMetric(typeMetric, nameMetric)
	if ok {
		_, err := res.Write([]byte(currentValue))
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}

func getNameMetricsHandler(res http.ResponseWriter, req *http.Request, storage *MemStorage) {
	if storage == nil {
		panic("Storage nil")
	}
	var list string
	for key, val := range storage.Counter {
		list += fmt.Sprintf("<li>%s: %d</li>", key, val)
	}
	for key, val := range storage.Gauge {
		list += fmt.Sprintf("<li>%s: %f</li>", key, val)
	}
	formFull := fmt.Sprintf(form, list)
	res.Header().Set("Content-Type", "text/html")
	_, err := io.WriteString(res, formFull)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
}
