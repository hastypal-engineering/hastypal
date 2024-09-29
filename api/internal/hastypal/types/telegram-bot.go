package types

type TelegramBotMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

type TelegramBotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type TelegramWebhook struct {
	Url string `json:"url"`
}

type AdminTelegramBotSetup struct {
	Commands []TelegramBotCommand `json:"commands"`
	Webhook  TelegramWebhook      `json:"webhook"`
}
