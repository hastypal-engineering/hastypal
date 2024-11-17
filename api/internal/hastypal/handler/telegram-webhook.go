package handler

import (
	"encoding/json"
	"github.com/adriein/hastypal/internal/hastypal/shared/service"
	types2 "github.com/adriein/hastypal/internal/hastypal/shared/types"
	"net/http"
)

type TelegramWebhookHandler struct {
	service *service.TelegramWebhookService
}

func NewTelegramWebhookHandler(
	service *service.TelegramWebhookService,
) *TelegramWebhookHandler {
	return &TelegramWebhookHandler{
		service: service,
	}
}

func (h *TelegramWebhookHandler) Handler(w http.ResponseWriter, r *http.Request) error {
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
