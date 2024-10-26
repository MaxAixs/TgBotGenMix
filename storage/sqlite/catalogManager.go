package sqlite

import (
	"BotMixology/lib/e"
	"BotMixology/storage"
	"fmt"
	"strings"
)

type CatalogManager interface {
	ShowTobaccoCatalog(tobaccoName string) string
}

// ShowTobaccoCatalog displays the catalog of a tobacco and all its associated flavors.
func (s *Storage) ShowTobaccoCatalog(tobaccoName string) string {
	var result strings.Builder
	var strength storage.Strength

	q := `SELECT strength FROM tobacco WHERE tobaccoName = ?`
	err := s.db.QueryRow(q, tobaccoName).Scan(&strength)
	if err = e.CheckErr("can't show tobacco", err); err != nil {
		return err.Error()
	}

	result.WriteString(fmt.Sprintf("Catalog for tobacco: %s\n", tobaccoName))
	result.WriteString(fmt.Sprintf("Tobacco strength: %s\n", strength))

	rows, err := s.db.Query(`SELECT flavorName, flavorType FROM flavor WHERE tobacco_id = (SELECT id FROM tobacco WHERE tobaccoName = ?)`, tobaccoName)
	if err = e.CheckErr("can't get flavors", err); err != nil {
		return err.Error()
	}
	defer rows.Close()

	for rows.Next() {
		var flavorName, flavorType string
		if err := rows.Scan(&flavorName, &flavorType); err != nil {
			return e.CheckErr("error scanning flavors", err).Error()
		}

		result.WriteString(fmt.Sprintf("Flavor: %s (Type: %s)\n", flavorName, flavorType))
	}

	return result.String()
}
