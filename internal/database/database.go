package database

import (
	"context"

	"github.com/drrhaos/metrics/internal/logger"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const pathMigrationsConst = "file://%s/migrations/"

type Database struct {
	Conn *pgx.Conn
	uri  string
}

func NewDatabase() *Database {
	return &Database{}
}

func (db *Database) Close() {
	db.Conn.Close(context.Background())
}

func (db *Database) IsClosed() bool {
	if db.Conn == nil || db.Conn.IsClosed() {
		return true
	}
	return false
}

func (db *Database) Ping() bool {
	if err := db.Conn.Ping(context.Background()); err != nil {
		return false
	}
	return true
}

func (db *Database) Connect(uri string) error {
	conn, err := pgx.Connect(context.Background(), uri)
	if err != nil {
		logger.Log.Info("Не удалось установить соединение с базой данных", zap.Error(err))
		return err
	}
	db.Conn = conn
	db.uri = uri
	logger.Log.Info("Соединение с базой успешно установлено")
	return nil
}

func (db *Database) Migrations() error {
	var exist bool
	err := db.Conn.QueryRow(context.Background(),
		`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'counters')`).Scan(&exist)
	if err != nil {
		logger.Log.Warn("Не удалось выполнить запрос", zap.Error(err))
		return err
	}
	if !exist {
		_, err := db.Conn.Exec(context.Background(),
			`CREATE TABLE counters (
				"id" SERIAL PRIMARY KEY,
				"name" character(40) NOT NULL UNIQUE,
				"value" integer NOT NULL
			)`)
		if err != nil {
			logger.Log.Warn("Не удалось создать таблицу", zap.Error(err))
			return err
		}
	}

	err = db.Conn.QueryRow(context.Background(),
		`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'gauges')`).Scan(&exist)
	if err != nil {
		logger.Log.Warn("Не удалось выполнить запрос", zap.Error(err))
		return err
	}
	if !exist {
		_, err := db.Conn.Exec(context.Background(),
			`CREATE TABLE gauges (
				"id" SERIAL PRIMARY KEY,
				"name" character(40) NOT NULL UNIQUE,
				"value" double precision NOT NULL
			)`)
		if err != nil {
			logger.Log.Warn("Не удалось создать таблицу", zap.Error(err))
			return err
		}
	}

	return nil
}

func (db *Database) UpdateGauge(nameMetric string, valueMetric float64) bool {
	_, err := db.Conn.Exec(context.Background(),
		`INSERT INTO gauges (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE
		SET value = EXCLUDED.value`, nameMetric, valueMetric)
	if err != nil {
		logger.Log.Warn("Не удалось обновить значение", zap.Error(err))
		return false
	}
	return true
}

func (db *Database) UpdateCounter(nameMetric string, valueMetric int64) bool {
	_, err := db.Conn.Exec(context.Background(),
		`INSERT INTO counters (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE
		SET value = EXCLUDED.value`, nameMetric, valueMetric)
	if err != nil {
		logger.Log.Warn("Не удалось обновить значение", zap.Error(err))
		return false
	}
	return true
}

func (db *Database) GetGauge(nameMetric string) (valueMetric float64, exists bool) {
	err := db.Conn.QueryRow(context.Background(),
		`SELECT value
		FROM gauges
		WHERE name = $1`, nameMetric).Scan(&valueMetric)
	if err != nil {
		logger.Log.Warn("Не удалось получить значение", zap.Error(err))
		return valueMetric, false
	}
	return valueMetric, true
}

func (db *Database) GetCounter(nameMetric string) (valueMetric int64, exists bool) {
	err := db.Conn.QueryRow(context.Background(),
		`SELECT value
		FROM counters
		WHERE name = $1`, nameMetric).Scan(&valueMetric)
	if err != nil {
		logger.Log.Warn("Не удалось получить значение", zap.Error(err))
		return valueMetric, false
	}
	return valueMetric, true
}

func (db *Database) GetGauges() (map[string]float64, bool) {
	valuesMetric := make(map[string]float64)
	rows, err := db.Conn.Query(context.Background(),
		`SELECT name, value
		FROM gauges`)
	if err != nil {
		logger.Log.Warn("Не удалось получить значения", zap.Error(err))
		return valuesMetric, false
	}
	defer rows.Close()

	for rows.Next() {
		var nameMetric string
		var valueMetric float64
		err = rows.Scan(&nameMetric, &valueMetric)
		if err != nil {
			return valuesMetric, false
		}

		valuesMetric[nameMetric] = valueMetric
	}
	return valuesMetric, true
}

func (db *Database) GetCounters() (map[string]int64, bool) {
	valuesMetric := make(map[string]int64)
	rows, err := db.Conn.Query(context.Background(),
		`SELECT name, value
		FROM counters`)
	if err != nil {
		logger.Log.Warn("Не удалось получить значения", zap.Error(err))
		return valuesMetric, false
	}
	defer rows.Close()

	for rows.Next() {
		var nameMetric string
		var valueMetric int64
		err = rows.Scan(&nameMetric, &valueMetric)
		if err != nil {
			return valuesMetric, false
		}

		valuesMetric[nameMetric] = valueMetric
	}
	return valuesMetric, true
}
