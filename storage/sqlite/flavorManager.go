package sqlite

import (
	"BotMixology/lib/e"
	"BotMixology/storage"
	"fmt"
)

// AddFlavor adds a new flavor associated with the specified tobacco.
func (s *Storage) AddFlavor(tobaccoName string, FlavorName string, flavorType storage.Flavors) error {
	var tobaccoID int
	err := s.db.QueryRow(`SELECT id FROM tobacco WHERE tobaccoName = ?`, tobaccoName).Scan(&tobaccoID)
	if err = e.CheckErr("tobacco not found", err); err != nil {
		return fmt.Errorf("tobacco not found %w", err)
	}

	q := `INSERT INTO flavor(tobacco_id,flavorName,flavorType) VALUES (?,?,?)`
	_, err = s.db.Exec(q, tobaccoID, FlavorName, flavorType)
	return e.CheckErr("cant add flavor", err)

}

// DeleteFlavor removes a flavor associated with the specified tobacco.
func (s *Storage) DeleteFlavor(tobaccoName string, flavorName string) error {
	var tobaccoID int
	err := s.db.QueryRow(`SELECT id FROM tobacco WHERE tobaccoName = ?`, tobaccoName).Scan(&tobaccoID)
	if err = e.CheckErr("tobacco not found", err); err != nil {
		return err
	}

	q := `DELETE FROM flavor WHERE tobacco_id = ? AND flavorName = ?`
	_, err = s.db.Exec(q, tobaccoID, flavorName)
	return e.CheckErr("can't delete Flavor", err)
}

// FlavorExists checks if a specific flavor exists for the specified tobacco.
func (s *Storage) FlavorExists(tobaccoName string, flavorName string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM flavor WHERE tobacco_id = (SELECT id FROM tobacco WHERE tobaccoName = ?) AND flavorName = ?)`
	err := s.db.QueryRow(query, tobaccoName, flavorName).Scan(&exists)
	if e.CheckErr("Cant check flavor existence", err) != nil {
		return false
	}

	return exists
}

// GetFlavors retrieves a map of flavors and their types for a specified tobacco and strength.
func (s *Storage) GetFlavors(tobaccoName string, strength storage.Strength) (map[string]string, error) {
	flavors := make(map[string]string)
	var tobaccoID int
	query := `SELECT id FROM tobacco WHERE tobaccoName = ? AND strength = ?`
	err := s.db.QueryRow(query, tobaccoName, strength).Scan(&tobaccoID)
	if err = e.CheckErr(fmt.Sprintf("can't find tobacco with name %s and strength %d", tobaccoName, strength), err); err != nil {
		return nil, err
	}

	query = `SELECT flavorName, flavorType FROM flavor WHERE tobacco_id = ?`
	rows, err := s.db.Query(query, tobaccoID)
	if err = e.CheckErr(fmt.Sprintf("can't get flavors by tobaccoName %s", tobaccoName), err); err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var flavorName, flavorType string
		if err = rows.Scan(&flavorName, &flavorType); err != nil {
			return nil, e.CheckErr(fmt.Sprintf("can't scan flavors by tobaccoName %s", tobaccoName), err)
		}

		flavors[flavorName] = flavorType
	}

	return flavors, nil
}
