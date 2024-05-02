package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"metrics/internal/handlers"
	"metrics/internal/server/configure"
	"metrics/internal/store"
	"metrics/internal/store/mocks"
	"metrics/internal/store/ramstorage"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
)

const (
	urlGetMetricsConst        = "/"
	urlGetPing                = "/ping"
	urlUpdateMetricConst      = "/update/{typeMetric}/{nameMetric}/{valueMetric}"
	urlUpdateMetricJSONConst  = "/update/"
	urlUpdatesMetricJSONConst = "/updates/"
	urlGetMetricConst         = "/value/{typeMetric}/{nameMetric}"
	urlGetMetricJSONConst     = "/value/"
)

var cfg = configure.Config{}

func ExampleMetricsHandler_UpdateMetricHandler() {
	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	metricHandler := handlers.NewMetricHandler(&cfg)

	r := chi.NewRouter()

	r.Get(urlGetMetricsConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetNameMetricsHandler(w, r, stMetrics)
	})
	r.Get(urlGetPing, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetPing(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricHandler(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricJSONHandler(w, r, stMetrics)
	})
	r.Post(urlUpdatesMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdatesMetricJSONHandler(w, r, stMetrics)
	})
	r.Get(urlGetMetricConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricHandler(w, r, stMetrics)
	})
	r.Post(urlGetMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricJSONHandler(w, r, stMetrics)
	})

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/NumGC/11.1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	req = httptest.NewRequest(http.MethodPost, "/update/counter/PoolCounret/11", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp = w.Result()
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	// Output:
	// 200
	// 200
}

func ExampleMetricsHandler_GetMetricHandler() {
	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())
	stMetrics.UpdateCounter(context.Background(), "PoolCounter", 11)
	stMetrics.UpdateGauge(context.Background(), "NumGC", 11.1)

	metricHandler := handlers.NewMetricHandler(&cfg)

	r := chi.NewRouter()

	r.Get(urlGetMetricsConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetNameMetricsHandler(w, r, stMetrics)
	})
	r.Get(urlGetPing, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetPing(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricHandler(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricJSONHandler(w, r, stMetrics)
	})
	r.Post(urlUpdatesMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdatesMetricJSONHandler(w, r, stMetrics)
	})
	r.Get(urlGetMetricConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricHandler(w, r, stMetrics)
	})
	r.Post(urlGetMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricJSONHandler(w, r, stMetrics)
	})

	req := httptest.NewRequest(http.MethodGet, "/value/gauge/NumGC", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	req = httptest.NewRequest(http.MethodGet, "/value/counter/PoolCounter", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp = w.Result()
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	// Output:
	// 200
	// 200
}

func ExampleMetricsHandler_GetNameMetricsHandler() {
	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())
	stMetrics.UpdateCounter(context.Background(), "PoolCounter", 11)
	stMetrics.UpdateGauge(context.Background(), "NumGC", 11.1)

	metricHandler := handlers.NewMetricHandler(&cfg)

	r := chi.NewRouter()

	r.Get(urlGetMetricsConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetNameMetricsHandler(w, r, stMetrics)
	})
	r.Get(urlGetPing, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetPing(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricHandler(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricJSONHandler(w, r, stMetrics)
	})
	r.Post(urlUpdatesMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdatesMetricJSONHandler(w, r, stMetrics)
	})
	r.Get(urlGetMetricConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricHandler(w, r, stMetrics)
	})
	r.Post(urlGetMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricJSONHandler(w, r, stMetrics)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	// Output:
	// 200
}

func ExampleMetricsHandler_GetPing() {
	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())
	stMetrics.UpdateCounter(context.Background(), "PoolCounter", 11)
	stMetrics.UpdateGauge(context.Background(), "NumGC", 11.1)

	metricHandler := handlers.NewMetricHandler(&cfg)

	r := chi.NewRouter()

	r.Get(urlGetMetricsConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetNameMetricsHandler(w, r, stMetrics)
	})
	r.Get(urlGetPing, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetPing(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricHandler(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricJSONHandler(w, r, stMetrics)
	})
	r.Post(urlUpdatesMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdatesMetricJSONHandler(w, r, stMetrics)
	})
	r.Get(urlGetMetricConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricHandler(w, r, stMetrics)
	})
	r.Post(urlGetMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricJSONHandler(w, r, stMetrics)
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	// Output:
	// 500
}

func ExampleMetricsHandler_UpdateMetricJSONHandler() {
	deltaCur := int64(11)
	metricCounter := handlers.Metrics{
		ID:    "PoolCounter",
		MType: "counter",
		Delta: &deltaCur,
	}
	valCur := float64(11.1)

	metricGauge := handlers.Metrics{
		ID:    "NumGC",
		MType: "gauge",
		Value: &valCur,
	}

	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	metricHandler := handlers.NewMetricHandler(&cfg)

	r := chi.NewRouter()

	r.Get(urlGetMetricsConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetNameMetricsHandler(w, r, stMetrics)
	})
	r.Get(urlGetPing, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetPing(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricHandler(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricJSONHandler(w, r, stMetrics)
	})
	r.Post(urlUpdatesMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdatesMetricJSONHandler(w, r, stMetrics)
	})
	r.Get(urlGetMetricConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricHandler(w, r, stMetrics)
	})
	r.Post(urlGetMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricJSONHandler(w, r, stMetrics)
	})

	bodyMetr, _ := json.Marshal(metricCounter)
	req := httptest.NewRequest(http.MethodPost, urlUpdateMetricJSONConst, bytes.NewReader(bodyMetr))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	bodyMetr, _ = json.Marshal(metricGauge)
	req = httptest.NewRequest(http.MethodPost, urlUpdateMetricJSONConst, bytes.NewReader(bodyMetr))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp = w.Result()
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	// Output:
	// 200
	// 200
}

func ExampleMetricsHandler_UpdatesMetricJSONHandler() {
	var metrics []handlers.Metrics
	deltaCur := int64(111)
	metric := handlers.Metrics{
		ID:    "PoolCount",
		MType: "counter",
		Delta: &deltaCur,
	}

	metrics = append(metrics, metric)
	valCur := float64(111.1)
	metric = handlers.Metrics{
		ID:    "PoolCount",
		MType: "gauge",
		Value: &valCur,
	}
	metrics = append(metrics, metric)

	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	metricHandler := handlers.NewMetricHandler(&cfg)

	r := chi.NewRouter()

	r.Get(urlGetMetricsConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetNameMetricsHandler(w, r, stMetrics)
	})
	r.Get(urlGetPing, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetPing(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricHandler(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricJSONHandler(w, r, stMetrics)
	})
	r.Post(urlUpdatesMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdatesMetricJSONHandler(w, r, stMetrics)
	})
	r.Get(urlGetMetricConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricHandler(w, r, stMetrics)
	})
	r.Post(urlGetMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricJSONHandler(w, r, stMetrics)
	})

	bodyMetr, _ := json.Marshal(metrics)
	req := httptest.NewRequest(http.MethodPost, urlUpdatesMetricJSONConst, bytes.NewReader(bodyMetr))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	// Output:
	// 200
}

func ExampleMetricsHandler_GetMetricJSONHandler() {
	stMetrics := &store.StorageContext{}
	mockStore := new(mocks.MockStore)
	mockStore.On("GetCounter", mock.Anything, "PoolCounter").Return(int64(1111), true)
	mockStore.On("GetCounter", mock.Anything, "PoolCounte").Return(int64(1111), false)
	mockStore.On("GetGauge", mock.Anything, "PoolGauge").Return(float64(1111.1), true)
	mockStore.On("GetGauge", mock.Anything, "PoolGaug").Return(float64(1111.1), false)
	stMetrics.SetStorage(mockStore)

	metricHandler := handlers.NewMetricHandler(&cfg)

	r := chi.NewRouter()

	r.Get(urlGetMetricsConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetNameMetricsHandler(w, r, stMetrics)
	})
	r.Get(urlGetPing, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetPing(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricHandler(w, r, stMetrics)
	})
	r.Post(urlUpdateMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricJSONHandler(w, r, stMetrics)
	})
	r.Post(urlUpdatesMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdatesMetricJSONHandler(w, r, stMetrics)
	})
	r.Get(urlGetMetricConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricHandler(w, r, stMetrics)
	})
	r.Post(urlGetMetricJSONConst, func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricJSONHandler(w, r, stMetrics)
	})

	dataMetric := handlers.Metrics{
		ID:    "PoolCounter",
		MType: "counter",
	}
	bodyMetr, _ := json.Marshal(dataMetric)
	req := httptest.NewRequest(http.MethodPost, urlGetMetricJSONConst, bytes.NewReader(bodyMetr))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	// Output:
	// 200
}
