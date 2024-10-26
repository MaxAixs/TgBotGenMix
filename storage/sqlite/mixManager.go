package sqlite

import (
	"BotMixology/storage"
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

// GenerateMix generates a mix by selecting a tobacco and two different flavors based on the specified strength and flavor type.
func (s *Storage) GenerateMix(strength storage.Strength, flavorType storage.Flavors) string {
	query := `SELECT t.tobaccoName, f.flavorName
				FROM flavor f
				JOIN tobacco t ON f.tobacco_id = t.id
				WHERE t.strength = ? AND f.FlavorType = ?`

	rows, err := s.db.Query(query, strength, flavorType)
	if err != nil {
		return fmt.Sprintf("cant get flavors: %v", err)
	}
	defer rows.Close()

	flavors, err := collectFlavors(rows)
	if err != nil {
		return fmt.Sprintf("error collecting flavors: %v", err)
	}

	if len(flavors) < 2 {
		return fmt.Sprintf("not enough flavors, to generate a mix")
	}

	tobacco, flavor1, flavor2 := RandomMix(flavors)

	return fmt.Sprintf("Tobacco: %s: %s and %s", tobacco, flavor1, flavor2)
}

func RandomMix(f map[string][]string) (string, string, string) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	tobaccos := make([]string, 0, len(f))
	for tobacco := range f {
		tobaccos = append(tobaccos, tobacco)
	}

	selectedTbc := tobaccos[rand.Intn(len(tobaccos))]
	flavors := f[selectedTbc]

	firstFlavor := flavors[rand.Intn(len(flavors))]
	secondFlavor := flavors[rand.Intn(len(flavors))]

	for secondFlavor == firstFlavor {
		secondFlavor = flavors[rand.Intn(len(flavors))]
	}

	return selectedTbc, firstFlavor, secondFlavor

}

func collectFlavors(rows *sql.Rows) (map[string][]string, error) {
	tobaccoFlavors := make(map[string][]string)
	var tobaccoName, flavorName string
	for rows.Next() {
		if err := rows.Scan(&tobaccoName, &flavorName); err != nil {
			return nil, fmt.Errorf("can't scan flavors: %v", err)
		}
		tobaccoFlavors[tobaccoName] = append(tobaccoFlavors[tobaccoName], flavorName)

	}

	return tobaccoFlavors, nil
}
