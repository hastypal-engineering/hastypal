package handler

import (
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/service"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"net/http"
)

type WelcomeHandler struct {
	bot       *service.TelegramBot
	presenter types.Presenter
}

func NewWelcomeHandler(
	bot *service.TelegramBot,
	presenter types.Presenter,
) *WelcomeHandler {
	return &WelcomeHandler{
		bot:       bot,
		presenter: presenter,
	}
}

func (h *WelcomeHandler) Handler(w http.ResponseWriter, r *http.Request) error {

	response, presenterErr := h.presenter.Format(nil)

	if presenterErr != nil {
		return presenterErr
	}

	if err := helper.Encode[types.ServerResponse](w, http.StatusOK, response); err != nil {
		return err
	}

	return nil
}
