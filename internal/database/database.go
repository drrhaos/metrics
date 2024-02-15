package database

import (
	"context"
	"time"

	"github.com/avast/retry-go"
	"github.com/drrhaos/metrics/internal/logger"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var sleepStep = map[uint]int64{0: 1, 1: 3, 2: 5}

type Database struct {
	Conn *pgxpool.Pool
}

func NewDatabase(uri string) *Database {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		logger.Log.Panic("Ошибка при парсинге конфигурации:", zap.Error(err))
		return nil
	}
	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Log.Panic("Не удалось подключиться к базе данных")
		return nil
	}
	db := &Database{Conn: conn}
	err = db.Migrations(ctx)
	if err != nil {
		logger.Log.Panic("Не удалось подключиться к базе данных")
		return nil
	}
	return db
}

func (db *Database) SaveMetrics(filePath string) bool {
	if db != nil {
		return true
	} else {
		return false
	}
}

func (db *Database) LoadMetrics(filePath string) bool {
	if db != nil {
		return true
	} else {
		return false
	}
}

func (db *Database) Close() {
	db.Conn.Close()
}

func (db *Database) Ping(ctx context.Context) bool {
	if err := db.Conn.Ping(ctx); err != nil {
		return false
	}
	return true
}

func (db *Database) Migrations(ctx context.Context) error {
	var exist bool
	err := db.Conn.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'counters')`).Scan(&exist)
	if err != nil {
		return err
	}
	if !exist {
		_, err := db.Conn.Exec(ctx,
			`CREATE TABLE counters (
				"id" SERIAL PRIMARY KEY,
				"name" character(40) NOT NULL UNIQUE,
				"value" bigint NOT NULL
			)`)
		if err != nil {
			return err
		}
	}

	err = db.Conn.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'gauges')`).Scan(&exist)
	if err != nil {
		return err
	}
	if !exist {
		_, err := db.Conn.Exec(ctx,
			`CREATE TABLE gauges (
				"id" SERIAL PRIMARY KEY,
				"name" character(40) NOT NULL UNIQUE,
				"value" double precision NOT NULL
			)`)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) UpdateGauge(ctx context.Context, nameMetric string, valueMetric float64) bool {
	err := retry.Do(
		func() error {
			_, err := db.Conn.Exec(ctx,
				`INSERT INTO gauges (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE
		SET value = EXCLUDED.value`, nameMetric, valueMetric)
			if err != nil {
				logger.Log.Warn("Не удалось обновить значение", zap.Error(err))
				return err
			}
			return nil
		},
		retry.Attempts(3),
		retry.DelayType(customDelay()),
	)
	if err != nil {
		logger.Log.Warn("Не удалось обновить значение", zap.Error(err))
		return false
	}
	return true
}

func (db *Database) UpdateCounter(ctx context.Context, nameMetric string, valueMetric int64) bool {
	currentValue, exist := db.GetCounter(ctx, nameMetric)
	if exist {
		valueMetric += currentValue
	}
	err := retry.Do(
		func() error {
			_, err := db.Conn.Exec(ctx,
				`INSERT INTO counters (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE
		SET value = EXCLUDED.value`, nameMetric, valueMetric)
			if err != nil {
				logger.Log.Warn("Не удалось обновить значение", zap.Error(err))
				return err
			}
			return nil
		},
		retry.Attempts(3),
		retry.DelayType(customDelay()),
	)
	if err != nil {
		logger.Log.Warn("Не удалось обновить значение", zap.Error(err))
		return false
	}
	return true
}

func (db *Database) GetGauge(ctx context.Context, nameMetric string) (valueMetric float64, exists bool) {
	err := retry.Do(
		func() error {
			err := db.Conn.QueryRow(ctx,
				`SELECT value
		FROM gauges
		WHERE name = $1`, nameMetric).Scan(&valueMetric)
			if err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(3),
		retry.DelayType(customDelay()),
	)
	if err != nil {
		logger.Log.Warn("Не удалось получить значение", zap.Error(err))
		return valueMetric, false
	}
	return valueMetric, true
}

func (db *Database) GetCounter(ctx context.Context, nameMetric string) (valueMetric int64, exists bool) {
	err := retry.Do(
		func() error {
			err := db.Conn.QueryRow(ctx,
				`SELECT value
		FROM counters
		WHERE name = $1`, nameMetric).Scan(&valueMetric)

			if err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(3),
		retry.DelayType(customDelay()),
	)
	if err != nil {
		logger.Log.Warn("Не удалось получить значение", zap.Error(err))
		return valueMetric, false
	}
	return valueMetric, true
}

func (db *Database) GetGauges(ctx context.Context) (map[string]float64, bool) {
	valuesMetric := make(map[string]float64)
	err := retry.Do(
		func() error {
			rows, err := db.Conn.Query(ctx,
				`SELECT name, value
				FROM gauges`)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var nameMetric string
				var valueMetric float64
				err = rows.Scan(&nameMetric, &valueMetric)
				if err != nil {
					return err
				}

				valuesMetric[nameMetric] = valueMetric
			}

			return nil
		},
		retry.Attempts(3),
		retry.DelayType(customDelay()),
	)
	if err != nil {
		logger.Log.Warn("Не удалось получить значение", zap.Error(err))
		return valuesMetric, false
	}

	return valuesMetric, true
}

func (db *Database) GetCounters(ctx context.Context) (map[string]int64, bool) {
	valuesMetric := make(map[string]int64)
	err := retry.Do(
		func() error {
			rows, err := db.Conn.Query(ctx,
				`SELECT name, value
		FROM counters`)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var nameMetric string
				var valueMetric int64
				err = rows.Scan(&nameMetric, &valueMetric)
				if err != nil {
					return err
				}

				valuesMetric[nameMetric] = valueMetric
			}
			return nil
		},
		retry.Attempts(3),
		retry.DelayType(customDelay()),
	)
	if err != nil {
		logger.Log.Warn("Не удалось получить значение", zap.Error(err))
		return valuesMetric, false
	}
	return valuesMetric, true
}

func customDelay() retry.DelayTypeFunc {
	return func(n uint, err error, config *retry.Config) time.Duration {
		return time.Duration(sleepStep[n])
	}
}
