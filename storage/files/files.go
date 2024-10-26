package files

import (
	"BotMixology/lib/e"
	"BotMixology/storage"
	"fmt"
	"math/rand"
	"time"
)

type Storage struct {
	tobaccos map[string]storage.BarOfTobacco
}

func NewStorage() *Storage {
	return &Storage{
		tobaccos: make(map[string]storage.BarOfTobacco),
	}
}

func (s *Storage) AddTobacco(tobaccoName string, strength storage.Strength) {

	s.tobaccos[tobaccoName] = storage.BarOfTobacco{
		Strength: storage.Strength(string(strength)),
		Flavor:   make(map[storage.Flavors][]string),
	}
}

func (s *Storage) DeleteTobacco(tobaccoName string) {
	delete(s.tobaccos, tobaccoName)
}

func (s *Storage) AddFlavor(tobaccoName string, flavorName string, flavorType storage.Flavors) error {
	tobacco := s.tobaccos[tobaccoName]

	if s.FlavorExists(tobacco, flavorName, flavorType) {
		return fmt.Errorf("вкус %s уже существует в типе %s у табака %s", flavorName, flavorType, tobaccoName)
	}

	if _, flavorExists := tobacco.Flavor[flavorType]; !flavorExists {
		tobacco.Flavor[flavorType] = []string{}
	}

	tobacco.Flavor[flavorType] = append(tobacco.Flavor[flavorType], flavorName)

	s.tobaccos[tobaccoName] = tobacco
	return nil
}

func (s *Storage) DeleteFlavor(tobaccoName string, flavorName string, flavorType storage.Flavors) error {
	tobacco := s.tobaccos[tobaccoName]

	flavors := tobacco.Flavor[flavorType]

	for i, flavor := range flavors {
		if flavor == flavorName {
			tobacco.Flavor[flavorType] = append(flavors[:i], flavors[i+1:]...)
			s.tobaccos[tobaccoName] = tobacco
			return nil
		}
	}

	return e.Wrap(fmt.Sprintf("У табака %s в типе вкуса %s нет вкуса %s", tobaccoName, flavorType, flavorName), nil)
}

func (s *Storage) ShowTobaccoCatalog(tobaccoName string) string {
	tobacco := s.tobaccos[tobaccoName]

	result := fmt.Sprintf("Каталог табака %s:\n", tobaccoName)
	result += fmt.Sprintf("Крепость: %s\n", tobacco.Strength)
	for flavorType, flavors := range tobacco.Flavor {
		result += fmt.Sprintf("Тип вкуса: %s - Вкусы: %v\n", flavorType, flavors)
	}

	return result
}

func (s *Storage) TobaccoExists(tobaccoName string) bool {
	_, exists := s.tobaccos[tobaccoName]
	return exists
}

func (s *Storage) FlavorExists(tobacco storage.BarOfTobacco, flavorName string, flavorType storage.Flavors) bool {
	if flavors, flavorExists := tobacco.Flavor[flavorType]; flavorExists {
		for _, existingFlavor := range flavors {
			if existingFlavor == flavorName {
				return true
			}
		}
	}

	return false
}

func (s *Storage) GenerateMix(strength storage.Strength, flavorType storage.Flavors) string {
	var mixResult string
	var allFlavors []string

	for tobaccoName, tobacco := range s.tobaccos {
		if tobacco.Strength == strength {
			if flavors, exists := tobacco.Flavor[flavorType]; exists && len(flavors) > 0 {
				mixResult += fmt.Sprintf("Табак: %s\n", tobaccoName)
				mixResult += fmt.Sprintf("Крепость: %s\n", strength)
				mixResult += fmt.Sprintf("Вкусы: (%s): %v\n", flavorType, flavors)

				allFlavors = append(allFlavors, flavors...)
			}
		}
	}

	if len(allFlavors) == 0 {
		return fmt.Sprintf("Нет табаков с такой крепостью '%s' и типом вкуса '%s'", strength, flavorType)
	}

	mix := generateFlavorMix(allFlavors)

	mixResult += fmt.Sprintf("\nСгенерированный микс вкусов: %v", mix)

	return mixResult
}

func generateFlavorMix(flavors []string) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	mixSize := 2

	if len(flavors) < mixSize {
		return flavors
	}

	flavor1 := r.Intn(len(flavors))
	flavor2 := r.Intn(len(flavors))

	for flavor1 == flavor2 {
		flavor2 = r.Intn(len(flavors))
	}

	return []string{flavors[flavor1], flavors[flavor2]}
}
