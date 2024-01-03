package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

var storage MemStorage

func main() {
	storage.counter = make(map[string]int64)
	storage.gauge = make(map[string]float64)
	// mux := http.NewServeMux()
	// mux.HandleFunc(`/update/`, updateMetricHandler)
	// err := http.ListenAndServe(`:8080`, mux)
	// if err != nil {
	// 	panic(err)
	// }

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", getNameMetricsHandler)
		r.Route("/update", func(r chi.Router) {
			r.Route("/{typeMetric}", func(r chi.Router) {
				r.Route("/{nameMetric}", func(r chi.Router) {
					r.Post("/{valueMetric}", updateMetricHandler)
				})
			})
		})
		r.Route("/value", func(r chi.Router) {
			r.Route("/{typeMetric}", func(r chi.Router) {
				r.Get("/{nameMetric}", getMetricHandler)
			})
		})
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}
