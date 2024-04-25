package handlers

import (
	"fmt"
	"metrics/internal/logger"
	"net/http"

	"go.uber.org/zap"
)

func ExampleMetricsHandler_UpdateMetricHandler() {

	typeMetric := "counter"
	nameMetric := "ddd"
	valueMetric := 11

	client := &http.Client{}
	urlStr := fmt.Sprintf("http://127.0.0.1:8080/update/%s/%s/%d", typeMetric, nameMetric, valueMetric)

	r, _ := http.NewRequest(http.MethodPost, urlStr, nil)
	r.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(r)
	if err != nil {
		logger.Log.Warn("Не удалось отправить запрос", zap.Error(err))
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	typeMetric = "gauge"
	nameMetric = "ddd"
	gaugeMetric := 11.1

	urlStr = fmt.Sprintf("http://127.0.0.1:8080/update/%s/%s/%f", typeMetric, nameMetric, gaugeMetric)

	r, _ = http.NewRequest(http.MethodPost, urlStr, nil)
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

func ExampleMetricsHandler_GetMetricHandler() {

	typeMetric := "counter"
	nameMetric := "ddd"

	client := &http.Client{}
	urlStr := fmt.Sprintf("http://127.0.0.1:8080/value/%s/%s", typeMetric, nameMetric)

	r, _ := http.NewRequest(http.MethodGet, urlStr, nil)
	r.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(r)
	if err != nil {
		logger.Log.Warn("Не удалось отправить запрос", zap.Error(err))
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	typeMetric = "gauge"
	nameMetric = "ddd"

	urlStr = fmt.Sprintf("http://127.0.0.1:8080/value/%s/%s", typeMetric, nameMetric)

	r, _ = http.NewRequest(http.MethodGet, urlStr, nil)
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

func ExampleMetricsHandler_GetNameMetricsHandler() {

	client := &http.Client{}
	urlStr := "http://127.0.0.1:8080/"

	r, _ := http.NewRequest(http.MethodGet, urlStr, nil)
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

func ExampleMetricsHandler_GetPing() {

	client := &http.Client{}
	urlStr := "http://127.0.0.1:8080/ping"

	r, _ := http.NewRequest(http.MethodGet, urlStr, nil)
	r.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(r)
	if err != nil {
		logger.Log.Warn("Не удалось отправить запрос", zap.Error(err))
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	// Output:
	// 500
}
