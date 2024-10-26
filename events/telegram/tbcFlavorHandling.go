package telegram

import (
	"BotMixology/events/telegram/buttons"
	"BotMixology/events/telegram/userState"
	"BotMixology/storage"
	"errors"
	"fmt"
	"strings"
)

func (p *Processor) SetTbcName(state *userState.UserState, text string) error {
	state.TobaccoName = strings.TrimSpace(text)

	if p.Storage.TobaccoExists(state.TobaccoName) {
		return fmt.Errorf("табак с именем %s уже существует", state.TobaccoName)
	}

	return nil
}

func SetStrengthTbc(state *userState.UserState, text string) error {
	state.TobaccoStrength = storage.Strength(strings.TrimSpace(text))

	if !storage.IsValidStrength(state.TobaccoStrength) {
		return errors.New("invalid strength")
	}

	return nil
}

func (p *Processor) ExistTbc(state *userState.UserState, text string) error {
	state.TobaccoName = strings.TrimSpace(text)

	if !p.Storage.TobaccoExists(state.TobaccoName) {
		return fmt.Errorf("табака с именем %s не сущестувет", state.TobaccoName)
	}

	return nil

}

func (p *Processor) SetFlavorName(state *userState.UserState, text string) error {
	state.FlavorName = strings.TrimSpace(text)

	if p.Storage.FlavorExists(state.TobaccoName, state.FlavorName) {
		return errors.New("this flavor already exists")
	}

	return nil
}

func (p *Processor) ExistFlavor(state *userState.UserState, text string) error {
	state.FlavorName = strings.TrimSpace(text)

	if !p.Storage.FlavorExists(state.TobaccoName, state.FlavorName) {
		return errors.New("this flavor does not exist")
	}

	return nil
}

func HandleFlavorType(state *userState.UserState, text string) error {
	state.FlavorType = storage.Flavors(strings.TrimSpace(text))

	if !storage.IsValidFlavorType(state.FlavorType) {
		return errors.New("invalid flavor type")
	}

	return nil
}

func (p *Processor) finishFlavorSelection(state *userState.UserState, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	mixResult := fmt.Sprintf("Ваш микс: Табак '%s' с вкусами %v.", state.TobaccoName, state.SelectedFlavors)

	state.Reset()

	return sendMsg(mixResult, buttons.CommandKeyboard())
}

func contains(flavors []string, text string) bool {
	for _, flavor := range flavors {
		if flavor == text {
			return true
		}
	}

	return false
}
