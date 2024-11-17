package business

import (
	"encoding/json"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	types2 "github.com/adriein/hastypal/internal/hastypal/shared/types"
	"net/http"
)

type CreateBusinessHandler struct {
	service *CreateBusinessService
}

func NewCreateBusinessHandler(
	service *CreateBusinessService,
) *CreateBusinessHandler {
	return &CreateBusinessHandler{
		service: service,
	}
}

func (h *CreateBusinessHandler) Handler(w http.ResponseWriter, r *http.Request) error {
	var request types2.Business

	if decodeErr := json.NewDecoder(r.Body).Decode(&request); decodeErr != nil {
		return types2.ApiError{
			Msg:      decodeErr.Error(),
			Function: "Handler -> json.NewDecoder().Decode()",
			File:     "create-business.go",
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
