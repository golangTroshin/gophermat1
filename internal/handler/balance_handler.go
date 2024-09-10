package handler

import (
	"encoding/json"
	"net/http"

	"github.com/golangTroshin/gophermat/internal/middleware"
	"github.com/golangTroshin/gophermat/internal/service"
)

type BalanceHandler struct {
	userService *service.UserService
}

func NewBalanceHandler(userService *service.UserService) *BalanceHandler {
	return &BalanceHandler{userService}
}

func (h *BalanceHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uint)

	if !ok {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	balance, err := h.userService.GetUserBalance(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"current":   balance.Current,
		"withdrawn": balance.Withdrawn,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	w.WriteHeader(http.StatusOK)
}
