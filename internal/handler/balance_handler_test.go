package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golangTroshin/gophermat/internal/handler"
	"github.com/golangTroshin/gophermat/internal/middleware"
	"github.com/golangTroshin/gophermat/internal/model"
)

type MockBalanceService struct {
	GetUserBalanceFunc     func(userID uint) (*model.UserBalance, error)
	GetUserWithBalanceFunc func(userID uint) (*model.User, error)
	UpdateUserBalanceFunc  func(user *model.User) error
}

func (m *MockBalanceService) GetUserBalance(userID uint) (*model.UserBalance, error) {
	if m.GetUserBalanceFunc != nil {
		return m.GetUserBalanceFunc(userID)
	}
	return nil, nil
}

func (m *MockBalanceService) GetUserWithBalance(userID uint) (*model.User, error) {
	if m.GetUserWithBalanceFunc != nil {
		return m.GetUserWithBalanceFunc(userID)
	}
	return nil, nil
}

func (m *MockBalanceService) UpdateUserBalance(user *model.User) error {
	if m.UpdateUserBalanceFunc != nil {
		return m.UpdateUserBalanceFunc(user)
	}
	return nil
}

func mockUserIDContext(userID uint) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, middleware.UserIDContextKey, userID)
}

func TestBalanceHandler_GetUserBalance(t *testing.T) {
	mockService := &MockBalanceService{}
	balanceHandler := handler.NewBalanceHandler(mockService)

	tests := map[string]struct {
		userID           uint
		mockBalanceFunc  func(userID uint) (*model.UserBalance, error)
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{
		"Valid Balance": {
			userID: 1,
			mockBalanceFunc: func(userID uint) (*model.UserBalance, error) {
				return &model.UserBalance{
					Current:   100.0,
					Withdrawn: 50.0,
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"current":   100.0,
				"withdrawn": 50.0,
			},
		},
		"Internal Server Error": {
			userID: 1,
			mockBalanceFunc: func(userID uint) (*model.UserBalance, error) {
				return nil, errors.New("some internal error")
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockService.GetUserBalanceFunc = tc.mockBalanceFunc

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
