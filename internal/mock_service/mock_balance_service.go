// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/service/balance_service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/golangTroshin/gophermat/internal/model"
)

// MockBalanceService is a mock of BalanceService interface.
type MockBalanceService struct {
	ctrl     *gomock.Controller
	recorder *MockBalanceServiceMockRecorder
}

// MockBalanceServiceMockRecorder is the mock recorder for MockBalanceService.
type MockBalanceServiceMockRecorder struct {
	mock *MockBalanceService
}

// NewMockBalanceService creates a new mock instance.
func NewMockBalanceService(ctrl *gomock.Controller) *MockBalanceService {
	mock := &MockBalanceService{ctrl: ctrl}
	mock.recorder = &MockBalanceServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBalanceService) EXPECT() *MockBalanceServiceMockRecorder {
	return m.recorder
}

// GetUserBalance mocks base method.
func (m *MockBalanceService) GetUserBalance(userID uint) (*model.UserBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBalance", userID)
	ret0, _ := ret[0].(*model.UserBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBalance indicates an expected call of GetUserBalance.
func (mr *MockBalanceServiceMockRecorder) GetUserBalance(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBalance", reflect.TypeOf((*MockBalanceService)(nil).GetUserBalance), userID)
}

// GetUserWithBalance mocks base method.
func (m *MockBalanceService) GetUserWithBalance(userID uint) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserWithBalance", userID)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserWithBalance indicates an expected call of GetUserWithBalance.
func (mr *MockBalanceServiceMockRecorder) GetUserWithBalance(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserWithBalance", reflect.TypeOf((*MockBalanceService)(nil).GetUserWithBalance), userID)
}

// UpdateUserBalance mocks base method.
func (m *MockBalanceService) UpdateUserBalance(user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserBalance", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserBalance indicates an expected call of UpdateUserBalance.
func (mr *MockBalanceServiceMockRecorder) UpdateUserBalance(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserBalance", reflect.TypeOf((*MockBalanceService)(nil).UpdateUserBalance), user)
}