// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	model "route256/checkout/internal/domain/model"

	gomock "github.com/golang/mock/gomock"
)

// MockCheckoutRepository is a mock of CheckoutRepository interface.
type MockCheckoutRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCheckoutRepositoryMockRecorder
}

// MockCheckoutRepositoryMockRecorder is the mock recorder for MockCheckoutRepository.
type MockCheckoutRepositoryMockRecorder struct {
	mock *MockCheckoutRepository
}

// NewMockCheckoutRepository creates a new mock instance.
func NewMockCheckoutRepository(ctrl *gomock.Controller) *MockCheckoutRepository {
	mock := &MockCheckoutRepository{ctrl: ctrl}
	mock.recorder = &MockCheckoutRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCheckoutRepository) EXPECT() *MockCheckoutRepositoryMockRecorder {
	return m.recorder
}

// AddToCart mocks base method.
func (m *MockCheckoutRepository) AddToCart(ctx context.Context, addToCartRequest *model.AddToCartRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToCart", ctx, addToCartRequest)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToCart indicates an expected call of AddToCart.
func (mr *MockCheckoutRepositoryMockRecorder) AddToCart(ctx, addToCartRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToCart", reflect.TypeOf((*MockCheckoutRepository)(nil).AddToCart), ctx, addToCartRequest)
}

// DeleteFromCart mocks base method.
func (m *MockCheckoutRepository) DeleteFromCart(ctx context.Context, deleteFromCartRequest *model.DeleteFromCartRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFromCart", ctx, deleteFromCartRequest)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFromCart indicates an expected call of DeleteFromCart.
func (mr *MockCheckoutRepositoryMockRecorder) DeleteFromCart(ctx, deleteFromCartRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFromCart", reflect.TypeOf((*MockCheckoutRepository)(nil).DeleteFromCart), ctx, deleteFromCartRequest)
}

// ListCart mocks base method.
func (m *MockCheckoutRepository) ListCart(ctx context.Context, listCartRequest *model.ListCartRequest) (*model.ListCartResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCart", ctx, listCartRequest)
	ret0, _ := ret[0].(*model.ListCartResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCart indicates an expected call of ListCart.
func (mr *MockCheckoutRepositoryMockRecorder) ListCart(ctx, listCartRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCart", reflect.TypeOf((*MockCheckoutRepository)(nil).ListCart), ctx, listCartRequest)
}
