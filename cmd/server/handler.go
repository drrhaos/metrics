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
		if metrics.Delta == nil {
			break
		}
		ok = storage.updateCounter(nameMetric, *metrics.Delta)
	case typeMetricGauge:
		if metrics.Value == nil {
			break
		}
		ok = storage.updateGauge(nameMetric, *metrics.Value)
	}
	res.Header().Set("Content-Type", "application/json")
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
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		respMetrics.Value = &curValue
	}

	resp, err := json.Marshal(respMetrics)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = res.Write(resp)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
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

func getMetricJSONHandler(res http.ResponseWriter, req *http.Request, storage *MemStorage) {
	if storage == nil {
		http.Error(res, "storage == nil", http.StatusNotFound)
		return
	}
	var metrics Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	typeMetric := metrics.MType
	nameMetric := metrics.ID

	if typeMetric == "" || nameMetric == "" {
		http.Error(res, "не задан тип или имя метрики", http.StatusNotFound)
		return
	}

	curStrValue, exist := storage.getMetric(typeMetric, nameMetric)
	if !exist {
		http.Error(res, "метрика не найдена", http.StatusNotFound)
		return

	}
	switch typeMetric {
	case typeMetricCounter:
		curValue, err := strconv.ParseInt(curStrValue, 10, 64)
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}
		metrics.Delta = &curValue
	case typeMetricGauge:
		curValue, err := strconv.ParseFloat(curStrValue, 64)
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}
		metrics.Value = &curValue
	}

	resp, err := json.Marshal(metrics)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(resp)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
	}
}

func getMetricHandler(res http.ResponseWriter, req *http.Request, storage *MemStorage) {
	if storage == nil {
		http.Error(res, "storage == nil", http.StatusNotFound)
		return
	}
	typeMetric := chi.URLParam(req, typeMetricConst)
	nameMetric := chi.URLParam(req, nameMetricConst)

	if typeMetric == "" || nameMetric == "" {
		http.Error(res, "не задан тип или имя метрики", http.StatusNotFound)
		return
	}

	currentValue, ok := storage.getMetric(typeMetric, nameMetric)
	if ok {
		_, err := res.Write([]byte(currentValue))
		if err != nil {
			http.Error(res, "ошибка записи", http.StatusNotFound)
			return
		}
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}

func getNameMetricsHandler(res http.ResponseWriter, req *http.Request, storage *MemStorage) {
	if storage == nil {
		http.Error(res, "storage == nil", http.StatusNotFound)
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
	_, err := io.WriteString(res, formFull)
	if err != nil {
		http.Error(res, "ошибка записи", http.StatusNotFound)
		return
	}
}
