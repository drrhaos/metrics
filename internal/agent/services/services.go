// Package services пакет сбора и отправки метрик
package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"syscall"
	"time"

	pb "metrics/internal/proto"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"metrics/internal/agent/configure"
	"metrics/internal/agent/gzip"
	"metrics/internal/logger"
	"metrics/internal/middlewares/cryptodata"
	"metrics/internal/store"
	"metrics/internal/store/ramstorage"

	"github.com/avast/retry-go"
	"go.uber.org/zap"
)

const (
	typeMetricCounter = "counter"
	typeMetricGauge   = "gauge"
	randomValueName   = "RandomValue"
	pollCountName     = "PollCount"
	gaugesTotalMem    = "TotalMemory"
	gaugesFreeMem     = "FreeMemory"
	gaugesCPUutil     = "CPUutilization1"
)

const urlUpdateMetricsJSONConst = "http://%s/updates/"

var sleepStep = map[uint]int64{0: 1, 1: 3, 2: 5}

var nameGauges = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

func customDelay() retry.DelayTypeFunc {
	return func(n uint, _ error, _ *retry.Config) time.Duration {
		return time.Duration(sleepStep[n])
	}
}

func getFloat64MemStats(m runtime.MemStats, name string) (float64, bool) {
	value := reflect.ValueOf(m).FieldByName(name)
	var floatValue float64
	switch value.Kind() {
	case reflect.Uint64:
		floatValue = float64(value.Uint())
	case reflect.Uint32:
		floatValue = float64(value.Uint())
	case reflect.Float64:
		floatValue = value.Float()
	default:
		logger.Log.Info("Тип значения не соответствует uint")
		return floatValue, false
	}
	return floatValue, true
}

func updateMertics(ctx context.Context, doneCh <-chan os.Signal, metricsCPU *store.StorageContext, cfg *configure.Config) {
	var PollCount int64
	timer := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Завершено обновление метрик")
			return
		case <-doneCh:
			logger.Log.Info("Завершено обновление метрик")
			return
		default:
			<-timer.C

			PollCount++
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			metricsCPU.UpdateGauge(ctx, randomValueName, rand.Float64())
			metricsCPU.UpdateCounter(ctx, pollCountName, PollCount)

			for _, name := range nameGauges {
				floatValue, ok := getFloat64MemStats(m, name)
				if ok {
					metricsCPU.UpdateGauge(ctx, name, floatValue)
				}
			}
		}
	}
}

func updateMerticsGops(ctx context.Context, doneCh <-chan os.Signal, metricsCPU *store.StorageContext, cfg *configure.Config) {
	timer := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Завершено обновление метрик")
			return
		case <-doneCh:
			logger.Log.Info("Завершено обновление метрик GOPS")
			return
		default:
			<-timer.C

			m, _ := mem.VirtualMemory()
			totalMem := m.Total
			metricsCPU.UpdateGauge(ctx, gaugesTotalMem, float64(totalMem))
			freeMem := m.Free
			metricsCPU.UpdateGauge(ctx, gaugesFreeMem, float64(freeMem))
			countCPU, _ := cpu.Counts(false)
			metricsCPU.UpdateGauge(ctx, gaugesCPUutil, float64(countCPU))
		}
	}
}

func getRealIP() (string, error) {
	var ipAddresses string

	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if (ip.To4() != nil || ip.To16() != nil) && ipAddresses == "" && !ip.IsLoopback() {
				ipAddresses = ip.String()
			} else if (ip.To4() != nil || ip.To16() != nil) && ipAddresses != "" && !ip.IsLoopback() {
				ipAddresses = fmt.Sprintf("%s, %s", ipAddresses, ip.String())
			}
		}
	}

	return ipAddresses, nil
}

func sendRESTMetric(ctx context.Context, metrics []store.Metrics, cfg configure.Config) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	client := &http.Client{}

	urlStr := fmt.Sprintf(urlUpdateMetricsJSONConst, cfg.Address)
	reqData, err := json.Marshal(metrics)
	if err != nil {
		logger.Log.Warn("Не удалось создать JSON", zap.Error(err))
		return err
	}
	if cfg.CryptoKeyPath != "" {
		cryptoKeyByte, errCrypt := os.ReadFile(cfg.CryptoKeyPath)
		if errCrypt != nil {
			logger.Log.Warn("Не удалось прочитать файл ключа", zap.Error(err))
			return errCrypt
		}

		pemBlock, _ := pem.Decode(cryptoKeyByte)
		cryptoKey, errCrypt := x509.ParsePKIXPublicKey(pemBlock.Bytes)
		if errCrypt != nil {
			logger.Log.Warn("Не удалось распарсить файл ключа", zap.Error(err))
			return errCrypt
		}

		reqData, err = cryptodata.Encrypt(reqData, cryptoKey)
		if err != nil {
			logger.Log.Warn("Не удалось зашифровать данные", zap.Error(err))
			return err
		}

	}

	buf, err := gzip.CompressReqData(reqData)
	if err != nil {
		logger.Log.Warn("Не удалось сжать данные", zap.Error(err))
		return err
	}

	err = retry.Do(
		func() error {
			r, _ := http.NewRequest(http.MethodPost, urlStr, buf)
			r = r.WithContext(ctx)
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Content-Encoding", "gzip")
			realIP, errReal := getRealIP()
			if realIP != "" && errReal == nil {
				r.Header.Set("X-Real-IP", realIP)
			}
			if cfg.Key != "" {
				h := hmac.New(sha256.New, []byte(cfg.Key))
				h.Write(reqData)
				hashReq := h.Sum(nil)
				r.Header.Set("HashSHA256", base64.URLEncoding.EncodeToString(hashReq))
			}
			resp, errCLient := client.Do(r)
			if errCLient != nil {
				logger.Log.Warn("Не удалось отправить запрос", zap.Error(errCLient))
				return errCLient
			}
			defer resp.Body.Close()
			return nil
		},
		retry.Attempts(3),
		retry.DelayType(customDelay()),
	)
	if err != nil {
		logger.Log.Warn("Не удалось отправить данные", zap.Error(err))
		return err
	}
	return nil
}


func sendGRPCMetric(ctx context.Context, metrics []store.Metrics, cfg configure.Config) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	conn, err := grpc.NewClient(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Warn("Не удалось установить соединение с сервером", zap.Error(err))
	}
	defer conn.Close()

	c := pb.NewMetricsClient(conn)

	var pbMetrics []*pb.Metric
	for _, metric := range metrics {
		switch metric.MType {
		case "counter":
			pbMetrics = append(pbMetrics, &pb.Metric{
				Id: metric.ID,
				Type: metric.MType,
				Delta: *metric.Delta,
				})
		case "gauge":			
			pbMetrics = append(pbMetrics, &pb.Metric{
				Id: metric.ID,
				Type: metric.MType,
				Value: *metric.Value,
				})
		}

	}

	_ , err = c.UpdateMetrics(ctx, &pb.AddMetricsRequest{
		Metric: pbMetrics,
	})

	return err
}

func sendMetricsWorker(ctx context.Context, workerID int, jobs <-chan []store.Metrics, cfg configure.Config) {
	for job := range jobs {
		logger.Log.Info(fmt.Sprintf("Воркер %d количество метрик %d", workerID, len(job)))
		sendRESTMetric(ctx, job, cfg)
		sendGRPCMetric(ctx, job, cfg)
	}
}

func prepareBatch(ctx context.Context, metricsCPU *store.StorageContext, cfg configure.Config) (metricsBatches [][]store.Metrics) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	metrics, _ := metricsCPU.GetBatchMetrics(ctx)

	lenMetrics := len(metrics)
	var countBatch int
	if cfg.RateLimit > lenMetrics {
		countBatch = lenMetrics
	} else if cfg.RateLimit <= lenMetrics {
		countBatch = cfg.RateLimit
	}
	metricsBatches = make([][]store.Metrics, countBatch)

	i := 0
	for j := 0; j < lenMetrics; j++ {
		if i >= cfg.RateLimit {
			i = 0
		}
		metricsBatches[i] = append(metricsBatches[i], metrics[j])
		i++
	}

	return metricsBatches
}

// CollectMetrics осуществляет сбор и отправку метрик на сервер
func CollectMetrics(ctx context.Context, cfg configure.Config) {
	doneChUpdate := make(chan os.Signal, 1)
	signal.Notify(doneChUpdate, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	doneChUpdateGops := make(chan os.Signal, 1)
	signal.Notify(doneChUpdateGops, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	doneChSend := make(chan os.Signal, 1)
	signal.Notify(doneChSend, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	jobs := make(chan []store.Metrics, cfg.RateLimit)

	metricsCPU := &store.StorageContext{}
	metricsCPU.SetStorage(ramstorage.NewStorage())

	go func() {
		updateMertics(ctx, doneChUpdate, metricsCPU, &cfg)
	}()

	go func() {
		updateMerticsGops(ctx, doneChUpdateGops, metricsCPU, &cfg)
	}()

	for w := 1; w <= cfg.RateLimit; w++ {
		go func(workerID int) {
			sendMetricsWorker(ctx, workerID, jobs, cfg)
		}(w)
	}
	var doneSend bool
	for !doneSend {
		select {
		case <-doneChSend:
			logger.Log.Info("Завершена отправка метрик")
			doneSend = true
		case <-ctx.Done():
			logger.Log.Info("Завершена отправка метрик")
			return
		default:
			for _, metrics := range prepareBatch(ctx, metricsCPU, cfg) {
				jobs <- metrics
			}
			time.Sleep(time.Duration(cfg.ReportInterval) * time.Second)
		}
	}
	logger.Log.Info("Агент остановлен")
}
