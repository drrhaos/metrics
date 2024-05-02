package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"metrics/internal/server/configure"
	"metrics/internal/store"
	"metrics/internal/store/ramstorage"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
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

func Benchmark_TestMetricsHandler_UpdateMetricHandler(b *testing.B) {
	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	metricHandler := NewMetricHandler(&cfg)

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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/update/gauge/NumGC/11", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

	}
}

func Benchmark_TestMetricsHandler_GetMetricHandler(b *testing.B) {
	ctx := context.Background()

	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	stMetrics.UpdateCounter(ctx, "testCounter", 10)
	stMetrics.UpdateGauge(ctx, "testGauge", 11.1)
	stMetrics.UpdateGauge(ctx, "testGauge2", 12.1)

	metricHandler := NewMetricHandler(&cfg)

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetNameMetricsHandler(w, r, stMetrics)
	})
	r.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricHandler(w, r, stMetrics)
	})
	r.Get("/value/{typeMetric}/{nameMetric}", func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricHandler(w, r, stMetrics)
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/value/gauge/testGauge", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

	}
}

func Benchmark_TestMetricsHandler_GetNameMetricsHandler(b *testing.B) {
	ctx := context.Background()

	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	stMetrics.UpdateCounter(ctx, "testCounter", 10)
	stMetrics.UpdateGauge(ctx, "testGauge", 11.1)
	stMetrics.UpdateGauge(ctx, "testGauge2", 12.1)

	metricHandler := NewMetricHandler(&cfg)

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetNameMetricsHandler(w, r, stMetrics)
	})
	r.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricHandler(w, r, stMetrics)
	})
	r.Get("/value/{typeMetric}/{nameMetric}", func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricHandler(w, r, stMetrics)
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}

func TestMetricsHandler_UpdateMetricHandler(t *testing.T) {
	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	metricHandler := NewMetricHandler(&cfg)

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
		url        string
		typeReqest string
		want       want
	}{
		{
			name:       "positive test update counter#1",
			url:        "/update/counter/PollCount/10",
			typeReqest: http.MethodPost,
			want: want{
				code: 200,
			},
		},
		{
			name:       "positive test update gauge #2",
			url:        "/update/gauge/NumGC/0",
			typeReqest: http.MethodPost,
			want: want{
				code: 200,
			},
		},
		{
			name:       "negative test typeMetric not gauge and not counter #3",
			url:        "/update/ddd/NumGC/0",
			typeReqest: http.MethodPost,
			want: want{
				code: 400,
			},
		},
		{
			name:       "negative test valueMetric counter not int  #4",
			url:        "/update/counter/NumGC/dddd",
			typeReqest: http.MethodPost,
			want: want{
				code: 400,
			},
		},
		{
			name:       "negative test valueMetric gauge not float  #5",
			url:        "/update/gauge/NumGC/dddd",
			typeReqest: http.MethodPost,
			want: want{
				code: 400,
			},
		},
		{
			name:       "negative test method get #6",
			url:        "/update/gauge/NumGC/11",
			typeReqest: http.MethodGet,
			want: want{
				code: 405,
			},
		},
		{
			name:       "negative test not found nameMetric #6",
			url:        "/update/gauge//11",
			typeReqest: http.MethodPost,
			want: want{
				code: 404,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.typeReqest, test.url, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != test.want.code {
				t.Errorf("expected status OK; got %v", w.Code)
			}

			assert.Equal(t, test.want.code, w.Code)
		})
	}
}

func TestMetricsHandler_GetMetricHandler(t *testing.T) {
	ctx := context.Background()

	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	stMetrics.UpdateCounter(ctx, "testCounter", 10)
	stMetrics.UpdateGauge(ctx, "testGauge", 11.1)
	stMetrics.UpdateGauge(ctx, "testGauge2", 12.1)

	metricHandler := NewMetricHandler(&cfg)

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetNameMetricsHandler(w, r, stMetrics)
	})
	r.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricHandler(w, r, stMetrics)
	})
	r.Get("/value/{typeMetric}/{nameMetric}", func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricHandler(w, r, stMetrics)
	})

	type want struct {
		code  int
		value string
	}
	tests := []struct {
		name       string
		url        string
		typeReqest string
		want       want
	}{
		{
			name:       "positive test get counter #1",
			url:        "/value/counter/testCounter",
			typeReqest: http.MethodGet,
			want: want{
				code:  200,
				value: "10",
			},
		},
		{
			name:       "positive test get gauge #2",
			url:        "/value/gauge/testGauge",
			typeReqest: http.MethodGet,
			want: want{
				code:  200,
				value: "11.1",
			},
		},
		{
			name:       "negative test empty type #3",
			url:        "/value//testGauge",
			typeReqest: http.MethodGet,
			want: want{
				code:  404,
				value: "",
			},
		},
		{
			name:       "negative test empty name #4",
			url:        "/value/gauge/",
			typeReqest: http.MethodGet,
			want: want{
				code:  404,
				value: "404 page not found\n",
			},
		},
		{
			name:       "negative test bad gauge name #5",
			url:        "/value/gauge/testGauge4",
			typeReqest: http.MethodGet,
			want: want{
				code: 404,
			},
		},
		{
			name:       "negative test bad counter name #6",
			url:        "/value/counter/testGauge",
			typeReqest: http.MethodGet,
			want: want{
				code: 404,
			},
		},
		{
			name:       "negative test not found nameMetric #7",
			url:        "/update/gauge//11",
			typeReqest: http.MethodPost,
			want: want{
				code: 404,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.typeReqest, test.url, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != test.want.code {
				t.Errorf("expected status OK; got %v", w.Code)
			}

			assert.Equal(t, test.want.code, w.Code)
			assert.Equal(t, test.want.value, w.Body.String())
		})
	}
}

func TestMetricsHandler_GetNameMetricsHandler(t *testing.T) {
	ctx := context.Background()

	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	stMetrics.UpdateCounter(ctx, "testCounter", 10)
	stMetrics.UpdateGauge(ctx, "testGauge", 11.1)
	stMetrics.UpdateGauge(ctx, "testGauge2", 12.1)

	metricHandler := NewMetricHandler(&cfg)

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetNameMetricsHandler(w, r, stMetrics)
	})
	r.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", func(w http.ResponseWriter, r *http.Request) {
		metricHandler.UpdateMetricHandler(w, r, stMetrics)
	})
	r.Get("/value/{typeMetric}/{nameMetric}", func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetMetricHandler(w, r, stMetrics)
	})

	type want struct {
		code   int
		lenRes int
	}
	tests := []struct {
		name       string
		url        string
		typeReqest string
		want       want
	}{
		{
			name:       "positive test get metrics #1",
			url:        "/",
			typeReqest: http.MethodGet,
			want: want{
				code:   200,
				lenRes: 178,
			},
		},
		{
			name:       "negative test post metrics #2",
			url:        "/",
			typeReqest: http.MethodPost,
			want: want{
				code:   405,
				lenRes: 0,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.typeReqest, test.url, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != test.want.code {
				t.Errorf("expected status OK; got %v", w.Code)
			}
			assert.Equal(t, test.want.code, w.Code)
			assert.Equal(t, test.want.lenRes, len(w.Body.String()))
		})
	}
}
