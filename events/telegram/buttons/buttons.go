package buttons

type ReplyMarkUp struct {
	InLineKeyBoard [][]InLineKeyBoardButton `json:"inline_keyboard"`
}

type InLineKeyBoardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

func CommandKeyboard() ReplyMarkUp {
	return ReplyMarkUp{
		InLineKeyBoard: [][]InLineKeyBoardButton{
			{
				{Text: "Добавить табак", CallbackData: "addTobacco"},
				{Text: "Удалить табак", CallbackData: "deleteTobacco"},
			},
			{
				{Text: "Добавить вкус", CallbackData: "addFlavor"},
				{Text: "Удалить вкус", CallbackData: "deleteFlavor"},
			},
			{
				{Text: "Показать каталог табака", CallbackData: "showTobaccoCatalog"},
			},
		},
	}
}
