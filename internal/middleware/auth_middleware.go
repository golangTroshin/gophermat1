package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golangTroshin/gophermat/internal/service"
)

var (
	ErrNoAuthToken  = errors.New("no auth token provided")
	ErrInvalidToken = errors.New("invalid auth token")
)

type ContextKey string

const UserIDContextKey = ContextKey("userID")

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie.Value == "" {
			http.Error(w, ErrNoAuthToken.Error(), http.StatusUnauthorized)
			return
		}

		token, err := validateToken(cookie.Value, service.SecretKey)
		if err != nil {
			http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims["UserID"].(float64)
		if !ok {
			http.Error(w, "userID not found in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, uint(userIDFloat))

		h.ServeHTTP(w, r.WithContext(ctx))
	})

}

func validateToken(tokenString, secretKey string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
}
