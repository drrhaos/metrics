// Package grpcmode запускает сервер.
package grpcmode

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	pb "metrics/internal/proto"
	"metrics/internal/server/configure"
	"metrics/internal/store"
	"metrics/internal/store/ramstorage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		name        string
		wantCounter int64
		wantGauge   float64
	}{
		{
			name:        "positive test #1",
			wantCounter: 1111,
			wantGauge:   1111.1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port, _ := getFreePort()
			cfg.Address = fmt.Sprintf("127.0.0.1:%d", port)

			go Run(cfg, stMetrics)

			time.Sleep(5 * time.Second)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()

			conn, err := grpc.NewClient(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			defer conn.Close()

			c := pb.NewMetricsClient(conn)

			var pbMetrics []*pb.Metric
			pbMetrics = append(pbMetrics, &pb.Metric{
				Id:    "test",
				Type:  "counter",
				Delta: 1111,
			})
			pbMetrics = append(pbMetrics, &pb.Metric{
				Id:    "test2",
				Type:  "gauge",
				Value: 1111.1,
			})

			_, err = c.UpdateMetrics(ctx, &pb.AddMetricsRequest{
				Metric: pbMetrics,
			})

			curCounter, _ := stMetrics.GetCounter(ctx, "test")
			curGauge, _ := stMetrics.GetGauge(ctx, "test2")

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			assert.Equal(t, curCounter, tt.wantCounter)
			assert.Equal(t, curGauge, tt.wantGauge)
		})
	}
}
