package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func (db *DB) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS projects (
		id TEXT PRIMARY KEY,
		slug TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	`

	_, err := db.conn.Exec(query)
	return err
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) GetConnection() *sql.DB {
	return db.conn
}
