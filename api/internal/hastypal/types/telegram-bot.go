package types

type TelegramBotMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}
