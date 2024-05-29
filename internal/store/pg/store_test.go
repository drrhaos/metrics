//go:build integration
// +build integration

package pg

import (
	"context"
	"fmt"
	"testing"

	"github.com/adrianbrad/psqldocker"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestDatabase_SaveMetrics(t *testing.T) {
	type fields struct {
		Conn *pgxpool.Pool
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
			args: args{
				filePath: "",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.fields.Conn,
			}
			if got := db.SaveMetrics(tt.args.filePath); got != tt.want {
				t.Errorf("Database.SaveMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_LoadMetrics(t *testing.T) {
	type fields struct {
		Conn *pgxpool.Pool
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
			args: args{
				filePath: "",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.fields.Conn,
			}
			if got := db.LoadMetrics(tt.args.filePath); got != tt.want {
				t.Errorf("Database.LoadMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_Migrations(t *testing.T) {
	const (
		usr           = "usr"
		password      = "pass"
		dbName        = "tst"
		containerName = "psql_docker_tests"
	)

	c, _ := psqldocker.NewContainer(
		usr,
		password,
		dbName,
		psqldocker.WithContainerName(containerName),
	)
	defer func() {
		c.Close()
	}()

	dsn := fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"host=localhost "+
			"port=%s "+
			"sslmode=disable",
		usr,
		password,
		dbName,
		c.Port(),
	)

	db := NewDatabase(dsn)
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		wantLenTables int
	}{
		{
			name: "positive test",
			args: args{
				ctx: context.Background(),
			},
			wantErr:       false,
			wantLenTables: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.Migrations(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Database.Migrations() error = %v, wantErr %v", err, tt.wantErr)
			}

			rows, _ := db.Conn.Query(context.Background(), "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")

			defer rows.Close()

			var tables []string
			for rows.Next() {
				var tableName string
				rows.Scan(&tableName)

				tables = append(tables, tableName)
			}

			assert.Equal(t, tt.wantLenTables, len(tables))
		})
	}
}

func TestDatabase_Ping(t *testing.T) {
	const (
		usr           = "usr"
		password      = "pass"
		dbName        = "tst"
		containerName = "psql_docker_tests"
	)

	c, _ := psqldocker.NewContainer(
		usr,
		password,
		dbName,
		psqldocker.WithContainerName(containerName),
	)
	defer func() {
		c.Close()
	}()

	dsn := fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"host=localhost "+
			"port=%s "+
			"sslmode=disable",
		usr,
		password,
		dbName,
		c.Port(),
	)

	db := NewDatabase(dsn)

	type fields struct {
		Conn *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
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
			if got := db.Ping(tt.args.ctx); got != tt.want {
				t.Errorf("Database.Ping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_UpdateGauge(t *testing.T) {
	const (
		usr           = "usr"
		password      = "pass"
		dbName        = "tst"
		containerName = "psql_docker_tests"
	)

	c, _ := psqldocker.NewContainer(
		usr,
		password,
		dbName,
		psqldocker.WithContainerName(containerName),
	)
	defer func() {
		c.Close()
	}()

	dsn := fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"host=localhost "+
			"port=%s "+
			"sslmode=disable",
		usr,
		password,
		dbName,
		c.Port(),
	)

	db := NewDatabase(dsn)

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
				nameMetric:  "test",
				valueMetric: 111.1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := db.UpdateGauge(tt.args.ctx, tt.args.nameMetric, tt.args.valueMetric); got != tt.want {
				t.Errorf("Database.UpdateGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_UpdateCounter(t *testing.T) {
	const (
		usr           = "usr"
		password      = "pass"
		dbName        = "tst"
		containerName = "psql_docker_tests"
	)

	c, _ := psqldocker.NewContainer(
		usr,
		password,
		dbName,
		psqldocker.WithContainerName(containerName),
	)
	defer func() {
		c.Close()
	}()

	dsn := fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"host=localhost "+
			"port=%s "+
			"sslmode=disable",
		usr,
		password,
		dbName,
		c.Port(),
	)

	db := NewDatabase(dsn)

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
				nameMetric:  "test",
				valueMetric: 111,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := db.UpdateCounter(tt.args.ctx, tt.args.nameMetric, tt.args.valueMetric); got != tt.want {
				t.Errorf("Database.UpdateCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_GetGauge(t *testing.T) {
	const (
		usr           = "usr"
		password      = "pass"
		dbName        = "tst"
		containerName = "psql_docker_tests"
	)

	c, _ := psqldocker.NewContainer(
		usr,
		password,
		dbName,
		psqldocker.WithContainerName(containerName),
	)
	defer func() {
		c.Close()
	}()

	dsn := fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"host=localhost "+
			"port=%s "+
			"sslmode=disable",
		usr,
		password,
		dbName,
		c.Port(),
	)

	db := NewDatabase(dsn)

	db.UpdateGauge(context.Background(), "test", 111.1)

	type args struct {
		ctx        context.Context
		nameMetric string
	}
	tests := []struct {
		name            string
		args            args
		wantValueMetric float64
		wantExists      bool
	}{
		{
			name: "positive test",
			args: args{
				ctx:        context.Background(),
				nameMetric: "test",
			},
			wantValueMetric: 111.1,
			wantExists:      true,
		},
		{
			name: "negative test",
			args: args{
				ctx:        context.Background(),
				nameMetric: "test2",
			},
			wantValueMetric: 0,
			wantExists:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValueMetric, gotExists := db.GetGauge(tt.args.ctx, tt.args.nameMetric)
			if gotValueMetric != tt.wantValueMetric {
				t.Errorf("Database.GetGauge() gotValueMetric = %v, want %v", gotValueMetric, tt.wantValueMetric)
			}
			if gotExists != tt.wantExists {
				t.Errorf("Database.GetGauge() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}

func TestDatabase_GetGauges(t *testing.T) {
	const (
		usr           = "usr"
		password      = "pass"
		dbName        = "tst"
		containerName = "psql_docker_tests"
	)

	c, _ := psqldocker.NewContainer(
		usr,
		password,
		dbName,
		psqldocker.WithContainerName(containerName),
	)
	defer func() {
		c.Close()
	}()

	dsn := fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"host=localhost "+
			"port=%s "+
			"sslmode=disable",
		usr,
		password,
		dbName,
		c.Port(),
	)

	db := NewDatabase(dsn)

	db.UpdateGauge(context.Background(), "test", 111.1)
	db.UpdateGauge(context.Background(), "test2", 22.1)

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
				"test":  111.1,
				"test2": 22.1,
			},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1 := db.GetGauges(tt.args.ctx)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Database.GetGauges() got = %v, want %v", got, tt.want)
			// }
			if got1 != tt.want1 {
				t.Errorf("Database.GetGauges() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDatabase_GetCounters(t *testing.T) {
	const (
		usr           = "usr"
		password      = "pass"
		dbName        = "tst"
		containerName = "psql_docker_tests"
	)

	c, _ := psqldocker.NewContainer(
		usr,
		password,
		dbName,
		psqldocker.WithContainerName(containerName),
	)
	defer func() {
		c.Close()
	}()

	dsn := fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"host=localhost "+
			"port=%s "+
			"sslmode=disable",
		usr,
		password,
		dbName,
		c.Port(),
	)

	db := NewDatabase(dsn)

	db.UpdateCounter(context.Background(), "test", 111)
	db.UpdateCounter(context.Background(), "test2", 22)

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
			want:  2,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := db.GetCounters(tt.args.ctx)

			assert.Equal(t, tt.want, len(got))
			if got1 != tt.want1 {
				t.Errorf("Database.GetCounters() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDatabase_GetBatchMetrics(t *testing.T) {
	const (
		usr           = "usr"
		password      = "pass"
		dbName        = "tst"
		containerName = "psql_docker_tests"
	)

	c, _ := psqldocker.NewContainer(
		usr,
		password,
		dbName,
		psqldocker.WithContainerName(containerName),
	)
	defer func() {
		c.Close()
	}()

	dsn := fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"host=localhost "+
			"port=%s "+
			"sslmode=disable",
		usr,
		password,
		dbName,
		c.Port(),
	)

	db := NewDatabase(dsn)

	db.UpdateCounter(context.Background(), "testCounter", 111)
	db.UpdateGauge(context.Background(), "testGauge", 111.1)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name          string
		args          args
		wantLenTables int
		wantExists    bool
	}{
		{
			name: "positive test",
			args: args{
				ctx: context.Background(),
			},
			wantExists:    true,
			wantLenTables: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMetrics, gotExists := db.GetBatchMetrics(tt.args.ctx)

			assert.Equal(t, tt.wantLenTables, len(gotMetrics))
			if gotExists != tt.wantExists {
				t.Errorf("Database.GetBatchMetrics() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}

func TestDatabase_GetCounter(t *testing.T) {
	const (
		usr           = "usr"
		password      = "pass"
		dbName        = "tst"
		containerName = "psql_docker_tests"
	)

	c, _ := psqldocker.NewContainer(
		usr,
		password,
		dbName,
		psqldocker.WithContainerName(containerName),
	)
	defer func() {
		c.Close()
	}()

	dsn := fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"host=localhost "+
			"port=%s "+
			"sslmode=disable",
		usr,
		password,
		dbName,
		c.Port(),
	)

	db := NewDatabase(dsn)

	db.UpdateCounter(context.Background(), "test", 1111)

	type args struct {
		ctx        context.Context
		nameMetric string
	}
	tests := []struct {
		name            string
		args            args
		wantValueMetric int64
		wantExists      bool
	}{
		{
			name: "positive test",
			args: args{
				ctx:        context.Background(),
				nameMetric: "test",
			},
			wantValueMetric: 1111,
			wantExists:      true,
		},
		{
			name: "negative test",
			args: args{
				ctx:        context.Background(),
				nameMetric: "test2",
			},
			wantValueMetric: 0,
			wantExists:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValueMetric, gotExists := db.GetCounter(tt.args.ctx, tt.args.nameMetric)
			if gotValueMetric != tt.wantValueMetric {
				t.Errorf("Database.GetCounter() gotValueMetric = %v, want %v", gotValueMetric, tt.wantValueMetric)
			}
			if gotExists != tt.wantExists {
				t.Errorf("Database.GetCounter() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}
