package google

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
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
	parsedUrl, urlParseErr := url.Parse(r.URL.String())

	if urlParseErr != nil {
		return exception.New(urlParseErr.Error()).
			Trace("url.Parse", "auth-google-handler.go").
			WithValues([]string{r.URL.String()})
	}

	businessId := parsedUrl.Query().Get("businessId")

	googleAuthUrl := h.service.Execute(businessId)

	response := types.ServerResponse{
		Ok:   true,
		Data: googleAuthUrl,
	}

	if err := helper.Encode[types.ServerResponse](w, http.StatusOK, response); err != nil {
		return exception.Wrap(
			"helper.Encode",
			"auth-google-handler.go",
			err,
		)
	}

	return nil
}
