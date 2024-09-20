// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/service/order_service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/golangTroshin/gophermat/internal/model"
)

// MockOrderService is a mock of OrderService interface.
type MockOrderService struct {
	ctrl     *gomock.Controller
	recorder *MockOrderServiceMockRecorder
}

// MockOrderServiceMockRecorder is the mock recorder for MockOrderService.
type MockOrderServiceMockRecorder struct {
	mock *MockOrderService
}

// NewMockOrderService creates a new mock instance.
func NewMockOrderService(ctrl *gomock.Controller) *MockOrderService {
	mock := &MockOrderService{ctrl: ctrl}
	mock.recorder = &MockOrderServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderService) EXPECT() *MockOrderServiceMockRecorder {
	return m.recorder
}

// CreateOrder mocks base method.
func (m *MockOrderService) CreateOrder(userID uint, orderNumber string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", userID, orderNumber)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockOrderServiceMockRecorder) CreateOrder(userID, orderNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockOrderService)(nil).CreateOrder), userID, orderNumber)
}

// GetOrdersByUserID mocks base method.
func (m *MockOrderService) GetOrdersByUserID(userID uint) ([]model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersByUserID", userID)
	ret0, _ := ret[0].([]model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersByUserID indicates an expected call of GetOrdersByUserID.
func (mr *MockOrderServiceMockRecorder) GetOrdersByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersByUserID", reflect.TypeOf((*MockOrderService)(nil).GetOrdersByUserID), userID)
}

// ValidateOrderNumber mocks base method.
func (m *MockOrderService) ValidateOrderNumber(orderNumber string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateOrderNumber", orderNumber)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ValidateOrderNumber indicates an expected call of ValidateOrderNumber.
func (mr *MockOrderServiceMockRecorder) ValidateOrderNumber(orderNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateOrderNumber", reflect.TypeOf((*MockOrderService)(nil).ValidateOrderNumber), orderNumber)
}