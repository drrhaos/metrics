package main

import (
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

func updateMetricHandler(res http.ResponseWriter, req *http.Request, storage MemStorage) {
	typeMetric := chi.URLParam(req, "typeMetric")
	nameMetric := chi.URLParam(req, "nameMetric")
	valueMetric := chi.URLParam(req, "valueMetric")
	if typeMetric != typeMetricCounter && typeMetric != typeMetricGauge {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if nameMetric == "" || valueMetric == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if typeMetric == typeMetricCounter {
		valueIntMetric, err := strconv.ParseInt(valueMetric, 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.updateCounter(nameMetric, valueIntMetric)
	}

	if typeMetric == typeMetricGauge {
		valueFloatMetric, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.updateGauge(nameMetric, valueFloatMetric)
	}
	res.WriteHeader(http.StatusOK)
}

func getMetricHandler(rw http.ResponseWriter, r *http.Request, storage MemStorage) {
	typeMetric := chi.URLParam(r, "typeMetric")
	nameMetric := chi.URLParam(r, "nameMetric")

	if typeMetric == "" || nameMetric == "" {
		http.Error(rw, "Not found", http.StatusNotFound)
		return
	}
	var currentValue string

	switch typeMetric {
	case typeMetricCounter:
		cur, ok := storage.getCounter(nameMetric)
		if ok {
			currentValue = strconv.FormatInt(cur, 10)
		} else {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	case typeMetricGauge:
		cur, ok := storage.getGauge(nameMetric)
		if ok {
			currentValue = strconv.FormatFloat(cur, 'f', -1, 64)
		} else {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	}
	rw.Write([]byte(currentValue))
}

func getNameMetricsHandler(rw http.ResponseWriter, r *http.Request, storage MemStorage) {
	var list string
	for key, val := range storage.counter {
		list += fmt.Sprintf("<li>%s: %d</li>", key, val)
	}
	for key, val := range storage.gauge {
		list += fmt.Sprintf("<li>%s: %f</li>", key, val)
	}
	formFull := fmt.Sprintf(form, list)
	io.WriteString(rw, formFull)
}
