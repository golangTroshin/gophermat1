package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golangTroshin/gophermat/internal/handler"
	"github.com/golangTroshin/gophermat/internal/mock_service"
	"github.com/golangTroshin/gophermat/internal/model"
	"github.com/golangTroshin/gophermat/internal/service"
)

func TestAuthHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthService(ctrl)
	authHandler := handler.NewAuthHandler(mockAuthService)

	tests := map[string]struct {
		requestBody     string
		mockRegister    func(m *mock_service.MockAuthService)
		expectedStatus  int
		expectedMessage string
	}{
		"Valid_Registration": {
			requestBody: `{"login":"valid_user", "password":"valid_password"}`,
			mockRegister: func(m *mock_service.MockAuthService) {
				m.EXPECT().RegisterUser("valid_user", "valid_password").Return(&model.User{ID: 1}, nil)
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "",
		},
		"Missing_Login_Or_Password": {
			requestBody:     `{"login":"", "password":"password"}`,
			mockRegister:    nil,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "login and password are required",
		},
		"User_Already_Exists": {
			requestBody: `{"login":"existing_user", "password":"password"}`,
			mockRegister: func(m *mock_service.MockAuthService) {
				m.EXPECT().RegisterUser("existing_user", "password").Return(nil, service.ErrUserExists)
			},
			expectedStatus:  http.StatusConflict,
			expectedMessage: "user already exists",
		},
		"Internal_Server_Error": {
			requestBody: `{"login":"valid_user", "password":"valid_password"}`,
			mockRegister: func(m *mock_service.MockAuthService) {
				m.EXPECT().RegisterUser("valid_user", "valid_password").Return(nil, errors.New("internal error"))
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "server error",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.mockRegister != nil {
				tc.mockRegister(mockAuthService)
			}

			req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(tc.requestBody))
			rr := httptest.NewRecorder()

			authHandler.Register(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			if tc.expectedMessage != "" {
				if rr.Body.String() != tc.expectedMessage+"\n" {
					t.Errorf("expected message %q, got %q", tc.expectedMessage, rr.Body.String())
				}
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthService(ctrl)
	authHandler := handler.NewAuthHandler(mockAuthService)

	tests := map[string]struct {
		requestBody     string
		mockLogin       func(m *mock_service.MockAuthService)
		expectedStatus  int
		expectedMessage string
	}{
		"Valid_Login": {
			requestBody: `{"login":"valid_user", "password":"valid_password"}`,
			mockLogin: func(m *mock_service.MockAuthService) {
				m.EXPECT().AuthenticateUser("valid_user", "valid_password").Return(&model.User{ID: 1}, nil)
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "",
		},
		"Missing_Login_Or_Password": {
			requestBody:     `{"login":"", "password":"password"}`,
			mockLogin:       nil,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "login and password are required",
		},
		"Invalid_Credentials": {
			requestBody: `{"login":"invalid_user", "password":"wrong_password"}`,
			mockLogin: func(m *mock_service.MockAuthService) {
				m.EXPECT().AuthenticateUser("invalid_user", "wrong_password").Return(nil, service.ErrInvalidCreds)
			},
			expectedStatus:  http.StatusUnauthorized,
			expectedMessage: "invalid credentials",
		},
		"Internal_Server_Error": {
			requestBody: `{"login":"valid_user", "password":"valid_password"}`,
			mockLogin: func(m *mock_service.MockAuthService) {
				m.EXPECT().AuthenticateUser("valid_user", "valid_password").Return(nil, errors.New("internal error"))
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "server error",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.mockLogin != nil {
				tc.mockLogin(mockAuthService)
			}

			req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(tc.requestBody))
			rr := httptest.NewRecorder()

			authHandler.Login(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			if tc.expectedMessage != "" {
				if rr.Body.String() != tc.expectedMessage+"\n" {
					t.Errorf("expected message %q, got %q", tc.expectedMessage, rr.Body.String())
				}
			}
		})
	}
}
