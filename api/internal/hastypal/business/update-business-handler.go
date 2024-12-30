package business

import (
	"encoding/json"
	"net/http"

	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
)

type UpdateBusinessHandler struct {
	service *UpdateBusinessService
}

func NewUpdateBusinessHandler(
	service *UpdateBusinessService,
) *UpdateBusinessHandler {
	return &UpdateBusinessHandler{
		service: service,
	}
}

func (h *UpdateBusinessHandler) Handler(w http.ResponseWriter, r *http.Request) error {
	var request types.Business

	if decodeErr := json.NewDecoder(r.Body).Decode(&request); decodeErr != nil {
		return exception.New(decodeErr.Error()).Trace("json.NewDecoder", "update-business-handler.go")
	}

	if serviceErr := h.service.Execute(request); serviceErr != nil {
		return exception.Wrap("h.service.Execute", "update-business-handler.go", serviceErr)
	}

	response := types.ServerResponse{Ok: true}

	if err := helper.Encode[types.ServerResponse](w, http.StatusOK, response); err != nil {
		return exception.Wrap("helper.Encode", "update-business-hanlder.go", err)
	}

	return nil
}
