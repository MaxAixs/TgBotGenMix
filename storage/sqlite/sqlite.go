package sqlite

import (
	"database/sql"
	"fmt"
)

// Storage represents the database for storing tobacco and flavor information.
type Storage struct {
	db *sql.DB
}

// NewDb opens the database at the specified path, checks the connection, and performs migrations.
func NewDb(path string) (*Storage, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("cant open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cant connect database: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("cant migrate database: %w", err)
	}

	return &Storage{db: db}, nil
}
