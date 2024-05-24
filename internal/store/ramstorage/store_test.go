package ramstorage

import (
	"context"
	"os"
	"reflect"
	"testing"
)

func TestRAMStorage_UpdateCounter(t *testing.T) {
	type args struct {
		ctx         context.Context
		nameMetric  string
		valueMetric int64
	}
	tests := []struct {
		name    string
		storage *RAMStorage
		args    args
		want    bool
	}{
		{
			name:    "Positive test update counter",
			storage: NewStorage(),
			args: args{
				ctx:         context.Background(),
				nameMetric:  "PoolCounter",
				valueMetric: 11,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.storage.UpdateCounter(tt.args.ctx, tt.args.nameMetric, tt.args.valueMetric); got != tt.want {
				t.Errorf("RAMStorage.UpdateCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRAMStorage_UpdateGauge(t *testing.T) {
	type args struct {
		ctx         context.Context
		nameMetric  string
		valueMetric float64
	}
	tests := []struct {
		name    string
		storage *RAMStorage
		args    args
		want    bool
	}{
		{
			name:    "Positive test update gauge",
			storage: NewStorage(),
			args: args{
				ctx:         context.Background(),
				nameMetric:  "Test1",
				valueMetric: 11,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.storage.UpdateGauge(tt.args.ctx, tt.args.nameMetric, tt.args.valueMetric); got != tt.want {
				t.Errorf("RAMStorage.UpdateGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRAMStorage_GetGauges(t *testing.T) {
	st := NewStorage()
	st.UpdateGauge(context.Background(), "Test1", 11.1)
	st.UpdateGauge(context.Background(), "Test2", 22.1)

	data := make(map[string]float64)
	data["Test1"] = 11.1
	data["Test2"] = 22.1

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		storage *RAMStorage
		want    map[string]float64
		args    args
		name    string
		want1   bool
	}{
		{
			name:    "Positive test get gauges",
			storage: st,
			args: args{
				ctx: context.Background(),
			},
			want:  data,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.storage.GetGauges(tt.args.ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RAMStorage.GetGauges() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("RAMStorage.GetGauges() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRAMStorage_GetCounters(t *testing.T) {
	st := NewStorage()
	st.UpdateCounter(context.Background(), "Test1", 11)
	st.UpdateCounter(context.Background(), "Test2", 22)

	data := make(map[string]int64)
	data["Test1"] = 11
	data["Test2"] = 22

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		storage *RAMStorage
		want    map[string]int64
		args    args
		name    string
		want1   bool
	}{
		{
			name:    "Positive test get counters",
			storage: st,
			args: args{
				ctx: context.Background(),
			},
			want:  data,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.storage.GetCounters(tt.args.ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RAMStorage.GetCounters() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("RAMStorage.GetCounters() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRAMStorage_GetCounter(t *testing.T) {
	st := NewStorage()
	st.UpdateCounter(context.Background(), "Test1", 11)
	st.UpdateCounter(context.Background(), "Test2", 22)

	type args struct {
		ctx        context.Context
		nameMetric string
	}
	tests := []struct {
		name             string
		storage          *RAMStorage
		args             args
		wantCurrentValue int64
		wantExists       bool
	}{
		{
			name:    "positive test get counter #1",
			storage: st,
			args: args{
				ctx:        context.Background(),
				nameMetric: "Test1",
			},
			wantCurrentValue: 11,
			wantExists:       true,
		},
		{
			name:    "positive test get counter #2",
			storage: st,
			args: args{
				ctx:        context.Background(),
				nameMetric: "Test2",
			},
			wantCurrentValue: 22,
			wantExists:       true,
		},
		{
			name:    "negative test get counter #3",
			storage: st,
			args: args{
				ctx:        context.Background(),
				nameMetric: "Test",
			},
			wantExists: false,
		},
		{
			name:    "negative test get counter #4",
			storage: nil,
			args: args{
				ctx:        context.Background(),
				nameMetric: "Test",
			},
			wantExists: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCurrentValue, gotExists := tt.storage.GetCounter(tt.args.ctx, tt.args.nameMetric)
			if gotCurrentValue != tt.wantCurrentValue {
				t.Errorf("RAMStorage.GetCounter() gotCurrentValue = %v, want %v", gotCurrentValue, tt.wantCurrentValue)
			}
			if gotExists != tt.wantExists {
				t.Errorf("RAMStorage.GetCounter() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}

func TestRAMStorage_GetGauge(t *testing.T) {
	st := NewStorage()
	st.UpdateGauge(context.Background(), "Test1", 11.1)
	st.UpdateGauge(context.Background(), "Test2", 22.1)

	type args struct {
		ctx        context.Context
		nameMetric string
	}
	tests := []struct {
		name             string
		storage          *RAMStorage
		args             args
		wantCurrentValue float64
		wantExists       bool
	}{
		{
			name:    "positive test get gauge #1",
			storage: st,
			args: args{
				ctx:        context.Background(),
				nameMetric: "Test1",
			},
			wantCurrentValue: 11.1,
			wantExists:       true,
		},
		{
			name:    "positive test get gauge #2",
			storage: st,
			args: args{
				ctx:        context.Background(),
				nameMetric: "Test2",
			},
			wantCurrentValue: 22.1,
			wantExists:       true,
		},
		{
			name:    "negative test get gauge #3",
			storage: st,
			args: args{
				ctx:        context.Background(),
				nameMetric: "Test",
			},
			wantExists: false,
		},
		{
			name:    "negative test get gauge #4",
			storage: nil,
			args: args{
				ctx:        context.Background(),
				nameMetric: "Test",
			},
			wantExists: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCurrentValue, gotExists := tt.storage.GetGauge(tt.args.ctx, tt.args.nameMetric)
			if gotCurrentValue != tt.wantCurrentValue {
				t.Errorf("RAMStorage.GetGauge() gotCurrentValue = %v, want %v", gotCurrentValue, tt.wantCurrentValue)
			}
			if gotExists != tt.wantExists {
				t.Errorf("RAMStorage.GetGauge() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}

func TestRAMStorage_Ping(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		storage *RAMStorage
		args    args
		name    string
		want    bool
	}{
		{
			name:    "Positive test ping",
			storage: NewStorage(),
			args: args{
				ctx: context.Background(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.storage.Ping(tt.args.ctx); got != tt.want {
				t.Errorf("RAMStorage.Ping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRAMStorage_SaveMetrics(t *testing.T) {
	fileTmp := "/tmp/tmp.json"
	fileTmpBad := "/teeeeemp/tmp.json"
	defer os.Remove(fileTmp)
	defer os.Remove(fileTmpBad)

	type fields struct {
		Gauge   map[string]float64
		Counter map[string]int64
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "positive test",
			fields: fields{
				Gauge:   map[string]float64{"first": 34.1},
				Counter: map[string]int64{"first": 34},
			},
			args: args{
				filePath: fileTmp,
			},
			want: true,
		},
		{
			name: "negative test",
			fields: fields{
				Gauge:   map[string]float64{"first": 34.1},
				Counter: map[string]int64{"first": 34},
			},
			args: args{
				filePath: fileTmpBad,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &RAMStorage{
				Gauge:   tt.fields.Gauge,
				Counter: tt.fields.Counter,
			}
			if got := storage.SaveMetrics(tt.args.filePath); got != tt.want {
				t.Errorf("RAMStorage.SaveMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRAMStorage_LoadMetrics(t *testing.T) {
	fileTmp := "/tmp/tmp.json"
	fileTmpBad := "/teeeeemp/tmp.json"
	defer os.Remove(fileTmp)
	defer os.Remove(fileTmpBad)

	storage := &RAMStorage{
		Gauge:   map[string]float64{"first": 34.1},
		Counter: map[string]int64{"first": 34},
	}
	storage.SaveMetrics(fileTmp)

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
				filePath: fileTmp,
			},
			want: true,
		},
		{
			name: "negative test",
			args: args{
				filePath: fileTmpBad,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := storage.LoadMetrics(tt.args.filePath); got != tt.want {
				t.Errorf("RAMStorage.LoadMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
