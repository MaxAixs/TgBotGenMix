package userState

import (
	"BotMixology/storage"
)

type UserState struct {
	State           string
	TobaccoName     string
	TobaccoStrength storage.Strength
	FlavorName      string
	FlavorType      storage.Flavors
	SelectedFlavors []string
}

func (state *UserState) Reset() {
	state.State = ""
	state.TobaccoName = ""
	state.FlavorName = ""
	state.TobaccoStrength = ""
	state.FlavorType = ""
	state.SelectedFlavors = []string{}
}

func (state *UserState) SetState(newState string) {
	state.State = newState
}

func (state *UserState) IsActive() bool {
	return state.State != ""
}
