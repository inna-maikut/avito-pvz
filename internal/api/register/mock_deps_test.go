// Code generated by MockGen. DO NOT EDIT.
// Source: deps.go
//
// Generated by this command:
//
//	mockgen -source deps.go -package register -typed -destination mock_deps_test.go
//

// Package register is a generated GoMock package.
package register

import (
	context "context"
	reflect "reflect"

	model "github.com/inna-maikut/avito-pvz/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// Mockregistering is a mock of registering interface.
type Mockregistering struct {
	ctrl     *gomock.Controller
	recorder *MockregisteringMockRecorder
	isgomock struct{}
}

// MockregisteringMockRecorder is the mock recorder for Mockregistering.
type MockregisteringMockRecorder struct {
	mock *Mockregistering
}

// NewMockregistering creates a new mock instance.
func NewMockregistering(ctrl *gomock.Controller) *Mockregistering {
	mock := &Mockregistering{ctrl: ctrl}
	mock.recorder = &MockregisteringMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockregistering) EXPECT() *MockregisteringMockRecorder {
	return m.recorder
}

// Register mocks base method.
func (m *Mockregistering) Register(ctx context.Context, email, password string, role model.UserRole) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, email, password, role)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockregisteringMockRecorder) Register(ctx, email, password, role any) *MockregisteringRegisterCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*Mockregistering)(nil).Register), ctx, email, password, role)
	return &MockregisteringRegisterCall{Call: call}
}

// MockregisteringRegisterCall wrap *gomock.Call
type MockregisteringRegisterCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockregisteringRegisterCall) Return(arg0 *model.User, arg1 error) *MockregisteringRegisterCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockregisteringRegisterCall) Do(f func(context.Context, string, string, model.UserRole) (*model.User, error)) *MockregisteringRegisterCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockregisteringRegisterCall) DoAndReturn(f func(context.Context, string, string, model.UserRole) (*model.User, error)) *MockregisteringRegisterCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
