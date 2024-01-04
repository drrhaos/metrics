package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

var storage MemStorage

func main() {
	endpoint := flag.String("a", "127.0.0.1:8080", "Net address endpoint host:port")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		endpoint = &envRunAddr
	}

	storage.counter = make(map[string]int64)
	storage.gauge = make(map[string]float64)
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
	log.Fatal(http.ListenAndServe(*endpoint, r))
}
