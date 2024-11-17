package google

import (
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"net/http"
)

type AuthGoogleHandler struct {
	service *AuthGoogleService
}

func NewGoogleAuthHandler(
	service *AuthGoogleService,
) *AuthGoogleHandler {
	return &AuthGoogleHandler{
		service: service,
	}
}

func (h *AuthGoogleHandler) Handler(w http.ResponseWriter, _ *http.Request) error {
	googleAuthUrl := h.service.Execute()

	response := types.ServerResponse{
		Ok:   true,
		Data: googleAuthUrl,
	}

	if err := helper.Encode[types.ServerResponse](w, http.StatusOK, response); err != nil {
		return err
	}

	return nil
}
