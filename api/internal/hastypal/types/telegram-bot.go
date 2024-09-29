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

type TelegramUser struct {
	Id           string `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	LanguageCode string `json:"language_code"`
}

type TelegramChat struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Type      string `json:"type"`
}

type TelegramChatMemberAdministrator struct {
	Status            string `json:"status"`
	CanManageChat     bool   `json:"can_manage_chat"`
	CanChangeInfo     bool   `json:"can_change_info"`
	CanPostMessages   bool   `json:"can_post_messages"`
	CanEditMessages   bool   `json:"can_edit_messages"`
	CanDeleteMessages bool   `json:"can_delete_messages"`
	CanInviteUsers    bool   `json:"can_invite_users"`
	CanPostStories    bool   `json:"can_post_stories"`
	CanEditStories    bool   `json:"can_edit_stories"`
}

type TelegramMessageUpdate struct {
	MessageId string       `json:"message_id"`
	From      TelegramUser `json:"from"`
	Chat      TelegramChat `json:"chat"`
	Date      string       `json:"date"`
	Text      string       `json:"text"`
}

type BotMemberUpdated struct {
	Chat          TelegramChat                    `json:"chat"`
	From          TelegramUser                    `json:"from"`
	Date          string                          `json:"date"`
	NewChatMember TelegramChatMemberAdministrator `json:"new_chat_member"`
}

type TelegramUpdate struct {
	UpdateId     string                `json:"update_id"`
	Message      TelegramMessageUpdate `json:"message,omitempty"`
	MyChatMember BotMemberUpdated      `json:"my_chat_member,omitempty"`
}
