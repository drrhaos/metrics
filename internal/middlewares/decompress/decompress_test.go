// Модуль decompress предназначен распаковки упакованных в запросе данных.

package decompress

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
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

func TestGzipDecompressMiddleware(t *testing.T) {
	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	metricHandler := handlers.NewMetricHandler(&cfg)

	r := chi.NewRouter()

	r.Use(GzipDecompressMiddleware)
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

	var buf bytes.Buffer
	zipF := gzip.NewWriter(&buf)
	zipF.Write(bodyMetr)

	type want struct {
		code int
	}
	tests := []struct {
		name string
		data io.Reader
		want want
	}{
		{
			name: "rest positive test update metrics #1",
			data: &buf,
			want: want{
				code: 200,
			},
		},
		{
			name: "rest negative test update metrics #2",
			data: bytes.NewReader(bodyMetr),
			want: want{
				code: 500,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/update", test.data)
			req.Header.Set("Content-Encoding", "gzip")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			fmt.Println()
			assert.Equal(t, test.want.code, w.Code)
		})
	}
}
