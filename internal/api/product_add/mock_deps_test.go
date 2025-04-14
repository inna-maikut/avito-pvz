// Code generated by MockGen. DO NOT EDIT.
// Source: deps.go
//
// Generated by this command:
//
//	mockgen -source deps.go -package product_add -typed -destination mock_deps_test.go
//

// Package product_add is a generated GoMock package.
package product_add

import (
	context "context"
	reflect "reflect"

	model "github.com/inna-maikut/avito-pvz/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockproductAdding is a mock of productAdding interface.
type MockproductAdding struct {
	ctrl     *gomock.Controller
	recorder *MockproductAddingMockRecorder
	isgomock struct{}
}

// MockproductAddingMockRecorder is the mock recorder for MockproductAdding.
type MockproductAddingMockRecorder struct {
	mock *MockproductAdding
}

// NewMockproductAdding creates a new mock instance.
func NewMockproductAdding(ctrl *gomock.Controller) *MockproductAdding {
	mock := &MockproductAdding{ctrl: ctrl}
	mock.recorder = &MockproductAddingMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockproductAdding) EXPECT() *MockproductAddingMockRecorder {
	return m.recorder
}

// AddProduct mocks base method.
func (m *MockproductAdding) AddProduct(ctx context.Context, pvzID model.PVZID, category model.ProductCategory) (model.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProduct", ctx, pvzID, category)
	ret0, _ := ret[0].(model.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProduct indicates an expected call of AddProduct.
func (mr *MockproductAddingMockRecorder) AddProduct(ctx, pvzID, category any) *MockproductAddingAddProductCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProduct", reflect.TypeOf((*MockproductAdding)(nil).AddProduct), ctx, pvzID, category)
	return &MockproductAddingAddProductCall{Call: call}
}

// MockproductAddingAddProductCall wrap *gomock.Call
type MockproductAddingAddProductCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockproductAddingAddProductCall) Return(arg0 model.Product, arg1 error) *MockproductAddingAddProductCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockproductAddingAddProductCall) Do(f func(context.Context, model.PVZID, model.ProductCategory) (model.Product, error)) *MockproductAddingAddProductCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockproductAddingAddProductCall) DoAndReturn(f func(context.Context, model.PVZID, model.ProductCategory) (model.Product, error)) *MockproductAddingAddProductCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
