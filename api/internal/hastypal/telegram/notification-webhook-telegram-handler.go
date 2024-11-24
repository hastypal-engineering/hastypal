package telegram

import (
	"encoding/json"
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
		return types.ApiError{
			Msg:      decodeErr.Error(),
			Function: "Handler -> json.NewDecoder().Decode()",
			File:     "handler/telegram-webhook.go",
		}
	}

	if serviceErr := h.service.Execute(update); serviceErr != nil {
		return types.WrapErr("h.service.Execute", "notification-webhook-telegram-handler.go", serviceErr)
	}

	return nil
}
