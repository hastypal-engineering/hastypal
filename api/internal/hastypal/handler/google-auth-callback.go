package handler

import (
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/service"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"net/http"
)

type GoogleAuthCallbackHandler struct {
	service service.GoogleAuthCallbackService
}

func NewGoogleAuthCallbackHandler(
	service service.GoogleAuthCallbackService,
) *GoogleAuthCallbackHandler {
	return &GoogleAuthCallbackHandler{
		service: service,
	}
}

func (h *GoogleAuthCallbackHandler) Handler(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.Execute(r.RequestURI); err != nil {
		return err
	}

	response := types.ServerResponse{
		Ok:   true,
		Data: nil,
	}

	if err := helper.Encode[types.ServerResponse](w, http.StatusOK, response); err != nil {
		return err
	}

	return nil
}
