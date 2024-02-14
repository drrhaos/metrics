package main

type StorageInterface interface {
	SaveMetrics(filePath string) bool
	LoadMetrics(filePath string) bool
	UpdateCounter(nameMetric string, valueMetric int64) bool
	UpdateGauge(nameMetric string, valueMetric float64) bool
	GetGauges() (map[string]float64, bool)
	GetCounters() (map[string]int64, bool)
	GetCounter(nameMetric string) (currentValue int64, exists bool)
	GetGauge(nameMetric string) (currentValue float64, exists bool)
	Ping() bool
}

type StorageContext struct {
	storage StorageInterface
}

func (sc *StorageContext) setStorage(storage StorageInterface) {
	sc.storage = storage
}

func (sc *StorageContext) SaveMetrics(filePath string) bool {
	return sc.storage.SaveMetrics(filePath)
}

func (sc *StorageContext) LoadMetrics(filePath string) bool {
	return sc.storage.LoadMetrics(filePath)
}

func (sc *StorageContext) UpdateCounter(nameMetric string, valueMetric int64) bool {
	return sc.storage.UpdateCounter(nameMetric, valueMetric)
}

func (sc *StorageContext) UpdateGauge(nameMetric string, valueMetric float64) bool {
	return sc.storage.UpdateGauge(nameMetric, valueMetric)
}

func (sc *StorageContext) GetGauges() (map[string]float64, bool) {
	return sc.storage.GetGauges()
}

func (sc *StorageContext) GetCounters() (map[string]int64, bool) {
	return sc.storage.GetCounters()
}

func (sc *StorageContext) GetCounter(nameMetric string) (currentValue int64, exists bool) {
	return sc.storage.GetCounter(nameMetric)
}

func (sc *StorageContext) GetGauge(nameMetric string) (currentValue float64, exists bool) {
	return sc.storage.GetGauge(nameMetric)
}

func (sc *StorageContext) Ping() (exists bool) {
	return sc.storage.Ping()
}
