package main

import (
	"net/http"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

var storage MemStorage

func main() {
	storage.counter = make(map[string]int64)
	storage.gauge = make(map[string]float64)
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, updateMetricHandler)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
