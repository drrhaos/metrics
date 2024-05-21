// Модуль logger предназначен для логирования в агенте и сервере.

package logger

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestRequestLogger(t *testing.T) {
	r := chi.NewRouter()
	r.Use(RequestLogger)
	r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("test")
		w.Write([]byte("test"))
		w.WriteHeader(http.StatusOK)
	})

	type want struct {
		code int
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive positive check signature #1",
			want: want{
				code: 200,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Initialize("info")

			req := httptest.NewRequest(http.MethodPost, "/update", nil)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.want.code, w.Code)
		})
	}
}
