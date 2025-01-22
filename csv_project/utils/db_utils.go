package utils

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func ConnectDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://postgres:password@db:5432/csvdb?sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(25)           // Max open connections
	db.SetMaxIdleConns(10)           // Max idle connections
	db.SetConnMaxLifetime(time.Hour) // Max connection lifetime

	return db, nil
}

func EnsureTableExistsAndTruncate(db *sql.DB) error {
	_, err := db.Exec(`TRUNCATE TABLE devices RESTART IDENTITY CASCADE`)
	if err != nil && strings.Contains(err.Error(), "relation \"devices\" does not exist") {
		_, err = db.Exec(`
			CREATE TABLE devices (
				id SERIAL PRIMARY KEY,
				devicename VARCHAR(100),
				devicetype VARCHAR(100),
				brand VARCHAR(100),
				model VARCHAR(100),
				os VARCHAR(100),
				osversion VARCHAR(100),
				purchasedate VARCHAR(100),
				warrantyend VARCHAR(100),
				status VARCHAR(100),
				price FLOAT
			)`)
	}
	return err
}
