package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"metrics/internal/logger"
	"net/http"

	"go.uber.org/zap"
)

func ExampleMetricsHandler_UpdateMetricJSONHandler() {

	deltaCur := int64(111)
	metricCounter := Metrics{
		ID:    "PoolCount",
		MType: "counter",
		Delta: &deltaCur,
	}
	valCur := float64(111)

	metricGauge := Metrics{
		ID:    "PoolCount",
		MType: "gauge",
		Value: &valCur,
	}

	client := &http.Client{}
	urlStr := "http://127.0.0.1:8080/update/"

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(metricCounter)

	r, _ := http.NewRequest(http.MethodPost, urlStr, &buf)
	r.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(r)
	if err != nil {
		logger.Log.Warn("Не удалось отправить запрос", zap.Error(err))
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	json.NewEncoder(&buf).Encode(metricGauge)

	r, _ = http.NewRequest(http.MethodPost, urlStr, &buf)
	r.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(r)
	if err != nil {
		logger.Log.Warn("Не удалось отправить запрос", zap.Error(err))
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	// Output:
	// 200
	// 200
}

func ExampleMetricsHandler_UpdatesMetricJSONHandler() {

	var metrics []Metrics
	deltaCur := int64(111)
	metric := Metrics{
		ID:    "PoolCount",
		MType: "counter",
		Delta: &deltaCur,
	}

	metrics = append(metrics, metric)
	valCur := float64(111.1)
	metric = Metrics{
		ID:    "PoolCount",
		MType: "gauge",
		Value: &valCur,
	}
	metrics = append(metrics, metric)

	client := &http.Client{}
	urlStr := "http://127.0.0.1:8080/updates/"

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(metrics)

	r, _ := http.NewRequest(http.MethodPost, urlStr, &buf)
	r.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(r)
	if err != nil {
		logger.Log.Warn("Не удалось отправить запрос", zap.Error(err))
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	// Output:
	// 200
}

func ExampleMetricsHandler_GetMetricJSONHandler() {

	metric := Metrics{
		ID:    "PoolCount",
		MType: "counter",
	}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(metric)

	client := &http.Client{}
	urlStr := "http://127.0.0.1:8080//"

	r, _ := http.NewRequest(http.MethodPost, urlStr, &buf)
	r.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(r)
	if err != nil {
		logger.Log.Warn("Не удалось отправить запрос", zap.Error(err))
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	// Output:
	// 404
}
