// Package restmode пакет запуска http сервера
package restmode

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"metrics/internal/server/configure"
	"metrics/internal/store"
	"metrics/internal/store/ramstorage"

	"gotest.tools/v3/assert"
)

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func TestRun(t *testing.T) {
	var cfg configure.Config
	cfg.ReadConfig()

	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	tests := []struct {
		name       string
		typeReqest string
		urlStr     string
		want       int
	}{
		{
			name:       "positive test #1",
			typeReqest: http.MethodGet,
			urlStr:     "/",
			want:       http.StatusOK,
		},
		{
			name:       "positive test #2",
			typeReqest: http.MethodGet,
			urlStr:     "/ping",
			want:       http.StatusInternalServerError,
		},
		{
			name:       "positive test #3",
			typeReqest: http.MethodGet,
			urlStr:     "/value/test/test",
			want:       http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port, _ := getFreePort()
			cfg.Address = fmt.Sprintf("127.0.0.1:%d", port)

			go Run(cfg, stMetrics)

			time.Sleep(5 * time.Second)

			client := &http.Client{}
			r, _ := http.NewRequest(tt.typeReqest, fmt.Sprintf("http://%s%s", cfg.Address, tt.urlStr), nil)
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			assert.Equal(t, resp.StatusCode, tt.want)
		})
	}
}
