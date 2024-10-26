package sqlite

import (
	"database/sql"
	"fmt"
)

func migrate(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS tobacco (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tobaccoName TEXT UNIQUE NOT NULL,
			strength TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS flavor (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tobacco_id INTEGER,
			flavorName TEXT NOT NULL,
			flavorType TEXT NOT NULL,
			FOREIGN KEY(tobacco_id) REFERENCES tobacco(id) ON DELETE CASCADE
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("cant create table: %w", err)
		}
	}

	return nil
}
