package handler

import (
	"encoding/json"
	"github.com/adriein/hastypal/internal/hastypal/service"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"net/http"
)

type TelegramWebhookHandler struct {
	service *service.CreateBusinessService
}

func NewTelegramWebhookHandler(
	service *service.CreateBusinessService,
) *TelegramWebhookHandler {
	return &TelegramWebhookHandler{
		service: service,
	}
}

func (h *TelegramWebhookHandler) Handler(w http.ResponseWriter, r *http.Request) error {
	var update types.TelegramUpdate

	if decodeErr := json.NewDecoder(r.Body).Decode(&update); decodeErr != nil {
		return types.ApiError{
			Msg:      decodeErr.Error(),
			Function: "Handler -> json.NewDecoder().Decode()",
			File:     "telegram-webhook.go",
		}
	}

	return nil
}
