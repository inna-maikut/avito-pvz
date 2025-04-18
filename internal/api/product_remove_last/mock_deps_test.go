// Code generated by MockGen. DO NOT EDIT.
// Source: deps.go
//
// Generated by this command:
//
//	mockgen -source deps.go -package product_remove_last -typed -destination mock_deps_test.go
//

// Package product_remove_last is a generated GoMock package.
package product_remove_last

import (
	context "context"
	reflect "reflect"

	model "github.com/inna-maikut/avito-pvz/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockproductRemoving is a mock of productRemoving interface.
type MockproductRemoving struct {
	ctrl     *gomock.Controller
	recorder *MockproductRemovingMockRecorder
	isgomock struct{}
}

// MockproductRemovingMockRecorder is the mock recorder for MockproductRemoving.
type MockproductRemovingMockRecorder struct {
	mock *MockproductRemoving
}

// NewMockproductRemoving creates a new mock instance.
func NewMockproductRemoving(ctrl *gomock.Controller) *MockproductRemoving {
	mock := &MockproductRemoving{ctrl: ctrl}
	mock.recorder = &MockproductRemovingMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockproductRemoving) EXPECT() *MockproductRemovingMockRecorder {
	return m.recorder
}

// RemoveLastProduct mocks base method.
func (m *MockproductRemoving) RemoveLastProduct(ctx context.Context, pvzID model.PVZID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveLastProduct", ctx, pvzID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveLastProduct indicates an expected call of RemoveLastProduct.
func (mr *MockproductRemovingMockRecorder) RemoveLastProduct(ctx, pvzID any) *MockproductRemovingRemoveLastProductCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveLastProduct", reflect.TypeOf((*MockproductRemoving)(nil).RemoveLastProduct), ctx, pvzID)
	return &MockproductRemovingRemoveLastProductCall{Call: call}
}

// MockproductRemovingRemoveLastProductCall wrap *gomock.Call
type MockproductRemovingRemoveLastProductCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockproductRemovingRemoveLastProductCall) Return(arg0 error) *MockproductRemovingRemoveLastProductCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockproductRemovingRemoveLastProductCall) Do(f func(context.Context, model.PVZID) error) *MockproductRemovingRemoveLastProductCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockproductRemovingRemoveLastProductCall) DoAndReturn(f func(context.Context, model.PVZID) error) *MockproductRemovingRemoveLastProductCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
