package business

import (
	"encoding/json"
	"net/http"

	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
)

type LoginBusinessHandler struct {
	service *LoginBusinessService
}

func NewLoginBusinessHandler(
	service *LoginBusinessService,
) *LoginBusinessHandler {
	return &LoginBusinessHandler{
		service: service,
	}
}

func (h *LoginBusinessHandler) Handler(w http.ResponseWriter, r *http.Request) error {
	var request LoginBusiness

	if decodeErr := json.NewDecoder(r.Body).Decode(&request); decodeErr != nil {
		return types.ApiError{
			Msg:      decodeErr.Error(),
			Function: "Handler -> json.NewDecoder().Decode()",
			File:     "login-business-handler.go",
		}
	}

	if serviceErr := h.service.Execute(request); serviceErr != nil {
		return serviceErr
	}

	response := types.ServerResponse{Ok: true}

	if err := helper.Encode[types.ServerResponse](w, http.StatusAccepted, response); err != nil {
		return err
	}

	// Add bearer-token jwt authentication

	return nil
}
