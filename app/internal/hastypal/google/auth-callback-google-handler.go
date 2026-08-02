package google

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"net/http"
)

type AuthCallbackGoogleHandler struct {
	service *AuthCallbackGoogleService
}

func NewGoogleAuthCallbackHandler(
	service *AuthCallbackGoogleService,
) *AuthCallbackGoogleHandler {
	return &AuthCallbackGoogleHandler{
		service: service,
	}
}

func (h *AuthCallbackGoogleHandler) Handler(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.Execute(r.RequestURI); err != nil {
		return exception.Wrap(
			"h.service.Execute",
			"auth-callback-google-handler.go",
			err,
		)
	}

	response := types.ServerResponse{
		Ok:   true,
		Data: nil,
	}

	if err := helper.Encode[types.ServerResponse](w, http.StatusOK, response); err != nil {
		return exception.Wrap(
			"helper.Encode",
			"auth-callback-google-handler.go",
			err,
		)
	}

	return nil
}
