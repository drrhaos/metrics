package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/drrhaos/metrics/internal/logger"
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

func updateMetricJsonHandler(res http.ResponseWriter, req *http.Request, storage *MemStorage) {
	var metrics Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	typeMetric := metrics.MType
	nameMetric := metrics.ID
	var ok bool
	switch typeMetric {
	case typeMetricCounter:
		ok = storage.updateCounter(nameMetric, *metrics.Delta)
	case typeMetricGauge:
		ok = storage.updateGauge(nameMetric, *metrics.Value)
	}

	if ok {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

func updateMetricHandler(res http.ResponseWriter, req *http.Request, storage *MemStorage) {
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
	}

	if ok {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

func getMetricHandler(rw http.ResponseWriter, r *http.Request, storage *MemStorage) {
	if storage == nil {
		http.Error(rw, "Not found", http.StatusNotFound)
		return
	}
	typeMetric := chi.URLParam(r, typeMetricConst)
	nameMetric := chi.URLParam(r, nameMetricConst)

	if typeMetric == "" || nameMetric == "" {
		http.Error(rw, "Not found", http.StatusNotFound)
		return
	}

	currentValue, ok := storage.getMetric(typeMetric, nameMetric)
	if ok {
		_, err := rw.Write([]byte(currentValue))
		if err != nil {
			logger.Log.Info("Ошибка записи")
		}
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func getNameMetricsHandler(rw http.ResponseWriter, r *http.Request, storage *MemStorage) {
	if storage == nil {
		http.Error(rw, "Not found", http.StatusNotFound)
		return
	}
	var list string
	for key, val := range storage.counter {
		list += fmt.Sprintf("<li>%s: %d</li>", key, val)
	}
	for key, val := range storage.gauge {
		list += fmt.Sprintf("<li>%s: %f</li>", key, val)
	}
	formFull := fmt.Sprintf(form, list)
	_, err := io.WriteString(rw, formFull)
	if err != nil {
		logger.Log.Info("Ошибка записи")
	}
}
