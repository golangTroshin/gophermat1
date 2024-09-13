package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golangTroshin/gophermat/internal/handler"
	"github.com/golangTroshin/gophermat/internal/model"
)

type MockWithdrawService struct {
	RecordWithdrawalFunc   func(withdrawal *model.UserWithdrawal) error
	GetUserWithdrawalsFunc func(userID uint) ([]model.UserWithdrawal, error)
}

func (m *MockWithdrawService) RecordWithdrawal(withdrawal *model.UserWithdrawal) error {
	if m.RecordWithdrawalFunc != nil {
		return m.RecordWithdrawalFunc(withdrawal)
	}
	return nil
}

func (m *MockWithdrawService) GetUserWithdrawals(userID uint) ([]model.UserWithdrawal, error) {
	if m.GetUserWithdrawalsFunc != nil {
		return m.GetUserWithdrawalsFunc(userID)
	}
	return nil, nil
}

func TestWithdrawHandler_Withdraw(t *testing.T) {
	mockBalanceService := &MockBalanceService{}
	mockWithdrawService := &MockWithdrawService{}
	mockOrderService := &MockOrderService{}

	withdrawHandler := handler.NewWithdrawHandler(mockBalanceService, mockWithdrawService, mockOrderService)

	tests := map[string]struct {
		requestBody        string
		userID             uint
		mockOrderValidate  func(order string) bool
		mockGetUserBalance func(userID uint) (*model.User, error)
		mockUpdateBalance  func(user *model.User) error
		mockRecordWithdraw func(withdrawal *model.UserWithdrawal) error
		expectedStatus     int
	}{
		"Valid Withdrawal": {
			requestBody: `{"order": "123456789", "sum": 50}`,
			userID:      1,
			mockOrderValidate: func(order string) bool {
				return true
			},
			mockGetUserBalance: func(userID uint) (*model.User, error) {
				return &model.User{ID: 1, Balance: model.UserBalance{Current: 100.0, Withdrawn: 0.0}}, nil
			},
			mockUpdateBalance: func(user *model.User) error {
				return nil
			},
			mockRecordWithdraw: func(withdrawal *model.UserWithdrawal) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		"Invalid Order Number": {
			requestBody: `{"order": "invalid_order", "sum": 50}`,
			userID:      1,
			mockOrderValidate: func(order string) bool {
				return false
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		"Insufficient Funds": {
			requestBody: `{"order": "123456789", "sum": 150}`,
			userID:      1,
			mockOrderValidate: func(order string) bool {
				return true
			},
			mockGetUserBalance: func(userID uint) (*model.User, error) {
				return &model.User{ID: 1, Balance: model.UserBalance{Current: 100.0, Withdrawn: 0.0}}, nil
			},
			expectedStatus: http.StatusPaymentRequired,
		},
		"Internal Server Error on Balance Retrieval": {
			requestBody: `{"order": "123456789", "sum": 50}`,
			userID:      1,
			mockOrderValidate: func(order string) bool {
				return true
			},
			mockGetUserBalance: func(userID uint) (*model.User, error) {
				return nil, errors.New("internal error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Mock order service validation
			mockOrderService.ValidateOrderNumberFunc = tc.mockOrderValidate

			// Mock balance service functions
			mockBalanceService.GetUserWithBalanceFunc = tc.mockGetUserBalance
			mockBalanceService.UpdateUserBalanceFunc = tc.mockUpdateBalance

			// Mock withdraw service functions
			mockWithdrawService.RecordWithdrawalFunc = tc.mockRecordWithdraw

			// Create a request
			req := httptest.NewRequest("POST", "/withdraw", bytes.NewBufferString(tc.requestBody))
			req = req.WithContext(mockUserIDContext(tc.userID))
			rr := httptest.NewRecorder()

			// Call the handler
			withdrawHandler.Withdraw(rr, req)

			// Validate the status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}
		})
	}
}

func TestWithdrawHandler_GetWithdrawals(t *testing.T) {
	mockBalanceService := &MockBalanceService{}
	mockWithdrawService := &MockWithdrawService{}
	mockOrderService := &MockOrderService{}

	withdrawHandler := handler.NewWithdrawHandler(mockBalanceService, mockWithdrawService, mockOrderService)

	tests := map[string]struct {
		userID               uint
		mockGetUserWithdraws func(userID uint) ([]model.UserWithdrawal, error)
		expectedStatus       int
		expectedResponse     []handler.WithdrawalResponse
	}{
		"Valid Withdrawals": {
			userID: 1,
			mockGetUserWithdraws: func(userID uint) ([]model.UserWithdrawal, error) {
				return []model.UserWithdrawal{
					{
						Order:       "123456789",
						Sum:         50,
						ProcessedAt: time.Now(),
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: []handler.WithdrawalResponse{
				{
					Order:       "123456789",
					Sum:         50,
					ProcessedAt: time.Now().Format(time.RFC3339),
				},
			},
		},
		"No Withdrawals": {
			userID: 1,
			mockGetUserWithdraws: func(userID uint) ([]model.UserWithdrawal, error) {
				return []model.UserWithdrawal{}, nil
			},
			expectedStatus:   http.StatusNoContent,
			expectedResponse: nil,
		},
		"Internal Server Error": {
			userID: 1,
			mockGetUserWithdraws: func(userID uint) ([]model.UserWithdrawal, error) {
				return nil, errors.New("internal error")
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Mock withdraw service functions
			mockWithdrawService.GetUserWithdrawalsFunc = tc.mockGetUserWithdraws

			// Create a request
			req := httptest.NewRequest("GET", "/withdrawals", nil)
			req = req.WithContext(mockUserIDContext(tc.userID))
			rr := httptest.NewRecorder()

			// Call the handler
			withdrawHandler.GetWithdrawals(rr, req)

			// Validate the status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			// Validate the response if applicable
			if tc.expectedStatus == http.StatusOK {
				var response []handler.WithdrawalResponse
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				// Compare expected and actual response
				for i, expected := range tc.expectedResponse {
					if response[i].Order != expected.Order || response[i].Sum != expected.Sum {
						t.Errorf("expected %v, got %v", expected, response[i])
					}
				}
			}
		})
	}
}
