// Package grpcmode запускает сервер.
package grpcmode

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "metrics/internal/proto"
	"metrics/internal/server/configure"
	"metrics/internal/store"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Run запускает сервер
func Run(_ context.Context, cfg configure.Config, stMetric *store.StorageContext) {
	listen, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()
	// регистрируем сервис
	metricsServer := MetricsServer{
		storage: stMetric,
	}
	pb.RegisterMetricsServer(s, &metricsServer)

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

// MetricsServer поддерживает все необходимые методы сервера.
type MetricsServer struct {
	pb.UnimplementedMetricsServer

	storage *store.StorageContext
}

// UpdateMetrics обновляет значение пачки метрики
func (ms *MetricsServer) UpdateMetrics(ctx context.Context, in *pb.AddMetricsRequest) (*pb.AddMetricsResponse, error) {
	var response pb.AddMetricsResponse

	for _, metr := range in.Metric {
		switch metr.Type {
		case "gauge":
			ms.storage.UpdateGauge(ctx, metr.Id, metr.Value)
		case "counter":
			ms.storage.UpdateCounter(ctx, metr.Id, metr.Delta)
		default:
			return nil, status.Errorf(codes.NotFound, `Тип метрики %s не найден`, metr.Type)
		}
	}

	return &response, nil
}
