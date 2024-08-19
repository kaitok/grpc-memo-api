package db

import (
	"github.com/gocraft/dbr/v2"
	_ "github.com/lib/pq" // Import PostgreSQL driver
)

// NewDB initializes the database connection.
func NewDB(dataSourceName string) (*dbr.Connection, error) {
	// Open database connection
	conn, err := dbr.Open("postgres", dataSourceName, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
