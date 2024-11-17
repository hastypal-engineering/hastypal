package handler

import (
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/service"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"net/http"
)

type GoogleAuthHandler struct {
	googleApi *service.GoogleApi
}

func NewGoogleAuthHandler(
	googleApi *service.GoogleApi,
) *GoogleAuthHandler {
	return &GoogleAuthHandler{
		googleApi: googleApi,
	}
}

func (h *GoogleAuthHandler) Handler(w http.ResponseWriter, _ *http.Request) error {
	googleAuthUrl := h.googleApi.GetAuthCodeUrl()

	response := types.ServerResponse{
		Ok:   true,
		Data: googleAuthUrl,
	}

	if err := helper.Encode[types.ServerResponse](w, http.StatusOK, response); err != nil {
		return err
	}

	return nil
}
