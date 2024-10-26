package telegram

import (
	"BotMixology/events/telegram/buttons"
	"BotMixology/events/telegram/cmd"
	"BotMixology/events/telegram/userState"
)

func (p *Processor) handleAddTobacco(state *userState.UserState, sendMsg func(string, ...buttons.ReplyMarkUp) error) {
	state.SetState(AddTbcName)
	sendMsg(cmd.EnterNameTbc)
}

func (p *Processor) handleDeleteTobacco(state *userState.UserState, sendMsg func(string, ...buttons.ReplyMarkUp) error) {
	state.SetState(DeleteTbcName)
	sendMsg(cmd.DeleteTbc)
}

func (p *Processor) handleAddFlavor(state *userState.UserState, sendMsg func(string, ...buttons.ReplyMarkUp) error) {
	state.SetState(SelectTbcForFlavor)
	sendMsg(cmd.AddFlavorTbc)
}

func (p *Processor) handleDeleteFlavor(state *userState.UserState, sendMsg func(string, ...buttons.ReplyMarkUp) error) {
	state.SetState(SelectTbcForFlavorDelete)
	sendMsg(cmd.DeleteFlavorTbc)
}

func (p *Processor) handleShowTobaccoCatalog(state *userState.UserState, sendMsg func(string, ...buttons.ReplyMarkUp) error) {
	state.SetState(ShowTobaccoCatalog)
	sendMsg(cmd.ShowTbcCatalog)
}

func (p *Processor) handleGenerateMix(state *userState.UserState, sendMsg func(string, ...buttons.ReplyMarkUp) error) {
	state.SetState(SelectStrForMix)
	sendMsg(cmd.SelectStr)
}

func (p *Processor) handleCreateMix(state *userState.UserState, sendMsg func(string, ...buttons.ReplyMarkUp) error) {
	state.SetState(ChooseStrength)
	sendMsg(cmd.SelectStr)
}
