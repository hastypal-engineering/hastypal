package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	"github.com/golang-jwt/jwt/v5"
)

func NewAuthMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)

			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)

			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv(constants.JwtKey)), nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)

			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)

			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx := context.WithValue(r.Context(), constants.ClaimsContextKey, claims)
			r = r.WithContext(ctx)
		}

		handler.ServeHTTP(w, r)
	})
}
