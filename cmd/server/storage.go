package main

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func (storage *MemStorage) init() {
	storage.counter = make(map[string]int64)
	storage.gauge = make(map[string]float64)
}

func (storage *MemStorage) updateGauge(nameMetric string, valueMetric float64) {
	storage.gauge[nameMetric] = valueMetric
}

func (storage *MemStorage) updateCounter(nameMetric string, valueMetric int64) {
	storage.counter[nameMetric] += valueMetric
}

func (storage *MemStorage) getGauge(nameMetric string) (float64, bool) {
	cur, ok := storage.gauge[nameMetric]
	return cur, ok
}

func (storage *MemStorage) getCounter(nameMetric string) (int64, bool) {
	cur, ok := storage.counter[nameMetric]
	return cur, ok
}
