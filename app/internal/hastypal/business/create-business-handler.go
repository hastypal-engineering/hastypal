package business

import (
	"encoding/json"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
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
	var request types.Business

	if decodeErr := json.NewDecoder(r.Body).Decode(&request); decodeErr != nil {
		return exception.New(decodeErr.Error()).Trace("json.NewDecoder", "create-business-handler.go")
	}

	if serviceErr := h.service.Execute(request); serviceErr != nil {
		return exception.Wrap("h.service.Execute", "create-business-handler.go", serviceErr)
	}

	response := types.ServerResponse{Ok: true}

	if err := helper.Encode[types.ServerResponse](w, http.StatusCreated, response); err != nil {
		return exception.Wrap("helper.Encode", "create-business-handler.go", err)
	}

	return nil
}
