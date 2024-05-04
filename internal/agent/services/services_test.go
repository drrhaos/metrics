package services

import (
	"context"
	"runtime"
	"testing"

	"metrics/internal/agent/configure"
	"metrics/internal/store"
	"metrics/internal/store/ramstorage"

	"github.com/stretchr/testify/assert"
)

var cfg = configure.Config{}

func Test_prepareBatch(t *testing.T) {
	cfg.ReadConfig()
	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	stMetrics.UpdateCounter(context.Background(), "PoolCount", 111)
	stMetrics.UpdateGauge(context.Background(), "TestGug", 111.1)

	stMetrics2 := &store.StorageContext{}
	stMetrics2.SetStorage(ramstorage.NewStorage())

	stMetrics2.UpdateCounter(context.Background(), "PoolCount", 111)

	stMetrics3 := &store.StorageContext{}
	stMetrics3.SetStorage(ramstorage.NewStorage())

	tests := []struct {
		metr *store.StorageContext
		name string
	}{
		{
			name: "positive test #1",
			metr: stMetrics,
		},
		{
			name: "positive test #2",
			metr: stMetrics2,
		},
		{
			name: "positive test #3",
			metr: stMetrics3,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, metrics := range prepareBatch(context.Background(), test.metr, cfg) {
				for _, cur := range metrics {
					if cur.MType == "gauge" {
						curVal, _ := stMetrics.GetGauge(context.Background(), cur.ID)
						assert.Equal(t, *cur.Value, curVal)
					} else if cur.MType == "counter" {
						curVal, _ := stMetrics.GetCounter(context.Background(), cur.ID)
						assert.Equal(t, *cur.Delta, curVal)
					}
				}
			}
		})
	}
}

func Test_updateMertics(t *testing.T) {
	stMetrics := &store.StorageContext{}
	stMetrics.SetStorage(ramstorage.NewStorage())

	type args struct {
		ctx        context.Context
		metricsCPU *store.StorageContext
		PollCount  int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "positive test #1",
			args: args{
				ctx:        context.Background(),
				metricsCPU: stMetrics,
				PollCount:  1,
			},
			want: 1,
		},
		{
			name: "positive test #2",
			args: args{
				ctx:        context.Background(),
				metricsCPU: stMetrics,
				PollCount:  1,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateMertics(tt.args.ctx, tt.args.metricsCPU, tt.args.PollCount)
			poolC, _ := stMetrics.GetCounter(context.Background(), "PollCount")
			assert.Equal(t, tt.want, poolC)
		})
	}
}

func Test_getFloat64MemStats(t *testing.T) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	m.Alloc = 1000
	m.GCCPUFraction = 100.1
	type args struct {
		name string
		m    runtime.MemStats
	}
	tests := []struct {
		name  string
		args  args
		want  float64
		want1 bool
	}{
		{
			name: "positive test #1",
			args: args{
				m:    m,
				name: "Alloc",
			},
			want:  1000,
			want1: true,
		},
		{
			name: "positive test #2",
			args: args{
				m:    m,
				name: "GCCPUFraction",
			},
			want:  100.1,
			want1: true,
		},
		{
			name: "negative test #2",
			args: args{
				m:    m,
				name: "DebugGC",
			},
			want:  0,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getFloat64MemStats(tt.args.m, tt.args.name)
			if got != tt.want {
				t.Errorf("getFloat64MemStats() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getFloat64MemStats() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
