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
		wantLenTables int64
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
			fmt.Println(tables)
		})
	}
}
