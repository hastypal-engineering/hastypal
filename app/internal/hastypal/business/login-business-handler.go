package business

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/adriein/hastypal/internal/hastypal/shared/exception"

	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"github.com/golang-jwt/jwt/v5"
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
		return exception.New(decodeErr.Error()).Trace("json.NewDecoder", "login-business-handler.go")
	}

	business, serviceErr := h.service.Execute(request)

	if serviceErr != nil {
		return exception.Wrap(
			"h.service.Execute",
			"login-business-handler.go",
			serviceErr,
		)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    business.Id,
		"email": business.Email,
		"name":  business.Name,
	})

	signedJwt, jwtErr := token.SignedString([]byte(os.Getenv(constants.JwtKey)))

	if jwtErr != nil {
		return exception.Wrap(
			"token.SignedString",
			"login-business-handler.go",
			jwtErr,
		)
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", signedJwt))

	response := types.ServerResponse{Ok: true}

	if err := helper.Encode[types.ServerResponse](w, http.StatusAccepted, response); err != nil {
		return exception.Wrap(
			"helper.Encode",
			"login-business-handler.go",
			err,
		)
	}

	return nil
}
