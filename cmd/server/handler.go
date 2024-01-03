package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func updateMetricHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		typeMetric := chi.URLParam(req, "typeMetric")
		nameMetric := chi.URLParam(req, "nameMetric")
		valueMetric := chi.URLParam(req, "valueMetric")
		if typeMetric != "counter" && typeMetric != "gauge" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if nameMetric == "" || valueMetric == "" {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		if typeMetric == "counter" {
			valueIntMetric, err := strconv.ParseInt(valueMetric, 10, 64)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			storage.counter[nameMetric] += valueIntMetric
			fmt.Println(nameMetric, storage.counter[nameMetric])
		}

		if typeMetric == "gauge" {
			valueFloatMetric, err := strconv.ParseFloat(valueMetric, 10)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			storage.gauge[nameMetric] = valueFloatMetric
			fmt.Println(nameMetric, valueFloatMetric)
		}
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
	return
}

func getMetricHandler(rw http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "typeMetric")
	nameMetric := chi.URLParam(r, "nameMetric")
	fmt.Println(typeMetric, nameMetric)
	if typeMetric == "" || nameMetric == "" {
		http.Error(rw, "Not found", http.StatusNotFound)
		return
	}
	var currentValue string
	if typeMetric == "counter" {
		cur, ok := storage.counter[nameMetric]
		if ok {
			currentValue = strconv.FormatInt(cur, 10)
		} else {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	} else if typeMetric == "gauge" {
		cur, ok := storage.gauge[nameMetric]
		if ok {
			currentValue = strconv.FormatFloat(cur, 'f', -1, 64)
		} else {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	}
	rw.Write([]byte(currentValue))
}

// func getNameMetricsHandler(rw http.ResponseWriter, r *http.Request) {
// 	rw.Header().Set("Content-Type", "application/json")
// 	jsonStorage, _ := json.Marshal(storage.gauge)
// 	jsonStorageCounter, _ := json.Marshal(storage.counter)
// 	jsonStorage = append(jsonStorage, jsonStorageCounter...)
// 	rw.Write([]byte(jsonStorage))
// }

func getNameMetricsHandler(rw http.ResponseWriter, r *http.Request) {
	// rw.Header().Set("Content-Type", "application/json")
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
