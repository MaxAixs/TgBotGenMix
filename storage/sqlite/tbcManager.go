package sqlite

import (
	"BotMixology/lib/e"
	"BotMixology/storage"
	"fmt"
	"strings"
)

// AddTobacco adds a new tobacco entry to the database with the given name and strength.
func (s *Storage) AddTobacco(tobaccoName string, strength storage.Strength) error {
	q := `INSERT INTO tobacco (tobaccoName, strength) VALUES (?,?)`
	_, err := s.db.Exec(q, tobaccoName, strength)
	return e.CheckErr("can't delete Tobacco", err)
}

// DeleteTobacco removes a tobacco entry from the database by its name.
func (s *Storage) DeleteTobacco(tobaccoName string) error {
	q := `DELETE FROM tobacco WHERE TobaccoName = ?`
	_, err := s.db.Exec(q, tobaccoName)
	return e.CheckErr("can't delete Tobacco", err)
}

// GetTobaccoBar Get all tobaccos of a specific strength, returning their names as a formatted string.
func (s *Storage) GetTobaccoBar(strength storage.Strength) string {
	var result strings.Builder
	query := `SELECT tobaccoName FROM tobacco WHERE strength = ?`
	rows, err := s.db.Query(query, strength)
	if err != nil {
		return fmt.Sprintf("cant find tobacco with strength %d", strength)
	}
	defer rows.Close()

	for rows.Next() {
		var tobaccoName string
		if err := rows.Scan(&tobaccoName); err != nil {
			return fmt.Sprintf("cant get tobacco with strength %d", strength)
		}

		result.WriteString("Табак:" + tobaccoName + "\n")
	}

	return result.String()
}

// TobaccoExists checks if a tobacco with the specified name exists in the database.
func (s *Storage) TobaccoExists(tobaccoName string) bool {
	var exists bool
	q := `SELECT EXISTS(SELECT 1 FROM tobacco WHERE tobaccoName=?)`
	err := s.db.QueryRow(q, tobaccoName).Scan(&exists)
	if err = e.CheckErr("can't check tobacco existence", err); err != nil {
		return false
	}

	return exists
}
