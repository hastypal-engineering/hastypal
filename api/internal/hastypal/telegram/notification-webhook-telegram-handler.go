package telegram

import (
	"encoding/json"
	"github.com/adriein/hastypal/internal/hastypal/shared/service"
	types2 "github.com/adriein/hastypal/internal/hastypal/shared/types"
	"net/http"
)

type NotificationWebhookTelegramHandler struct {
	service *service.TelegramWebhookService
}

func NewNotificationWebhookTelegramHandler(
	service *service.TelegramWebhookService,
) *NotificationWebhookTelegramHandler {
	return &NotificationWebhookTelegramHandler{
		service: service,
	}
}

func (h *NotificationWebhookTelegramHandler) Handler(w http.ResponseWriter, r *http.Request) error {
	var update types2.TelegramUpdate

	if decodeErr := json.NewDecoder(r.Body).Decode(&update); decodeErr != nil {
		return types2.ApiError{
			Msg:      decodeErr.Error(),
			Function: "Handler -> json.NewDecoder().Decode()",
			File:     "handler/telegram-webhook.go",
		}
	}

	if serviceErr := h.service.Execute(update); serviceErr != nil {
		return serviceErr
	}

	return nil
}
