package telegram

import (
	"BotMixology/client/telegram"
	"BotMixology/events"
	"BotMixology/events/telegram/userState"
	"BotMixology/lib/e"
	"BotMixology/storage/sqlite"
	"errors"
)

type Processor struct {
	tg         *telegram.Client
	offset     int
	Storage    sqlite.Storage
	userStates map[int]*userState.UserState
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrorUnknownEventType = errors.New("unknown event")
	ErrorUnknownMetaType  = errors.New("unknown meta type")
)

func NewProcessor(client *telegram.Client, storage sqlite.Storage) *Processor {
	return &Processor{
		tg:         client,
		Storage:    storage,
		userStates: make(map[int]*userState.UserState),
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("cant get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	case events.СlickBtn:
		return p.processBtn(event)
	default:
		return e.Wrap("cant process Message", ErrorUnknownEventType)
	}
}

func (p *Processor) processBtn(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("cant process button", ErrorUnknownEventType)
	}

	state, exist := p.userStates[meta.ChatID]
	if !exist {
		state = &userState.UserState{}
		p.userStates[meta.ChatID] = state
	}

	err = p.doCmd(state, event.Text, meta.ChatID)

	return err
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("cant process Message", err)
	}

	state, exists := p.userStates[meta.ChatID]
	if !exists {
		state = &userState.UserState{}
		p.userStates[meta.ChatID] = state

	}

	return p.doCmd(state, event.Text, meta.ChatID)
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("Cant get meta", ErrorUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	switch updType {

	case events.Message:
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}

	case events.СlickBtn:
		res.Meta = Meta{
			ChatID:   upd.CallbackQuery.Message.Chat.ID,
			Username: upd.CallbackQuery.Message.From.Username,
		}
	default:
		e.Wrap("Cant event type", ErrorUnknownEventType)
	}

	return res
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message != nil {
		return events.Message
	}

	if upd.CallbackQuery != nil {
		return events.СlickBtn
	}

	return events.Unknown
}

func fetchText(upd telegram.Update) string {
	if upd.Message != nil {
		return upd.Message.Text
	}

	if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Data
	}
	return ""
}
