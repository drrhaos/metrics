package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func updateMetricHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		argsString := strings.TrimPrefix(req.URL.Path, "/update/")
		s := strings.Split(argsString, "/")
		if len(s) != 3 {
			res.WriteHeader(http.StatusNotFound)
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
			fmt.Println(nameMetric, valueFloatMetric)
		}
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
	return
}
