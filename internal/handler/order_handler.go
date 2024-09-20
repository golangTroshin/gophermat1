package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/golangTroshin/gophermat/internal/middleware"
	"github.com/golangTroshin/gophermat/internal/service"
)

type OrderHandler struct {
	orderService service.OrderService
}

type OrderResponse struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uint)
	if !ok {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	orderNumber := string(body)

	if !h.orderService.ValidateOrderNumber(orderNumber) {
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		return
	}

	err = h.orderService.CreateOrder(userID, orderNumber)
	if err != nil {
		if errors.Is(err, service.ErrOrderExistsSameUser) {
			w.WriteHeader(http.StatusOK)
			return
		} else if errors.Is(err, service.ErrOrderExistsDifferentUser) {
			http.Error(w, "Order already uploaded by another user", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uint)
	if !ok {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	orders, err := h.orderService.GetOrdersByUserID(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response := make([]OrderResponse, len(orders))
	for i, order := range orders {
		response[i] = OrderResponse{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.UploadedAt.Format(time.RFC3339),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
