package xrealip

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"gotest.tools/v3/assert"
)

func TestRealIP(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name          string
		trustedSubnet string
		xip           string
		want          want
	}{
		{
			name:          "positive test #1",
			trustedSubnet: "192.168.1.0/24",
			xip:           "192.168.1.1",
			want: want{
				code: 200,
			},
		},
		{
			name:          "positive test #2",
			trustedSubnet: "",
			xip:           "192.168.1.1",
			want: want{
				code: 200,
			},
		},
		{
			name:          "negative test #1",
			trustedSubnet: "192.168.1.0/24",
			xip:           "192.168.2.1",
			want: want{
				code: 403,
			},
		},
		{
			name:          "negative test #2",
			trustedSubnet: "192.168.1.0/24",
			want: want{
				code: 403,
			},
		},
		{
			name:          "negative test #3",
			trustedSubnet: "192.168.1.0/24",
			xip:           "192.168.2",
			want: want{
				code: 403,
			},
		},
		{
			name:          "negative test #4",
			trustedSubnet: "192.168..0/24",
			xip:           "192.168.2",
			want: want{
				code: 400,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Use(RealIP(test.trustedSubnet))
			r.Post("/update", func(w http.ResponseWriter, _ *http.Request) {
				w.Write([]byte("test"))
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodPost, "/update", bytes.NewReader([]byte("test")))
			req.Header.Add("X-Real-IP", test.xip)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.want.code, w.Code)
		})
	}
}
