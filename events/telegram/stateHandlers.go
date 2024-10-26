package telegram

import (
	"BotMixology/events/telegram/buttons"
	"BotMixology/events/telegram/userState"
)

const (
	AddTbcName               = "stateAddTobaccoName"
	AddTbcStrength           = "stateAddTobaccoStrength"
	DeleteTbcName            = "stateDeleteTobaccoName"
	SelectTbcForFlavor       = "stateSelectTobaccoForFlavor"
	AddFlavorName            = "stateAddFlavorName"
	AddFlavorType            = "stateAddFlavorType"
	SelectTbcForFlavorDelete = "stateSelectTobaccoForFlavorDelete"
	DeleteFlavorName         = "stateDeleteFlavorName"
	DeleteFlavorType         = "stateDeleteFlavorType"
	ShowTobaccoCatalog       = "stateShowTobaccoCatalog"
	SelectStrForMix          = "stateSelectStrengthForMix"
	SelectFlavorType         = "stateSelectFlavorTypeForMix"
	ChooseStrength           = "stateChooseStrength"
	ChooseTbc                = "stateChooseTbc"
	GetAndChooseFlavors      = "stateChooseFlavors"
)

var StateHandlers = map[string]func(*Processor, *userState.UserState, string, func(string, ...buttons.ReplyMarkUp) error) error{
	AddTbcName:               (*Processor).AddTobacco,
	AddTbcStrength:           (*Processor).AddTobacco,
	DeleteTbcName:            (*Processor).DeleteTobacco,
	SelectTbcForFlavor:       (*Processor).AddFlavor,
	AddFlavorName:            (*Processor).AddFlavor,
	AddFlavorType:            (*Processor).AddFlavor,
	SelectTbcForFlavorDelete: (*Processor).DeleteFlavor,
	DeleteFlavorName:         (*Processor).DeleteFlavor,
	DeleteFlavorType:         (*Processor).DeleteFlavor,
	ShowTobaccoCatalog:       (*Processor).ShowTobacco,
	SelectStrForMix:          (*Processor).GenerateMix,
	SelectFlavorType:         (*Processor).GenerateMix,
	ChooseStrength:           (*Processor).CreateMix,
	ChooseTbc:                (*Processor).CreateMix,
	GetAndChooseFlavors:      (*Processor).CreateMix,
}
