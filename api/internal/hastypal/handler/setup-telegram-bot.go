package handler

import (
	"encoding/json"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/service"
	types2 "github.com/adriein/hastypal/internal/hastypal/shared/types"
	"net/http"
)

type SetupTelegramBotHandler struct {
	service *service.SetupTelegramBotService
}

func NewSetupTelegramBotHandler(
	service *service.SetupTelegramBotService,
) *SetupTelegramBotHandler {
	return &SetupTelegramBotHandler{
		service: service,
	}
}

func (h *SetupTelegramBotHandler) Handler(w http.ResponseWriter, r *http.Request) error {
	var request types2.AdminTelegramBotSetup

	if decodeErr := json.NewDecoder(r.Body).Decode(&request); decodeErr != nil {
		return types2.ApiError{
			Msg:      decodeErr.Error(),
			Function: "Handler -> json.NewDecoder().Decode()",
			File:     "handler/setup-telegram.bot.go",
			Values:   []string{},
		}
	}

	if serviceErr := h.service.Execute(request); serviceErr != nil {
		return serviceErr
	}

	response := types2.ServerResponse{Ok: true}

	if err := helper.Encode[types2.ServerResponse](w, http.StatusCreated, response); err != nil {
		return err
	}

	return nil
}
