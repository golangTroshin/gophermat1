package handler

import (
	"encoding/json"
	"net/http"

	"github.com/golangTroshin/gophermat/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "login and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.authService.RegisterUser(req.Login, req.Password)
	if err != nil {
		if err.Error() == "login already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	service.SetAuthCookie(user.ID, w)

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "login and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.authService.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	service.SetAuthCookie(user.ID, w)

	w.WriteHeader(http.StatusOK)
}
