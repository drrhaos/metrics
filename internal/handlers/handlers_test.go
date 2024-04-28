package handlers_test

import (
	"fmt"
	_ "metrics/internal/handlers"
	"metrics/internal/server/configure"
	"net/http"
)

func ExampleMetricsHandler_UpdateMetricHandler() {
	var cfg configure.Config
	cfg.ReadStartParams()

	typeMetric := "counter"
	nameMetric := "ddd"
	valueMetric := 11

	client := &http.Client{}
	urlStr := fmt.Sprintf("http://%s/update/%s/%s/%d", cfg.Address, typeMetric, nameMetric, valueMetric)

	r, err := http.NewRequest(http.MethodPost, urlStr, nil)
	if err != nil {
		return
	}
	r.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	typeMetric = "gauge"
	nameMetric = "ddd"
	gaugeMetric := 11.1

	urlStr = fmt.Sprintf("http://%s/update/%s/%s/%f", cfg.Address, typeMetric, nameMetric, gaugeMetric)

	r, err = http.NewRequest(http.MethodPost, urlStr, nil)
	if err != nil {
		return
	}
	resp, err = client.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	// Output:
	// 200
	// 200
}

// func ExampleMetricsHandler_GetMetricHandler() {

// 	typeMetric := "counter"
// 	nameMetric := "ddd"

// 	client := &http.Client{}
// 	urlStr := fmt.Sprintf("http://127.0.0.1:8080/value/%s/%s", typeMetric, nameMetric)

// 	r, _ := http.NewRequest(http.MethodGet, urlStr, nil)
// 	r.Header.Set("Content-Type", "application/json")
// 	resp, err := client.Do(r)
// 	if err != nil {
// 		logger.Log.Warn("Не удалось отправить запрос", zap.Error(err))
// 	}
// 	defer resp.Body.Close()
// 	fmt.Println(resp.StatusCode)

// 	typeMetric = "gauge"
// 	nameMetric = "ddd"

// 	urlStr = fmt.Sprintf("http://127.0.0.1:8080/value/%s/%s", typeMetric, nameMetric)

// 	r, _ = http.NewRequest(http.MethodGet, urlStr, nil)
// 	resp, err = client.Do(r)
// 	if err != nil {
// 		logger.Log.Warn("Не удалось отправить запрос", zap.Error(err))
// 	}
// 	defer resp.Body.Close()
// 	fmt.Println(resp.StatusCode)

// 	// Output:
// 	// 200
// 	// 200
// }

// func ExampleMetricsHandler_GetNameMetricsHandler() {

// 	client := &http.Client{}
// 	urlStr := "http://127.0.0.1:8080/"

// 	r, _ := http.NewRequest(http.MethodGet, urlStr, nil)
// 	r.Header.Set("Content-Type", "application/json")
// 	resp, err := client.Do(r)
// 	if err != nil {
// 		logger.Log.Warn("Не удалось отправить запрос", zap.Error(err))
// 	}
// 	defer resp.Body.Close()
// 	fmt.Println(resp.StatusCode)

// 	// Output:
// 	// 200
// }

// func ExampleMetricsHandler_GetPing() {

// 	client := &http.Client{}
// 	urlStr := "http://127.0.0.1:8080/ping"

// 	r, _ := http.NewRequest(http.MethodGet, urlStr, nil)
// 	r.Header.Set("Content-Type", "application/json")
// 	resp, err := client.Do(r)
// 	if err != nil {
// 		logger.Log.Warn("Не удалось отправить запрос", zap.Error(err))
// 	}
// 	defer resp.Body.Close()
// 	fmt.Println(resp.StatusCode)

// 	// Output:
// 	// 500
// }

// func ExampleMetricsHandler_UpdateMetricJSONHandler() {

// 	deltaCur := int64(111)
// 	metricCounter := Metrics{
// 		ID:    "PoolCount",
// 		MType: "counter",
// 		Delta: &deltaCur,
// 	}
// 	valCur := float64(111)

// 	metricGauge := Metrics{
// 		ID:    "PoolCount",
// 		MType: "gauge",
// 		Value: &valCur,
// 	}

// 	client := &http.Client{}
// 	urlStr := "http://127.0.0.1:8080/update/"

// 	var buf bytes.Buffer
// 	json.NewEncoder(&buf).Encode(metricCounter)

// 	r, _ := http.NewRequest(http.MethodPost, urlStr, &buf)
// 	r.Header.Set("Content-Type", "application/json")
// 	resp, err := client.Do(r)
// 	if err != nil {
// 		logger.Log.Warn("Не удалось отправить запрос", zap.Error(err))
// 	}
// 	defer resp.Body.Close()
// 	fmt.Println(resp.StatusCode)

// 	json.NewEncoder(&buf).Encode(metricGauge)

// 	r, _ = http.NewRequest(http.MethodPost, urlStr, &buf)
// 	r.Header.Set("Content-Type", "application/json")
// 	resp, err = client.Do(r)
// 	if err != nil {
// 		logger.Log.Warn("Не удалось отправить запрос", zap.Error(err))
// 	}
// 	defer resp.Body.Close()
// 	fmt.Println(resp.StatusCode)

// 	// Output:
// 	// 200
// 	// 200
// }

// func ExampleMetricsHandler_UpdatesMetricJSONHandler() {

// 	var metrics []Metrics
// 	deltaCur := int64(111)
// 	metric := Metrics{
// 		ID:    "PoolCount",
// 		MType: "counter",
// 		Delta: &deltaCur,
// 	}

// 	metrics = append(metrics, metric)
// 	valCur := float64(111.1)
// 	metric = Metrics{
// 		ID:    "PoolCount",
// 		MType: "gauge",
// 		Value: &valCur,
// 	}
// 	metrics = append(metrics, metric)

// 	client := &http.Client{}
// 	urlStr := "http://127.0.0.1:8080/updates/"

// 	var buf bytes.Buffer
// 	json.NewEncoder(&buf).Encode(metrics)

// 	r, _ := http.NewRequest(http.MethodPost, urlStr, &buf)
// 	r.Header.Set("Content-Type", "application/json")
// 	resp, err := client.Do(r)
// 	if err != nil {
// 		logger.Log.Warn("Не удалось отправить запрос", zap.Error(err))
// 	}
// 	defer resp.Body.Close()
// 	fmt.Println(resp.StatusCode)

// 	// Output:
// 	// 200
// }

// func ExampleMetricsHandler_GetMetricJSONHandler() {

// 	metric := Metrics{
// 		ID:    "PoolCount",
// 		MType: "counter",
// 	}
// 	var buf bytes.Buffer
// 	json.NewEncoder(&buf).Encode(metric)

// 	client := &http.Client{}
// 	urlStr := "http://127.0.0.1:8080//"

// 	r, _ := http.NewRequest(http.MethodPost, urlStr, &buf)
// 	r.Header.Set("Content-Type", "application/json")
// 	resp, err := client.Do(r)
// 	if err != nil {
// 		logger.Log.Warn("Не удалось отправить запрос", zap.Error(err))
// 	}
// 	defer resp.Body.Close()
// 	fmt.Println(resp.StatusCode)

// 	// Output:
// 	// 404
// }
