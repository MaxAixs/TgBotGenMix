package telegram

import (
	"BotMixology/events/telegram/userState"
	"BotMixology/storage"
	"errors"
	"fmt"
	"strings"
)

func (p *Processor) SetTbcName(state *userState.UserState, text string) error {
	state.TobaccoName = strings.TrimSpace(text)

	if p.storage.TobaccoExists(state.TobaccoName) {
		return fmt.Errorf("табак с именем %s уже существует", state.TobaccoName)
	}

	return nil
}

func SetStrengthTbc(text string) (storage.Strength, error) {
	strength := storage.Strength(strings.TrimSpace(text))

	if !storage.IsValidStrength(strength) {
		return strength, errors.New("invalid strength")
	}

	return strength, nil
}

func (p *Processor) ExistTbc(state *userState.UserState, text string) error {
	state.TobaccoName = strings.TrimSpace(text)

	if !p.storage.TobaccoExists(state.TobaccoName) {
		return fmt.Errorf("табака с именем %s не сущестувет", state.TobaccoName)
	}

	return nil

}

func SetFlavorName(state *userState.UserState, text string) {
	state.FlavorName = strings.TrimSpace(text)
}

func HandleFlavorType(text string) (storage.Flavors, error) {
	flavors := storage.Flavors(strings.TrimSpace(text))

	if !storage.IsValidFlavorType(flavors) {
		return flavors, errors.New("invalid flavor type")
	}

	return flavors, nil
}
