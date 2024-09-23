package telegram

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID            int              `json:"update_id"`
	Message       *IncomingMessage `json:"message"`
	CallbackQuery *CallBackQuery   `json:"callback_query"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From
	Chat Chat
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type CallBackQuery struct {
	Data    string           `json:"data"`
	Message *IncomingMessage `json:"message"`
}
