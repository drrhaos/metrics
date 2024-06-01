// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"
	store "metrics/internal/store"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestMockStore_GetBatchMetrics(t *testing.T) {
	var mo []store.Metrics
	delt := int64(1111)
	mo = append(mo, store.Metrics{
		ID:    "PoolCount",
		Delta: &delt,
		MType: "counter",
	})
	mockStore := new(MockStore)
	mockStore.On("GetBatchMetrics", mock.Anything).Return(mo, true)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{
			name: "positive test",
			args: args{
				ctx: context.Background(),
			},
			want:  1,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, got1 := mockStore.GetBatchMetrics(tt.args.ctx)
			assert.Equal(t, tt.want, len(got))

			if got1 != tt.want1 {
				t.Errorf("MockStore.GetBatchMetrics() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMockStore_GetCounter(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("GetCounter", mock.Anything, "PoolCounter").Return(int64(10), true)
	mockStore.On("GetCounter", mock.Anything, "PoolCounters").Return(int64(0), false)

	type args struct {
		ctx        context.Context
		nameMetric string
	}
	tests := []struct {
		name  string
		args  args
		want  int64
		want1 bool
	}{
		{
			name: "positive test",
			args: args{
				ctx:        context.Background(),
				nameMetric: "PoolCounter",
			},
			want:  10,
			want1: true,
		},
		{
			name: "negative test",
			args: args{
				ctx:        context.Background(),
				nameMetric: "PoolCounters",
			},
			want:  0,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, got1 := mockStore.GetCounter(tt.args.ctx, tt.args.nameMetric)
			if got != tt.want {
				t.Errorf("MockStore.GetCounter() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MockStore.GetCounter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMockStore_GetCounters(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("GetCounters", mock.Anything).Return(map[string]int64{
		"Pool":  100,
		"Pools": 1001,
	}, true)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name  string
		args  args
		want  map[string]int64
		want1 bool
	}{
		{
			name: "positive test",
			args: args{
				ctx: context.Background(),
			},
			want: map[string]int64{
				"Pool":  100,
				"Pools": 1001,
			},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := mockStore.GetCounters(tt.args.ctx)
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("MockStore.GetCounters() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MockStore.GetCounters() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMockStore_GetGauge(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("GetGauge", mock.Anything, "Alloc").Return(float64(10.1), true)
	mockStore.On("GetGauge", mock.Anything, "Mem").Return(float64(0), false)

	type args struct {
		ctx        context.Context
		nameMetric string
	}
	tests := []struct {
		name  string
		args  args
		want  float64
		want1 bool
	}{
		{
			name: "positive test",
			args: args{
				ctx:        context.Background(),
				nameMetric: "Alloc",
			},
			want:  10.1,
			want1: true,
		},
		{
			name: "negative test",
			args: args{
				ctx:        context.Background(),
				nameMetric: "Mem",
			},
			want:  0,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, got1 := mockStore.GetGauge(tt.args.ctx, tt.args.nameMetric)
			if got != tt.want {
				t.Errorf("MockStore.GetGauge() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MockStore.GetGauge() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMockStore_GetGauges(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("GetGauges", mock.Anything).Return(map[string]float64{
		"Alloc": 100.1,
		"Mem":   1001.1,
	}, true)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name  string
		args  args
		want  map[string]float64
		want1 bool
	}{
		{
			name: "positive test",
			args: args{
				ctx: context.Background(),
			},
			want: map[string]float64{
				"Alloc": 100.1,
				"Mem":   1001.1,
			},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := mockStore.GetGauges(tt.args.ctx)
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("MockStore.GetGauges() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MockStore.GetGauges() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMockStore_LoadMetrics(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("LoadMetrics", "/tmp/metrics.json").Return(true)
	mockStore.On("LoadMetrics", "/tmp/metracs.json").Return(false)

	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive test",
			args: args{
				filePath: "/tmp/metrics.json",
			},
			want: true,
		},
		{
			name: "negative test",
			args: args{
				filePath: "/tmp/metracs.json",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mockStore.LoadMetrics(tt.args.filePath); got != tt.want {
				t.Errorf("MockStore.LoadMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockStore_Ping(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("Ping", mock.Anything).Return(true)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive test",
			args: args{
				ctx: context.Background(),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mockStore.Ping(tt.args.ctx); got != tt.want {
				t.Errorf("MockStore.Ping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockStore_SaveMetrics(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("SaveMetrics", "/tmp/metrics.json").Return(true)
	mockStore.On("SaveMetrics", "/tmp/metracs.json").Return(false)

	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive test",
			args: args{
				filePath: "/tmp/metrics.json",
			},
			want: true,
		},
		{
			name: "negative test",
			args: args{
				filePath: "/tmp/metracs.json",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mockStore.SaveMetrics(tt.args.filePath); got != tt.want {
				t.Errorf("MockStore.SaveMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockStore_UpdateCounter(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("UpdateCounter", mock.Anything, "PoolCounter", int64(100)).Return(true)

	type args struct {
		ctx         context.Context
		nameMetric  string
		valueMetric int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive test",
			args: args{
				ctx:         context.Background(),
				nameMetric:  "PoolCounter",
				valueMetric: int64(100),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mockStore.UpdateCounter(tt.args.ctx, tt.args.nameMetric, tt.args.valueMetric); got != tt.want {
				t.Errorf("MockStore.UpdateCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockStore_UpdateGauge(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("UpdateGauge", mock.Anything, "Alloc", float64(100.1)).Return(true)

	type args struct {
		ctx         context.Context
		nameMetric  string
		valueMetric float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive test",
			args: args{
				ctx:         context.Background(),
				nameMetric:  "Alloc",
				valueMetric: float64(100.1),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mockStore.UpdateGauge(tt.args.ctx, tt.args.nameMetric, tt.args.valueMetric); got != tt.want {
				t.Errorf("MockStore.UpdateCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}