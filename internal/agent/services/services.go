package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"metrics/internal/agent/configure"
	"metrics/internal/agent/gzip"
	"metrics/internal/logger"
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

type Metrics struct {
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
}

func customDelay() retry.DelayTypeFunc {
	return func(n uint, err error, config *retry.Config) time.Duration {
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

func updateMertics(ctx context.Context, metricsCPU *store.StorageContext, PollCount int64) {
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

func updateMerticsGops(ctx context.Context, metricsCPU *store.StorageContext) {
	m, _ := mem.VirtualMemory()
	totalMem := m.Total
	metricsCPU.UpdateGauge(ctx, gaugesTotalMem, float64(totalMem))
	freeMem := m.Free
	metricsCPU.UpdateGauge(ctx, gaugesFreeMem, float64(freeMem))
	countCPU, _ := cpu.Counts(false)
	metricsCPU.UpdateGauge(ctx, gaugesCPUutil, float64(countCPU))
}

func sendAllMetric(ctx context.Context, metrics []Metrics, cfg configure.Config) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*30))
	defer cancel()

	client := &http.Client{}

	urlStr := fmt.Sprintf(urlUpdateMetricsJSONConst, cfg.Address)
	reqData, err := json.Marshal(metrics)
	if err != nil {
		logger.Log.Warn("Не удалось создать JSON", zap.Error(err))
		return
	}

	buf, err := gzip.CompressReqData(reqData)
	if err != nil {
		logger.Log.Warn("Не удалось сжать данные", zap.Error(err))
		return
	}
	err = retry.Do(
		func() error {
			r, _ := http.NewRequest(http.MethodPost, urlStr, buf)
			r = r.WithContext(ctx)
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Content-Encoding", "gzip")
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
	}
}

func sendMetricsWorker(ctx context.Context, workerID int, jobs <-chan []Metrics, cfg configure.Config) {
	for job := range jobs {
		logger.Log.Info(fmt.Sprintf("Воркер %d количество метрик %d", workerID, len(job)))
		sendAllMetric(ctx, job, cfg)
	}
}

func prepareBatch(ctx context.Context, metricsCPU *store.StorageContext, cfg configure.Config) (metricsBatches [][]Metrics) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var metrics []Metrics

	currentGauges, ok := metricsCPU.GetGauges(ctx)
	if !ok {
		return metricsBatches
	}

	for nameMetric, valueMetric := range currentGauges {
		hd := valueMetric
		metrics = append(metrics, Metrics{ID: nameMetric, MType: typeMetricGauge, Value: &hd})
	}

	currentCounters, ok := metricsCPU.GetCounters(ctx)
	if !ok {
		return metricsBatches
	}
	for nameMetric, valueMetric := range currentCounters {
		hd := valueMetric
		metrics = append(metrics, Metrics{ID: nameMetric, MType: typeMetricCounter, Delta: &hd})
	}

	lenMetrics := len(metrics)
	var countBatch int
	if cfg.RateLimit > lenMetrics {
		countBatch = lenMetrics
	} else if cfg.RateLimit <= lenMetrics {
		countBatch = cfg.RateLimit
	}
	metricsBatches = make([][]Metrics, countBatch)

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

func CollectMetrics(cfg configure.Config) {
	jobs := make(chan []Metrics, cfg.RateLimit)

	metricsCPU := &store.StorageContext{}
	metricsCPU.SetStorage(ramstorage.NewStorage())

	ctx := context.Background()

	var mut sync.Mutex
	go func() {
		var PollCount int64
		for {
			PollCount++
			mut.Lock()
			updateMertics(ctx, metricsCPU, PollCount)
			mut.Unlock()
			time.Sleep(time.Duration(cfg.PollInterval) * time.Second)
		}
	}()

	go func() {
		for {
			mut.Lock()
			updateMerticsGops(ctx, metricsCPU)
			mut.Unlock()
			time.Sleep(time.Duration(cfg.PollInterval) * time.Second)
		}
	}()

	for w := 1; w <= cfg.RateLimit; w++ {
		go func(workerID int) {
			sendMetricsWorker(ctx, workerID, jobs, cfg)
		}(w)
	}

	for {
		for _, metrics := range prepareBatch(ctx, metricsCPU, cfg) {
			jobs <- metrics
		}
		time.Sleep(time.Duration(cfg.ReportInterval) * time.Second)
	}
}
