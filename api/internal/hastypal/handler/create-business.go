package handler

import (
	"encoding/json"
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/service"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"net/http"
)

type CreateBusinessHandler struct {
	service *service.CreateBusinessService
}

func NewCreateBusinessHandler(
	service *service.CreateBusinessService,
) *CreateBusinessHandler {
	return &CreateBusinessHandler{
		service: service,
	}
}

func (h *CreateBusinessHandler) Handler(w http.ResponseWriter, r *http.Request) error {
	var request types.Business

	if decodeErr := json.NewDecoder(r.Body).Decode(&request); decodeErr != nil {
		return types.ApiError{
			Msg:      decodeErr.Error(),
			Function: "Handler -> json.NewDecoder().Decode()",
			File:     "create-business.go",
		}
	}

	if serviceErr := h.service.Execute(request); serviceErr != nil {
		return serviceErr
	}

	response := types.ServerResponse{Ok: true}

	if err := helper.Encode[types.ServerResponse](w, http.StatusCreated, response); err != nil {
		return err
	}

	return nil
}
