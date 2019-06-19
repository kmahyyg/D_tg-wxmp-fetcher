package db

import (
	"database/sql"

	_ "github.com/lib/pq" // Load the PostgreSQL driver
)

var db *sql.DB

// Connect makes a connection to current database
func Connect(driver string, source string) error {
	var err error
	db, err = sql.Open(driver, source)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	return nil
}
