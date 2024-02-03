package main

import (
	"context"
	"fmt"

	"github.com/drrhaos/metrics/internal/logger"
	"github.com/jackc/pgx/v5"
)

type Database struct {
	Conn *pgx.Conn
	Uri  string
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
		fmt.Println(err)
		return false
	}
	return true
}

func (db *Database) Connect(uri string) error {
	conn, err := pgx.Connect(context.Background(), uri)
	if err != nil {
		return err
	}
	db.Conn = conn
	db.Uri = uri
	logger.Log.Info("Соединение с базой успешно установлено")
	return nil
}
