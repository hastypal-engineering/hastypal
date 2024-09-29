package types

type WebhookService interface {
	Execute(update TelegramUpdate) error
}
