package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/golangTroshin/gophermat/internal/handler"
	"github.com/golangTroshin/gophermat/internal/mock_service"
	"github.com/golangTroshin/gophermat/internal/model"
	"github.com/google/go-cmp/cmp"
)

func TestWithdrawHandler_Withdraw(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBalanceService := mock_service.NewMockBalanceService(ctrl)
	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	mockOrderService := mock_service.NewMockOrderService(ctrl)

	withdrawHandler := handler.NewWithdrawHandler(mockBalanceService, mockWithdrawService, mockOrderService)

	tests := map[string]struct {
		requestBody        string
		userID             uint
		mockOrderValidate  func(*mock_service.MockOrderService)
		mockGetUserBalance func(*mock_service.MockBalanceService)
		mockUpdateBalance  func(*mock_service.MockBalanceService)
		mockRecordWithdraw func(*mock_service.MockWithdrawService)
		expectedStatus     int
	}{
		"Valid_Withdrawal": {
			requestBody: `{"order": "123456789", "sum": 50}`,
			userID:      1,
			mockOrderValidate: func(m *mock_service.MockOrderService) {
				m.EXPECT().ValidateOrderNumber("123456789").Return(true)
			},
			mockGetUserBalance: func(m *mock_service.MockBalanceService) {
				m.EXPECT().GetUserWithBalance(uint(1)).Return(&model.User{
					ID: 1, Balance: model.UserBalance{Current: 100.0, Withdrawn: 0.0},
				}, nil)
			},
			mockUpdateBalance: func(m *mock_service.MockBalanceService) {
				m.EXPECT().UpdateUserBalance(gomock.Any()).Return(nil)
			},
			mockRecordWithdraw: func(m *mock_service.MockWithdrawService) {
				m.EXPECT().RecordWithdrawal(gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		"Empty_Body": {
			requestBody:       "",
			userID:            1,
			mockOrderValidate: nil,
			expectedStatus:    http.StatusBadRequest,
		},
		"Invalid_Order_Number": {
			requestBody: `{"order": "invalid_order", "sum": 50}`,
			userID:      1,
			mockOrderValidate: func(m *mock_service.MockOrderService) {
				m.EXPECT().ValidateOrderNumber("invalid_order").Return(false)
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		"Insufficient_Funds": {
			requestBody: `{"order": "123456789", "sum": 150}`,
			userID:      1,
			mockOrderValidate: func(m *mock_service.MockOrderService) {
				m.EXPECT().ValidateOrderNumber("123456789").Return(true)
			},
			mockGetUserBalance: func(m *mock_service.MockBalanceService) {
				m.EXPECT().GetUserWithBalance(uint(1)).Return(&model.User{
					ID: 1, Balance: model.UserBalance{Current: 100.0, Withdrawn: 0.0},
				}, nil)
			},
			expectedStatus: http.StatusPaymentRequired,
		},
		"Internal_Server_Error_on_Balance_Retrieval": {
			requestBody: `{"order": "123456789", "sum": 50}`,
			userID:      1,
			mockOrderValidate: func(m *mock_service.MockOrderService) {
				m.EXPECT().ValidateOrderNumber("123456789").Return(true)
			},
			mockGetUserBalance: func(m *mock_service.MockBalanceService) {
				m.EXPECT().GetUserWithBalance(uint(1)).Return(nil, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.mockOrderValidate != nil {
				tc.mockOrderValidate(mockOrderService)
			}

			if tc.mockGetUserBalance != nil {
				tc.mockGetUserBalance(mockBalanceService)
			}

			if tc.mockUpdateBalance != nil {
				tc.mockUpdateBalance(mockBalanceService)
			}

			if tc.mockRecordWithdraw != nil {
				tc.mockRecordWithdraw(mockWithdrawService)
			}

			req := httptest.NewRequest("POST", "/withdraw", bytes.NewBufferString(tc.requestBody))
			req = req.WithContext(mockUserIDContext(tc.userID))
			rr := httptest.NewRecorder()

			withdrawHandler.Withdraw(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}
		})
	}
}

func TestWithdrawHandler_GetWithdrawals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawService := mock_service.NewMockWithdrawService(ctrl)
	withdrawHandler := handler.NewWithdrawHandler(nil, mockWithdrawService, nil)

	tests := map[string]struct {
		userID               uint
		mockGetUserWithdraws func(*mock_service.MockWithdrawService)
		expectedStatus       int
		expectedResponse     []handler.WithdrawalResponse
	}{
		"Valid_Withdrawals": {
			userID: 1,
			mockGetUserWithdraws: func(m *mock_service.MockWithdrawService) {
				m.EXPECT().GetUserWithdrawals(uint(1)).Return([]model.UserWithdrawal{
					{
						OrderNumber: "123456789",
						Sum:         50,
						ProcessedAt: time.Now().Format(time.RFC3339),
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: []handler.WithdrawalResponse{
				{
					OrderNumber: "123456789",
					Sum:         50,
					ProcessedAt: time.Now().Format(time.RFC3339),
				},
			},
		},
		"No_Withdrawals": {
			userID: 1,
			mockGetUserWithdraws: func(m *mock_service.MockWithdrawService) {
				m.EXPECT().GetUserWithdrawals(uint(1)).Return([]model.UserWithdrawal{}, nil)
			},
			expectedStatus:   http.StatusNoContent,
			expectedResponse: nil,
		},
		"Internal_Server_Error": {
			userID: 1,
			mockGetUserWithdraws: func(m *mock_service.MockWithdrawService) {
				m.EXPECT().GetUserWithdrawals(uint(1)).Return(nil, errors.New("internal error"))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.mockGetUserWithdraws != nil {
				tc.mockGetUserWithdraws(mockWithdrawService)
			}

			req := httptest.NewRequest("GET", "/withdrawals", nil)
			req = req.WithContext(mockUserIDContext(tc.userID))
			rr := httptest.NewRecorder()

			withdrawHandler.GetWithdrawals(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			if tc.expectedStatus == http.StatusOK {
				var response []handler.WithdrawalResponse
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if diff := cmp.Diff(tc.expectedResponse, response); diff != "" {
					t.Errorf("Unexpected response (-expected +got):\n%s", diff)
				}
			}
		})
	}
}
