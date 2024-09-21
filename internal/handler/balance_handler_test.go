package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golangTroshin/gophermat/internal/handler"
	"github.com/golangTroshin/gophermat/internal/middleware"
	"github.com/golangTroshin/gophermat/internal/mock_service"
	"github.com/golangTroshin/gophermat/internal/model"
)

func mockUserIDContext(userID uint) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, middleware.UserIDContextKey, userID)
}

func TestBalanceHandler_GetUserBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockBalanceService(ctrl)
	balanceHandler := handler.NewBalanceHandler(mockService)

	tests := map[string]struct {
		userID           uint
		mockBalanceFunc  func(*mock_service.MockBalanceService)
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{
		"Valid_Balance": {
			userID: 1,
			mockBalanceFunc: func(m *mock_service.MockBalanceService) {
				m.EXPECT().GetUserBalance(uint(1)).Return(&model.UserBalance{
					Current:   100.0,
					Withdrawn: 50.0,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"current":   100.0,
				"withdrawn": 50.0,
			},
		},
		"Internal_Server_Error": {
			userID: 1,
			mockBalanceFunc: func(m *mock_service.MockBalanceService) {
				m.EXPECT().GetUserBalance(uint(1)).Return(nil, errors.New("some internal error"))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.mockBalanceFunc != nil {
				tc.mockBalanceFunc(mockService)
			}

			req := httptest.NewRequest("GET", "/balance", nil)
			req = req.WithContext(mockUserIDContext(tc.userID))
			rr := httptest.NewRecorder()

			balanceHandler.GetUserBalance(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			if tc.expectedResponse != nil {
				var response map[string]interface{}
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				for key, expectedValue := range tc.expectedResponse {
					if response[key] != expectedValue {
						t.Errorf("expected %s to be %v, got %v", key, expectedValue, response[key])
					}
				}
			}
		})
	}
}
