package telegram

import (
	"BotMixology/client/telegram"
	"BotMixology/events/telegram/buttons"
	"BotMixology/events/telegram/cmd"
	"BotMixology/events/telegram/userState"
	"fmt"
	"log"
)

func (p *Processor) doCmd(state *userState.UserState, text string, chatID int) error {
	log.Printf("Обработка команды: %s", text)

	sendMsg := NewMessageSender(chatID, p.tg)
	log.Print(state.State)
	if state.IsActive() {
		return p.processState(state, text, sendMsg)
	}

	return p.processCommand(state, text, sendMsg)
}

func (p *Processor) processState(state *userState.UserState, text string, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	handler, exists := StateHandlers[state.State]
	if !exists {
		return sendMsg(cmd.UnknownCommand)
	}

	return handler(p, state, text, sendMsg)
}

func (p *Processor) processCommand(state *userState.UserState, text string, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	cmdActions := map[string]func(){
		cmd.StartCmd:           func() { sendMsg(cmd.MsgStart, buttons.CommandKeyboard()) },
		cmd.AddTobacco:         func() { p.handleAddTobacco(state, sendMsg) },
		cmd.DeleteTobacco:      func() { p.handleDeleteTobacco(state, sendMsg) },
		cmd.AddFlavor:          func() { p.handleAddFlavor(state, sendMsg) },
		cmd.DeleteFlavor:       func() { p.handleDeleteFlavor(state, sendMsg) },
		cmd.ShowTobaccoCatalog: func() { p.handleShowTobaccoCatalog(state, sendMsg) },
		cmd.GenerateMix:        func() { p.handleGenerateMix(state, sendMsg) },
		cmd.CreateMix:          func() { p.handleCreateMix(state, sendMsg) },
	}

	if action, exists := cmdActions[text]; exists {
		action()
		return nil
	}

	return sendMsg(cmd.UnknownCommand)
}

func (p *Processor) AddTobacco(state *userState.UserState, text string, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	switch state.State {
	case AddTbcName:
		err := p.SetTbcName(state, text)
		if err != nil {
			return sendMsg(cmd.TbcAlreadyExists)
		}

		state.SetState(AddTbcStrength)
		return sendMsg(cmd.EnterStrengthTbc)

	case AddTbcStrength:
		err := SetStrengthTbc(state, text)
		if err != nil {
			return sendMsg(cmd.InvalidStrength)
		}

		p.Storage.AddTobacco(state.TobaccoName, state.TobaccoStrength)
	}

	state.Reset()
	return sendMsg(cmd.SuccessAddTobacco, buttons.CommandKeyboard())
}

func (p *Processor) DeleteTobacco(state *userState.UserState, text string, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	switch state.State {
	case DeleteTbcName:
		err := p.ExistTbc(state, text)
		if err != nil {
			return sendMsg(cmd.TbcNotFound)
		}

		p.Storage.DeleteTobacco(state.TobaccoName)
	}

	state.Reset()
	return sendMsg(cmd.SuccessDeleteTobacco, buttons.CommandKeyboard())
}

func (p *Processor) AddFlavor(state *userState.UserState, text string, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	switch state.State {
	case SelectTbcForFlavor:
		err := p.ExistTbc(state, text)
		if err != nil {
			return sendMsg(cmd.TbcNotFound)
		}

		state.SetState(AddFlavorName)
		return sendMsg(cmd.EnterFlavorName)

	case AddFlavorName:
		err := p.SetFlavorName(state, text)
		if err != nil {
			return sendMsg(cmd.FlavorAlreadyExists)
		}

		state.SetState(AddFlavorType)
		return sendMsg(cmd.EnterFlavorType)

	case AddFlavorType:
		err := HandleFlavorType(state, text)
		if err != nil {
			return sendMsg(cmd.InvalidFlavorType)
		}

		err = p.Storage.AddFlavor(state.TobaccoName, state.FlavorName, state.FlavorType)
		if err != nil {
			return sendMsg(fmt.Sprintf("Ошибка: %v", err))
		}

		state.Reset()
		return sendMsg(cmd.SuccessAddFlavor, buttons.CommandKeyboard())
	}

	return nil
}

func (p *Processor) DeleteFlavor(state *userState.UserState, text string, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	switch state.State {
	case SelectTbcForFlavorDelete:
		err := p.ExistTbc(state, text)
		if err != nil {
			return sendMsg(cmd.TbcNotFound)
		}

		state.SetState(DeleteFlavorName)
		return sendMsg(cmd.EnterFlavorName)

	case DeleteFlavorName:
		err := p.ExistFlavor(state, text)
		if err != nil {
			return sendMsg(cmd.FlavorNotExists)
		}

		state.SetState(DeleteFlavorType)
		return sendMsg(cmd.EnterFlavorType)

	case DeleteFlavorType:
		err := HandleFlavorType(state, text)
		if err != nil {
			return sendMsg(cmd.InvalidFlavorType)
		}

		err = p.Storage.DeleteFlavor(state.TobaccoName, state.FlavorName)
		if err != nil {
			state.Reset()
			return sendMsg(fmt.Sprintf("Ошибка: %v", err))
		}

		state.Reset()
		return sendMsg(cmd.SuccessDeleteFlavor, buttons.CommandKeyboard())
	}

	return nil
}

func (p *Processor) ShowTobacco(state *userState.UserState, text string, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	switch state.State {
	case ShowTobaccoCatalog:
		err := p.ExistTbc(state, text)
		if err != nil {
			return sendMsg(cmd.TbcNotFound)
		}

		catalog := p.Storage.ShowTobaccoCatalog(state.TobaccoName)

		state.Reset()
		return sendMsg(catalog, buttons.CommandKeyboard())
	}

	return nil
}

func (p *Processor) GenerateMix(state *userState.UserState, text string, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	switch state.State {
	case SelectStrForMix:
		err := SetStrengthTbc(state, text)
		if err != nil {
			return sendMsg(cmd.InvalidStrength)
		}

		state.SetState(SelectFlavorType)
		return sendMsg(cmd.SelectFlavorType)

	case SelectFlavorType:
		err := HandleFlavorType(state, text)
		if err != nil {
			return sendMsg(cmd.InvalidFlavorType)
		}

		mixResult := p.Storage.GenerateMix(state.TobaccoStrength, state.FlavorType)

		state.Reset()
		return sendMsg(mixResult, buttons.CommandKeyboard())
	}

	return nil
}

func (p *Processor) CreateMix(state *userState.UserState, text string, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	switch state.State {
	case ChooseStrength:
		err := SetStrengthTbc(state, text)
		if err != nil {
			return sendMsg(cmd.InvalidStrength)
		}

		barOfTbc := p.Storage.GetTbcBar(state.TobaccoStrength)
		sendMsg(barOfTbc)

		state.SetState(ChooseTbc)
		return sendMsg(cmd.ChooseTobacco)

	case ChooseTbc:
		err := p.ExistTbc(state, text)
		if err != nil {
			return sendMsg(cmd.TbcNotFound)
		}

		state.SetState(GetAndChooseFlavors)
		fallthrough
	case GetAndChooseFlavors:
		flavors, err := p.Storage.GetFlavors(state.TobaccoName, state.TobaccoStrength)
		if err != nil {
			return sendMsg(cmd.FlavorsIsEmpty)
		}

		return p.chooseFlavors(flavors, text, state, sendMsg)
	}

	return nil
}

func (p *Processor) chooseFlavors(flavors map[string]string, text string, state *userState.UserState, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	var flavorList string

	for flavorName, flavorType := range flavors {
		flavorList += fmt.Sprintf("Вкус: %s, Тип: %s\n", flavorName, flavorType)
	}

	sendMsg(flavorList)

	if flavorType, exists := flavors[text]; exists {
		if contains(state.SelectedFlavors, text) {
			return sendMsg("Вы уже выбрали этот вкус, пожалуйста, выберите другой.")
		}
		state.SelectedFlavors = append(state.SelectedFlavors, text)
		sendMsg(fmt.Sprintf("Вы выбрали вкус: %s (Тип: %s)", text, flavorType))
	} else {
		return sendMsg("Такого вкуса нет в списке, попробуйте ещё раз.")
	}

	if len(state.SelectedFlavors) == 3 {
		return p.finishFlavorSelection(state, sendMsg)
	}

	return sendMsg(fmt.Sprintf("Вы выбрали: %v. Можете выбрать ещё %d вкуса.", state.SelectedFlavors, 3-len(state.SelectedFlavors)))
}

func NewMessageSender(chatID int, tg *telegram.Client) func(string, ...buttons.ReplyMarkUp) error {
	return func(msg string, markUps ...buttons.ReplyMarkUp) error {
		markUp := buttons.ReplyMarkUp{}
		if len(markUps) > 0 {
			markUp = markUps[0]
		}

		if markUp.InLineKeyBoard == nil {
			markUp.InLineKeyBoard = [][]buttons.InLineKeyBoardButton{}
		}

		log.Printf("Отправка сообщения: %s с клавиатурой: %v", msg, markUp)
		return tg.SendMessage(chatID, msg, markUp)
	}
}
