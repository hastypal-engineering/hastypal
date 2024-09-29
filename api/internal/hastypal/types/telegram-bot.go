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

type TelegramMessageUpdate struct {
	MessageId string `json:"message_id"`
	From      struct {
		Id           string `json:"id"`
		IsBot        bool   `json:"is_bot"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		LanguageCode string `json:"language_code"`
	}
	Chat struct {
		Id        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Type      string `json:"type"`
	}
	Date string `json:"date"`
	Text string `json:"text"`
}

type TelegramUpdate struct {
	UpdateId string                `json:"update_id"`
	Message  TelegramMessageUpdate `json:"message"`
}
