package storage

type Storage interface {
	AddTobacco(tobaccoName string, strength Strength)
	DeleteTobacco(tobaccoName string)
	AddFlavor(tobaccoName, flavorName string, flavorType Flavors) error
	DeleteFlavor(tobaccoName, flavorName string, flavorType Flavors) error
	ShowTobaccoCatalog(tobaccoName string) string
	TobaccoExists(tobaccoName string) bool
	FlavorExists(tobaccoName BarOfTobacco, flavorName string, flavorType Flavors) bool
}

type Strength string

const (
	Light  Strength = "Light"
	Medium Strength = "Medium"
	Strong Strength = "Strong"
)

type Flavors string

const (
	Sour      Flavors = "Sour"
	Sweet     Flavors = "Sweet"
	SourSweet Flavors = "SourSweet"
)

type BarOfTobacco struct {
	Strength string
	Flavor   map[Flavors][]string
}

func IsValidStrength(strength Strength) bool {
	switch strength {
	case Light, Medium, Strong:
		return true
	default:
		return false
	}
}

func IsValidFlavorType(flavorType Flavors) bool {
	switch flavorType {
	case Sour, Sweet, SourSweet:
		return true
	default:
		return false
	}

}
