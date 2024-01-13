package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

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

func updateMetricHandler(res http.ResponseWriter, req *http.Request, storage *MemStorage) {
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
	ok := storage.updateMetric(typeMetric, nameMetric, valueMetric)
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
	typeMetric := chi.URLParam(r, "typeMetric")
	nameMetric := chi.URLParam(r, "nameMetric")

	if typeMetric == "" || nameMetric == "" {
		http.Error(rw, "Not found", http.StatusNotFound)
		return
	}

	currentValue, ok := storage.getMetric(typeMetric, nameMetric)
	if ok {
		_, err := rw.Write([]byte(currentValue))
		if err != nil {
			log.Printf("Ошибка записи: %v", err)
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
		log.Printf("Ошибка записи: %v", err)
	}
}
