package google

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"net/http"
	"net/url"
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

func (h *AuthGoogleHandler) Handler(w http.ResponseWriter, r *http.Request) error {
	url, urlParseErr := url.Parse(r.URL.String())

	if urlParseErr != nil {
		return types.ApiError{
			Msg:      urlParseErr.Error(),
			Function: "Handler -> url.Parse()",
			File:     "auth-google-handler.go",
			Values:   []string{r.URL.String()},
		}
	}

	businessId := url.Query().Get("businessId")

	googleAuthUrl := h.service.Execute(businessId)

	response := types.ServerResponse{
		Ok:   true,
		Data: googleAuthUrl,
	}

	if err := helper.Encode[types.ServerResponse](w, http.StatusOK, response); err != nil {
		return err
	}

	return nil
}
