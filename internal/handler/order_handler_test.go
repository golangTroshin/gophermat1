package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golangTroshin/gophermat/internal/handler"
	"github.com/golangTroshin/gophermat/internal/middleware"
	"github.com/golangTroshin/gophermat/internal/model"
	"github.com/golangTroshin/gophermat/internal/service"
)

type MockOrderService struct {
	ValidateOrderNumberFunc func(orderNumber string) bool
	CreateOrderFunc         func(userID uint, orderNumber string) error
	GetOrdersByUserIDFunc   func(userID uint) ([]model.Order, error)
}

func (m *MockOrderService) ValidateOrderNumber(orderNumber string) bool {
	if m.ValidateOrderNumberFunc != nil {
		return m.ValidateOrderNumberFunc(orderNumber)
	}
	return false
}

func (m *MockOrderService) CreateOrder(userID uint, orderNumber string) error {
	if m.CreateOrderFunc != nil {
		return m.CreateOrderFunc(userID, orderNumber)
	}
	return nil
}

func (m *MockOrderService) GetOrdersByUserID(userID uint) ([]model.Order, error) {
	if m.GetOrdersByUserIDFunc != nil {
		return m.GetOrdersByUserIDFunc(userID)
	}
	return nil, nil
}

func TestOrderHandler_UploadOrder(t *testing.T) {
	mockService := &MockOrderService{}
	orderHandler := handler.NewOrderHandler(mockService)

	tests := map[string]struct {
		body                string
		userID              uint
		validateOrderNumber func(orderNumber string) bool
		createOrder         func(userID uint, orderNumber string) error
		expectedStatus      int
	}{
		"Valid Order Upload": {
			body:   "123456789",
			userID: 1,
			validateOrderNumber: func(orderNumber string) bool {
				return true
			},
			createOrder: func(userID uint, orderNumber string) error {
				return nil
			},
			expectedStatus: http.StatusAccepted,
		},
		"Invalid Order Number": {
			body:   "invalidorder",
			userID: 1,
			validateOrderNumber: func(orderNumber string) bool {
				return false
			},
			createOrder: func(userID uint, orderNumber string) error {
				return nil
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		"Order Exists Same User": {
			body:   "123456789",
			userID: 1,
			validateOrderNumber: func(orderNumber string) bool {
				return true
			},
			createOrder: func(userID uint, orderNumber string) error {
				return service.ErrOrderExistsSameUser
			},
			expectedStatus: http.StatusOK,
		},
		"Internal Server Error": {
			body:   "123456789",
			userID: 1,
			validateOrderNumber: func(orderNumber string) bool {
				return true
			},
			createOrder: func(userID uint, orderNumber string) error {
				return errors.New("some internal error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockService.ValidateOrderNumberFunc = tc.validateOrderNumber
			mockService.CreateOrderFunc = tc.createOrder

			req := httptest.NewRequest("POST", "/api/user/orders", nil)
			ctx := context.WithValue(req.Context(), middleware.UserIDContextKey, tc.userID)
			req = req.WithContext(ctx)
			req.Body = io.NopCloser(bytes.NewBufferString(tc.body))
			rr := httptest.NewRecorder()

			orderHandler.UploadOrder(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}
		})
	}
}

func TestOrderHandler_GetOrders(t *testing.T) {
	mockService := &MockOrderService{}
	orderHandler := handler.NewOrderHandler(mockService)

	tests := map[string]struct {
		userID           uint
		getOrdersByUser  func(userID uint) ([]model.Order, error)
		expectedStatus   int
		expectedResponse []handler.OrderResponse
	}{
		"Valid Orders Retrieval": {
			userID: 1,
			getOrdersByUser: func(userID uint) ([]model.Order, error) {
				return []model.Order{
					{
						Number:     "123456789",
						Status:     "PROCESSED",
						Accrual:    100.0,
						UploadedAt: time.Now(),
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: []handler.OrderResponse{
				{
					Number:     "123456789",
					Status:     "PROCESSED",
					Accrual:    100.0,
					UploadedAt: time.Now().Format(time.RFC3339),
				},
			},
		},
		"No Orders Found": {
			userID: 1,
			getOrdersByUser: func(userID uint) ([]model.Order, error) {
				return []model.Order{}, nil
			},
			expectedStatus:   http.StatusNoContent,
			expectedResponse: []handler.OrderResponse{},
		},
		"Internal Server Error": {
			userID: 1,
			getOrdersByUser: func(userID uint) ([]model.Order, error) {
				return nil, errors.New("some internal error")
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockService.GetOrdersByUserIDFunc = tc.getOrdersByUser

			req := httptest.NewRequest("GET", "/api/user/orders", nil)
			ctx := context.WithValue(req.Context(), middleware.UserIDContextKey, tc.userID)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			orderHandler.GetOrders(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			if tc.expectedStatus == http.StatusOK {
				var orders []handler.OrderResponse
				if err := json.NewDecoder(rr.Body).Decode(&orders); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if len(orders) != len(tc.expectedResponse) {
					t.Errorf("expected %d orders, got %d", len(tc.expectedResponse), len(orders))
				}

				for i, expectedOrder := range tc.expectedResponse {
					if orders[i].Number != expectedOrder.Number || orders[i].Status != expectedOrder.Status {
						t.Errorf("expected order %v, got %v", expectedOrder, orders[i])
					}
				}
			}
		})
	}
}
