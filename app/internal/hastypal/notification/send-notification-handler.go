package notification

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"net/http"
)

type SendNotificationHandler struct {
	service *SendNotificationService
}

func NewSendNotificationHandler(
	service *SendNotificationService,
) *SendNotificationHandler {
	return &SendNotificationHandler{
		service: service,
	}
}

func (h *SendNotificationHandler) Handler(w http.ResponseWriter, _ *http.Request) error {
	if serviceErr := h.service.Execute(); serviceErr != nil {
		return exception.Wrap("h.service.Execute", "send-notification-handler.go", serviceErr)
	}

	response := types.ServerResponse{Ok: true}

	if err := helper.Encode[types.ServerResponse](w, http.StatusCreated, response); err != nil {
		return exception.Wrap("helper.Encode", "send-notification-handler.go", err)
	}

	return nil
}
