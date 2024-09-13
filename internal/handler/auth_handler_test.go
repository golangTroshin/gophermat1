package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golangTroshin/gophermat/internal/handler"
	"github.com/golangTroshin/gophermat/internal/model"
)

type MockAuthService struct {
	RegisterUserFunc     func(login, password string) (*model.User, error)
	AuthenticateUserFunc func(login, password string) (*model.User, error)
}

func (m *MockAuthService) RegisterUser(login, password string) (*model.User, error) {
	if m.RegisterUserFunc != nil {
		return m.RegisterUserFunc(login, password)
	}
	return nil, nil
}

func (m *MockAuthService) AuthenticateUser(login, password string) (*model.User, error) {
	if m.AuthenticateUserFunc != nil {
		return m.AuthenticateUserFunc(login, password)
	}
	return nil, nil
}

func MockSetAuthCookie(userID uint, w http.ResponseWriter) {
	w.Header().Set("Set-Cookie", "auth_token=mocked_token; HttpOnly")
}

func TestAuthHandler_Register(t *testing.T) {
	mockService := &MockAuthService{}

	tests := map[string]struct {
		requestBody      string
		mockRegisterFunc func(login, password string) (*model.User, error)
		expectedStatus   int
		expectedCookie   bool
	}{
		"Valid Registration": {
			requestBody: `{"login":"valid_user", "password":"valid_password"}`,
			mockRegisterFunc: func(login, password string) (*model.User, error) {
				return &model.User{ID: 1}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCookie: true,
		},
		"Login Already Exists": {
			requestBody: `{"login":"existing_user", "password":"password"}`,
			mockRegisterFunc: func(login, password string) (*model.User, error) {
				return nil, errors.New("login already exists")
			},
			expectedStatus: http.StatusConflict,
			expectedCookie: false,
		},
		"Invalid Request Body": {
			requestBody:    `{"login":"missing_password"}`,
			expectedStatus: http.StatusBadRequest,
			expectedCookie: false,
		},
		"Server Error": {
			requestBody: `{"login":"valid_user", "password":"password"}`,
			mockRegisterFunc: func(login, password string) (*model.User, error) {
				return nil, errors.New("some internal error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCookie: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.mockRegisterFunc != nil {
				mockService.RegisterUserFunc = tc.mockRegisterFunc
			}

			req := httptest.NewRequest("POST", "/api/user/register", bytes.NewBufferString(tc.requestBody))
			rr := httptest.NewRecorder()

			authHandler := handler.NewAuthHandler(mockService)

			authHandler.Register(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			cookie := rr.Header().Get("Set-Cookie")
			if tc.expectedCookie && cookie == "" {
				t.Errorf("expected Set-Cookie header, but none was set")
			}
			if !tc.expectedCookie && cookie != "" {
				t.Errorf("did not expect Set-Cookie header, but it was set")
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	mockService := &MockAuthService{}

	tests := map[string]struct {
		requestBody    string
		mockLoginFunc  func(login, password string) (*model.User, error)
		expectedStatus int
		expectedCookie bool
	}{
		"Valid Login": {
			requestBody: `{"login":"valid_user", "password":"valid_password"}`,
			mockLoginFunc: func(login, password string) (*model.User, error) {
				return &model.User{ID: 1}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCookie: true,
		},
		"Invalid Credentials": {
			requestBody: `{"login":"invalid_user", "password":"wrong_password"}`,
			mockLoginFunc: func(login, password string) (*model.User, error) {
				return nil, errors.New("invalid credentials")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedCookie: false,
		},
		"Invalid Request Body": {
			requestBody:    `{"login":"missing_password"}`,
			mockLoginFunc:  nil,
			expectedStatus: http.StatusBadRequest,
			expectedCookie: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.mockLoginFunc != nil {
				mockService.AuthenticateUserFunc = tc.mockLoginFunc
			}

			req := httptest.NewRequest("POST", "/api/user/login", bytes.NewBufferString(tc.requestBody))
			rr := httptest.NewRecorder()

			authHandler := handler.NewAuthHandler(mockService)
			authHandler.Login(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			cookie := rr.Header().Get("Set-Cookie")
			if tc.expectedCookie && cookie == "" {
				t.Errorf("expected Set-Cookie header, but none was set")
			}
			if !tc.expectedCookie && cookie != "" {
				t.Errorf("did not expect Set-Cookie header, but it was set")
			}
		})
	}
}
