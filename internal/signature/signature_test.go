// Модуль signature предназначен для проверки целостности запроса.

package signature

import (
	"bytes"
	"encoding/json"
	"fmt"
	"metrics/internal/handlers"
	"metrics/internal/server/configure"
	"metrics/internal/store"
	"metrics/internal/store/ramstorage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

var cfg = configure.Config{}

func TestCheckSignaturMiddleware(t *testing.T) {
	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	metricHandler := handlers.NewMetricHandler(&cfg)

	r := chi.NewRouter()
	key := "test"
	r.Use(CheckSignaturMiddleware(key))
	r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetNameMetricsHandler(w, r, stMetrics)
	})

	valGauge := float64(1111.1)
	dataMetric := handlers.Metrics{
		ID:    "PoolGauge",
		MType: "gaue",
		Value: &valGauge,
	}
	bodyMetr, _ := json.Marshal(dataMetric)
	type want struct {
		code int
	}
	tests := []struct {
		name string
		hash string
		data []byte
		want want
	}{
		{
			name: "rest positive test update metrics #1",
			hash: "MWWZx-8HPTidG99XrdIPVICaB1okBYCEaUZ2lxrcVdI=",
			data: bodyMetr,
			want: want{
				code: 200,
			},
		},
		{
			name: "rest negative test update metrics #2",
			hash: "ZchLggTYnSeSsLHDE_3aI19eZEf2KL-06jh1mbG9hIY=",
			data: bodyMetr,
			want: want{
				code: 400,
			},
		},
		{
			name: "rest negative test update metrics #3",
			hash: "ZchLggTYnSeSsLHDE_3aI19eZEf2KL-06jh1mbGhIY=",
			data: bodyMetr,
			want: want{
				code: 400,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/update", bytes.NewReader(bodyMetr))

			req.Header.Set("HashSHA256", test.hash)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.want.code, w.Code)
		})
	}

}

func TestAddSignatureMiddleware(t *testing.T) {
	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	metricHandler := handlers.NewMetricHandler(&cfg)

	r := chi.NewRouter()
	key := "test"
	r.Use(AddSignatureMiddleware(key))
	r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
		metricHandler.GetNameMetricsHandler(w, r, stMetrics)
	})

	valGauge := float64(1111.1)
	dataMetric := handlers.Metrics{
		ID:    "PoolGauge",
		MType: "gaue",
		Value: &valGauge,
	}
	bodyMetr, _ := json.Marshal(dataMetric)
	type want struct {
		hash string
	}
	tests := []struct {
		name string
		data []byte
		want want
	}{
		{
			name: "rest positive test update metrics #1",

			data: bodyMetr,
			want: want{
				hash: "9Win5C3nSnfe85m8mati4_VeNk4D1OCjwUxea7nLG6Q=",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/update", bytes.NewReader(bodyMetr))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			fmt.Println()
			assert.Equal(t, test.want.hash, w.Header().Get("HashSHA256"))
		})
	}
}
