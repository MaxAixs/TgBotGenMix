package telegram

import (
	"BotMixology/client/telegram"
	"BotMixology/events/telegram/buttons"
	"BotMixology/events/telegram/userState"
	"fmt"
	"log"
)

func (p *Processor) doCmd(state *userState.UserState, text string, chatID int) error {
	log.Printf("Обработка команды: %s", text)

	sendMsg := NewMessageSender(chatID, p.tg)

	if state.IsActive() {
		return p.processState(state, text, sendMsg)
	}

	return p.processCommand(state, text, sendMsg)
}

func (p *Processor) processState(state *userState.UserState, text string, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	handler, exists := StateHandlers[state.State]
	if !exists {
		return sendMsg(unknownCommand)
	}

	return handler(p, state, text, sendMsg)
}

func (p *Processor) processCommand(state *userState.UserState, text string, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	switch text {
	case startCmd:
		return sendMsg(msgStart, buttons.CommandKeyboard())
	case addTobacco:
		state.SetState(AddTbcName)
		return sendMsg(enterNameTbc)
	case deleteTobacco:
		state.SetState(DeleteTbcName)
		return sendMsg(deleteTbc)
	case addFlavor:
		state.SetState(SelectTbcForFlavor)
		return sendMsg(addFlavorTbc)
	case deleteFlavor:
		state.SetState(SelectTbcForFlavorDelete)
		return sendMsg(deleteFlavorTbc)
	case showTobaccoCatalog:
		state.SetState(ShowTbcCatalog)
		return sendMsg(showTbcCatalog)
	default:
		return sendMsg(unknownCommand)
	}
}

func (p *Processor) AddTobacco(state *userState.UserState, text string, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	switch state.State {
	case AddTbcName:
		if err := p.SetTbcName(state, text); err != nil {
			return sendMsg(TbcAlreadyExists)
		}

		state.SetState(AddTbcStrength)
		return sendMsg(EnterStrengthTbc)

	case AddTbcStrength:
		strength, err := SetStrengthTbc(text)
		if err != nil {
			return sendMsg(invalidStrength)
		}

		p.storage.AddTobacco(state.TobaccoName, strength)
	}

	state.Reset()
	return sendMsg(successAddTobacco, buttons.CommandKeyboard())
}

func (p *Processor) DeleteTobacco(state *userState.UserState, text string, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	switch state.State {
	case DeleteTbcName:
		if err := p.ExistTbc(state, text); err != nil {
			return sendMsg(TbcNotFound)
		}

		p.storage.DeleteTobacco(state.TobaccoName)
	}

	state.Reset()
	return sendMsg(SuccessDeleteTobacco, buttons.CommandKeyboard())
}

func (p *Processor) AddFlavor(state *userState.UserState, text string, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	switch state.State {
	case SelectTbcForFlavor:
		if err := p.ExistTbc(state, text); err != nil {
			return sendMsg(TbcNotFound)
		}

		state.SetState(AddFlavorName)
		return sendMsg(enterFlavorName)

	case AddFlavorName:
		SetFlavorName(state, text)

		state.SetState(AddFlavorType)
		return sendMsg(enterFlavorType)

	case AddFlavorType:
		flavorType, err := HandleFlavorType(text)
		if err != nil {
			return sendMsg(invalidFlavorType)
		}

		err = p.storage.AddFlavor(state.TobaccoName, state.FlavorName, flavorType)
		if err != nil {
			return sendMsg(fmt.Sprintf("Ошибка: %v", err))
		}

		state.Reset()
		return sendMsg(successAddFlavor, buttons.CommandKeyboard())
	}

	return nil
}

func (p *Processor) DeleteFlavor(state *userState.UserState, text string, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	switch state.State {
	case SelectTbcForFlavorDelete:
		if err := p.ExistTbc(state, text); err != nil {
			return sendMsg(TbcNotFound)
		}

		state.SetState(DeleteFlavorName)
		return sendMsg(enterFlavorName)

	case DeleteFlavorName:
		SetFlavorName(state, text)

		state.SetState(DeleteFlavorType)
		return sendMsg(enterFlavorType)

	case DeleteFlavorType:
		flavorType, err := HandleFlavorType(text)
		if err != nil {
			return err
		}

		if err := p.storage.DeleteFlavor(state.TobaccoName, state.FlavorName, flavorType); err != nil {
			state.Reset()
			return sendMsg(fmt.Sprintf("Ошибка: %v", err))
		}

		state.Reset()
		return sendMsg(successDeleteFlavor, buttons.CommandKeyboard())
	}

	return nil
}

func (p *Processor) ShowTobacco(state *userState.UserState, text string, sendMsg func(string, ...buttons.ReplyMarkUp) error) error {
	switch state.State {
	case ShowTbcCatalog:
		if err := p.ExistTbc(state, text); err != nil {
			return sendMsg(TbcNotFound)
		}

		catalog := p.storage.ShowTobaccoCatalog(state.TobaccoName)

		state.Reset()
		return sendMsg(catalog, buttons.CommandKeyboard())
	}

	return nil
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
