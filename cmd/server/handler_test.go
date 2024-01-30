package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func Test_updateMetricHandler(t *testing.T) {
	storage := &MemStorage{
		Counter: make(map[string]int64),
		Gauge:   make(map[string]float64),
	}

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

func Test_getMetricHandler(t *testing.T) {
	storage := &MemStorage{
		Counter: make(map[string]int64),
		Gauge:   make(map[string]float64),
	}
	storage.updateCounter("testCounter", 10)
	storage.updateGauge("testGauge", 11.1)
	storage.updateGauge("testGauge2", 12.1)

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

func Test_getNameMetricsHandler(t *testing.T) {
	storage := &MemStorage{
		Counter: make(map[string]int64),
		Gauge:   make(map[string]float64),
	}
	storage.updateCounter("testCounter", 10)
	storage.updateGauge("testGauge", 11.1)
	storage.updateGauge("testGauge2", 12.1)

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
