package constants

//Env var keys

const (
	DatabaseUser             = "DATABASE_USER"
	DatabasePassword         = "DATABASE_PASSWORD"
	DatabaseName             = "DATABASE_NAME"
	ServerPort               = "SERVER_PORT"
	Env                      = "ENV"
	Production               = "PRODUCTION"
	WhatsappBusinessApiToken = "WHATSAPP_BUSINESS_API_TOKEN"
	TelegramApiToken         = "TELEGRAM_API_TOKEN"
	TelegramApiBotUrl        = "TELEGRAM_BOT_API_URL"
)

// Criteria

const (
	Equal              = "="
	GreaterThanOrEqual = ">="
	LessThanOrEqual    = "<="
)

// Errors

const (
	ServerGenericError = "SERVER_ERROR"
)

// Domain constants

const (
	TelegramMessageField string = "Message"
)

// Telegram commands

const (
	StartCommand string = "/start"
)
