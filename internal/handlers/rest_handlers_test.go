package handlers

import (
	"bytes"
	"encoding/json"
	"metrics/internal/server/configure"
	"metrics/internal/store"
	"metrics/internal/store/mocks"
	"metrics/internal/store/ramstorage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMetricsHandler_UpdateMetricJSONHandler(t *testing.T) {
	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	var cfg configure.Config
	cfg.ReadStartParams()

	metricHandler := NewMetricHandler(cfg)

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

	delt := int64(1111)
	valGauge := float64(1111.1)
	type want struct {
		code int
	}
	tests := []struct {
		name       string
		typeReqest string
		dataMetric Metrics
		want       want
	}{
		{
			name:       "rest positive test update metrics #1",
			typeReqest: http.MethodPost,
			dataMetric: Metrics{
				ID:    "PoolCounter",
				MType: "counter",
				Delta: &delt},
			want: want{
				code: 200,
			},
		},
		{
			name:       "rest positive test update metrics #2",
			typeReqest: http.MethodPost,
			dataMetric: Metrics{
				ID:    "PoolGauge",
				MType: "gauge",
				Value: &valGauge},
			want: want{
				code: 200,
			},
		},
		{
			name:       "rest positive test update metrics #3",
			typeReqest: http.MethodPost,
			dataMetric: Metrics{
				ID:    "PoolGauge",
				MType: "gaue",
				Value: &valGauge},
			want: want{
				code: 400,
			},
		},
		{
			name:       "rest positive test update metrics #3",
			typeReqest: http.MethodPost,
			dataMetric: Metrics{
				ID:    "PoolGauge",
				MType: "gaue",
				Value: &valGauge},
			want: want{
				code: 400,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bodyMetr, _ := json.Marshal(test.dataMetric)

			req := httptest.NewRequest(test.typeReqest, urlUpdateMetricJSONConst, bytes.NewReader(bodyMetr))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != test.want.code {
				t.Errorf("expected status OK; got %v", w.Code)
			}

			assert.Equal(t, test.want.code, w.Code)
		})
	}
}

func TestMetricsHandler_UpdatesMetricJSONHandler(t *testing.T) {
	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	var cfg configure.Config
	cfg.ReadStartParams()

	metricHandler := NewMetricHandler(cfg)

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

	delt := int64(1111)
	valGauge := float64(1111.1)
	var metrics []Metrics
	metrics = append(metrics, Metrics{
		ID:    "PoolCounter",
		MType: "counter",
		Delta: &delt})
	metrics = append(metrics, Metrics{
		ID:    "PoolGauge",
		MType: "gauge",
		Value: &valGauge})

	var metricsBad []Metrics
	metricsBad = append(metricsBad, Metrics{
		ID:    "PoolCounter",
		MType: "countr",
		Delta: &delt})
	metricsBad = append(metricsBad, Metrics{
		ID:    "PoolGauge",
		MType: "gauge",
		Value: &valGauge})
	type want struct {
		code int
	}
	tests := []struct {
		name       string
		typeReqest string
		dataMetric []Metrics
		want       want
	}{
		{
			name:       "rest positive test update metrics #1",
			typeReqest: http.MethodPost,
			dataMetric: metrics,
			want: want{
				code: 200,
			},
		},
		{
			name:       "rest positive test update metrics #2",
			typeReqest: http.MethodPost,
			want: want{
				code: 200,
			},
		},
		{
			name:       "rest positive test update metrics #3",
			typeReqest: http.MethodGet,
			want: want{
				code: 405,
			},
		},
		{
			name:       "rest positive test update metrics #4",
			typeReqest: http.MethodPost,
			dataMetric: metricsBad,
			want: want{
				code: 400,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bodyMetr, _ := json.Marshal(test.dataMetric)

			req := httptest.NewRequest(test.typeReqest, urlUpdatesMetricJSONConst, bytes.NewReader(bodyMetr))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != test.want.code {
				t.Errorf("expected status OK; got %v", w.Code)
			}

			assert.Equal(t, test.want.code, w.Code)
		})
	}
}

func TestMetricsHandler_GetMetricJSONHandler(t *testing.T) {
	stMetrics := &store.StorageContext{}
	mockStore := new(mocks.MockStore)
	mockStore.On("GetCounter", mock.Anything, "PoolCounter").Return(int64(1111), true)
	mockStore.On("GetCounter", mock.Anything, "PoolCounte").Return(int64(1111), false)
	mockStore.On("GetGauge", mock.Anything, "PoolGauge").Return(float64(1111.1), true)
	mockStore.On("GetGauge", mock.Anything, "PoolGaug").Return(float64(1111.1), false)
	stMetrics.SetStorage(mockStore)

	var cfg configure.Config
	cfg.ReadStartParams()

	metricHandler := NewMetricHandler(cfg)

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

	type want struct {
		code int
	}
	tests := []struct {
		name       string
		typeReqest string
		dataMetric Metrics
		want       want
	}{
		{
			name:       "rest positive test update metrics #1",
			typeReqest: http.MethodPost,
			dataMetric: Metrics{
				ID:    "PoolCounter",
				MType: "counter"},
			want: want{
				code: 200,
			},
		},
		{
			name:       "rest positive test update metrics #2",
			typeReqest: http.MethodPost,
			dataMetric: Metrics{
				ID:    "PoolCounte",
				MType: "counter"},
			want: want{
				code: 404,
			},
		},
		{
			name:       "rest positive test update metrics #3",
			typeReqest: http.MethodPost,
			dataMetric: Metrics{
				ID:    "PoolGauge",
				MType: "gauge"},
			want: want{
				code: 200,
			},
		},
		{
			name:       "rest positive test update metrics #3",
			typeReqest: http.MethodPost,
			dataMetric: Metrics{
				ID:    "PoolGaug",
				MType: "gauge"},
			want: want{
				code: 404,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bodyMetr, _ := json.Marshal(test.dataMetric)

			req := httptest.NewRequest(test.typeReqest, urlGetMetricJSONConst, bytes.NewReader(bodyMetr))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != test.want.code {
				t.Errorf("expected status OK; got %v", w.Code)
			}

			assert.Equal(t, test.want.code, w.Code)
		})
	}
}
