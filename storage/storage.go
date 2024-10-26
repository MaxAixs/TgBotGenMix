package storage

type Storage interface {
	TobaccoManager
	FlavorManager
	CatalogManager
	GenerateMix
}

type FlavorManager interface {
	AddFlavor(tobaccoName, flavorName string, flavorType Flavors) error
	DeleteFlavor(tobaccoName, flavorName string) error
	FlavorExists(tobaccoName string, flavorName string) bool
	GetFlavors(tobaccoName string, strength Strength) (map[string]string, error)
}

type CatalogManager interface {
	ShowTobaccoCatalog(tobaccoName string) string
}

type TobaccoManager interface {
	AddTobacco(tobaccoName string, strength Strength) error
	DeleteTobacco(tobaccoName string) error
	TobaccoExists(tobaccoName string) bool
	GetTbcBar(strength Strength) string
}

type GenerateMix interface {
	GenerateMix(strength Strength, flavorType Flavors) string
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
	Strength Strength
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
