package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golangTroshin/gophermat/internal/middleware"
	"github.com/golangTroshin/gophermat/internal/model"
	"github.com/golangTroshin/gophermat/internal/service"
)

type WithdrawHandler struct {
	balanceService  service.BalanceService
	withdrawService service.WithdrawService
	orderService    service.OrderService
}

type WithdrawalResponse struct {
	OrderNumber string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func NewWithdrawHandler(balanceService service.BalanceService, withdrawService service.WithdrawService, orderService service.OrderService) *WithdrawHandler {
	return &WithdrawHandler{balanceService, withdrawService, orderService}
}

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (h *WithdrawHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uint)

	if !ok {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	var req WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if len(req.Order) == 0 || !h.orderService.ValidateOrderNumber(req.Order) {
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		return
	}

	user, err := h.balanceService.GetUserWithBalance(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if user.Balance.Current < req.Sum {
		http.Error(w, "Insufficient funds", http.StatusPaymentRequired)
		return
	}

	user.Balance.Current -= req.Sum
	user.Balance.Withdrawn += req.Sum

	if err := h.balanceService.UpdateUserBalance(user); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	withdrawal := &model.UserWithdrawal{
		UserID:      user.ID,
		OrderNumber: req.Order,
		Sum:         req.Sum,
		ProcessedAt: time.Now().Format(time.RFC3339),
	}
	if err := h.withdrawService.RecordWithdrawal(withdrawal); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WithdrawHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uint)

	if !ok {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	withdrawals, err := h.withdrawService.GetUserWithdrawals(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var response []WithdrawalResponse
	for _, withdrawal := range withdrawals {
		response = append(response, WithdrawalResponse{
			OrderNumber: withdrawal.OrderNumber,
			Sum:         withdrawal.Sum,
			ProcessedAt: withdrawal.ProcessedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
