package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

var storage MemStorage

func updateMetric(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		argsString := strings.TrimPrefix(req.URL.Path, "/update/")
		s := strings.Split(argsString, "/")
		if len(s) != 3 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		typeMetric := s[0]
		nameMetric := s[1]
		valueMetric := s[2]
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
			// fmt.Println(nameMetric, valueFloatMetric)
		}
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
	return
}

func main() {
	storage.counter = make(map[string]int64)
	storage.gauge = make(map[string]float64)
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, updateMetric)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
