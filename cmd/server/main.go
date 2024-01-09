package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func main() {

	endpoint := flag.String("a", "127.0.0.1:8080", "Net address endpoint host:port")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		endpoint = &envRunAddr
	}
	var storage MemStorage
	storage.init()

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		getNameMetricsHandler(w, r, storage)
	})
	r.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", func(w http.ResponseWriter, r *http.Request) {
		updateMetricHandler(w, r, storage)
	})
	r.Get("/value/{typeMetric}/{nameMetric}", func(w http.ResponseWriter, r *http.Request) {
		getMetricHandler(w, r, storage)
	})
	log.Fatal(http.ListenAndServe(*endpoint, r))
}
