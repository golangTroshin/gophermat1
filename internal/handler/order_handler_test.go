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

	"github.com/golang/mock/gomock"
	"github.com/golangTroshin/gophermat/internal/handler"
	"github.com/golangTroshin/gophermat/internal/middleware"
	"github.com/golangTroshin/gophermat/internal/mock_service"
	"github.com/golangTroshin/gophermat/internal/model"
	"github.com/golangTroshin/gophermat/internal/service"
)

func TestOrderHandler_UploadOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderService := mock_service.NewMockOrderService(ctrl)
	orderHandler := handler.NewOrderHandler(mockOrderService)

	tests := map[string]struct {
		body              string
		userID            uint
		mockValidateOrder func(m *mock_service.MockOrderService)
		mockCreateOrder   func(m *mock_service.MockOrderService)
		expectedStatus    int
	}{
		"Valid_Order_Upload": {
			body:   "123456789",
			userID: 1,
			mockValidateOrder: func(m *mock_service.MockOrderService) {
				m.EXPECT().ValidateOrderNumber("123456789").Return(true)
			},
			mockCreateOrder: func(m *mock_service.MockOrderService) {
				m.EXPECT().CreateOrder(uint(1), "123456789").Return(nil)
			},
			expectedStatus: http.StatusAccepted,
		},
		"Invalid_Order_Number": {
			body:   "invalidorder",
			userID: 1,
			mockValidateOrder: func(m *mock_service.MockOrderService) {
				m.EXPECT().ValidateOrderNumber("invalidorder").Return(false)
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		"Empty_Body": {
			body:           "",
			userID:         1,
			expectedStatus: http.StatusBadRequest,
		},
		"Order_Exists_Same_User": {
			body:   "123456789",
			userID: 1,
			mockValidateOrder: func(m *mock_service.MockOrderService) {
				m.EXPECT().ValidateOrderNumber("123456789").Return(true)
			},
			mockCreateOrder: func(m *mock_service.MockOrderService) {
				m.EXPECT().CreateOrder(uint(1), "123456789").Return(service.ErrOrderExistsSameUser)
			},
			expectedStatus: http.StatusOK,
		},
		"Order_Exists_Different_User": {
			body:   "123456789",
			userID: 1,
			mockValidateOrder: func(m *mock_service.MockOrderService) {
				m.EXPECT().ValidateOrderNumber("123456789").Return(true)
			},
			mockCreateOrder: func(m *mock_service.MockOrderService) {
				m.EXPECT().CreateOrder(uint(1), "123456789").Return(service.ErrOrderExistsDifferentUser)
			},
			expectedStatus: http.StatusConflict,
		},
		"Internal_Server_Error": {
			body:   "123456789",
			userID: 1,
			mockValidateOrder: func(m *mock_service.MockOrderService) {
				m.EXPECT().ValidateOrderNumber("123456789").Return(true)
			},
			mockCreateOrder: func(m *mock_service.MockOrderService) {
				m.EXPECT().CreateOrder(uint(1), "123456789").Return(errors.New("some internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.mockValidateOrder != nil {
				tc.mockValidateOrder(mockOrderService)
			}

			if tc.mockCreateOrder != nil {
				tc.mockCreateOrder(mockOrderService)
			}

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderService := mock_service.NewMockOrderService(ctrl)
	orderHandler := handler.NewOrderHandler(mockOrderService)

	tests := map[string]struct {
		userID              uint
		mockGetOrdersByUser func(m *mock_service.MockOrderService)
		expectedStatus      int
		expectedResponse    []handler.OrderResponse
	}{
		"Valid_Orders_Retrieval": {
			userID: 1,
			mockGetOrdersByUser: func(m *mock_service.MockOrderService) {
				m.EXPECT().GetOrdersByUserID(uint(1)).Return([]model.Order{
					{
						Number:     "123456789",
						Status:     "PROCESSED",
						Accrual:    100.0,
						UploadedAt: time.Now(),
					},
				}, nil)
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
		"No_Orders_Found": {
			userID: 1,
			mockGetOrdersByUser: func(m *mock_service.MockOrderService) {
				m.EXPECT().GetOrdersByUserID(uint(1)).Return([]model.Order{}, nil)
			},
			expectedStatus:   http.StatusNoContent,
			expectedResponse: nil,
		},
		"Internal_Server_Error": {
			userID: 1,
			mockGetOrdersByUser: func(m *mock_service.MockOrderService) {
				m.EXPECT().GetOrdersByUserID(uint(1)).Return(nil, errors.New("some internal error"))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.mockGetOrdersByUser != nil {
				tc.mockGetOrdersByUser(mockOrderService)
			}

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
