package handler

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/service"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"net/http"
)

type GoogleAuthCallbackHandler struct {
	googleApi *service.GoogleApi
}

func NewGoogleAuthCallbackHandler(
	googleApi *service.GoogleApi,
) *GoogleAuthCallbackHandler {
	return &GoogleAuthCallbackHandler{
		googleApi: googleApi,
	}
}

func (h *GoogleAuthCallbackHandler) Handler(w http.ResponseWriter, r *http.Request) error {
	fmt.Println(r)

	response := types.ServerResponse{
		Ok:   true,
		Data: nil,
	}

	if err := helper.Encode[types.ServerResponse](w, http.StatusOK, response); err != nil {
		return err
	}

	return nil
}
