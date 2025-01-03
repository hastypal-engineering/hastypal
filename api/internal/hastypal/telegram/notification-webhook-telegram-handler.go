package telegram

import (
	"encoding/json"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"net/http"
)

type NotificationWebhookTelegramHandler struct {
	service *NotificationWebhookTelegramService
}

func NewNotificationWebhookTelegramHandler(
	service *NotificationWebhookTelegramService,
) *NotificationWebhookTelegramHandler {
	return &NotificationWebhookTelegramHandler{
		service: service,
	}
}

func (h *NotificationWebhookTelegramHandler) Handler(w http.ResponseWriter, r *http.Request) error {
	var update types.TelegramUpdate

	if decodeErr := json.NewDecoder(r.Body).Decode(&update); decodeErr != nil {
		return exception.New(decodeErr.Error()).
			Trace("json.NewDecoder", "notification-webhook-telegram-handler.go")
	}

	if serviceErr := h.service.Execute(update); serviceErr != nil {
		return exception.Wrap("h.service.Execute", "notification-webhook-telegram-handler.go", serviceErr)
	}

	return nil
}
