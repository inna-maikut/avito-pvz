// Code generated by MockGen. DO NOT EDIT.
// Source: deps.go
//
// Generated by this command:
//
//	mockgen -source deps.go -package pvz_register -typed -destination mock_deps_test.go
//

// Package pvz_register is a generated GoMock package.
package pvz_register

import (
	context "context"
	reflect "reflect"

	model "github.com/inna-maikut/avito-pvz/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockpvzRegistering is a mock of pvzRegistering interface.
type MockpvzRegistering struct {
	ctrl     *gomock.Controller
	recorder *MockpvzRegisteringMockRecorder
	isgomock struct{}
}

// MockpvzRegisteringMockRecorder is the mock recorder for MockpvzRegistering.
type MockpvzRegisteringMockRecorder struct {
	mock *MockpvzRegistering
}

// NewMockpvzRegistering creates a new mock instance.
func NewMockpvzRegistering(ctrl *gomock.Controller) *MockpvzRegistering {
	mock := &MockpvzRegistering{ctrl: ctrl}
	mock.recorder = &MockpvzRegisteringMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockpvzRegistering) EXPECT() *MockpvzRegisteringMockRecorder {
	return m.recorder
}

// RegisterPVZ mocks base method.
func (m *MockpvzRegistering) RegisterPVZ(ctx context.Context, city string) (model.PVZ, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterPVZ", ctx, city)
	ret0, _ := ret[0].(model.PVZ)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterPVZ indicates an expected call of RegisterPVZ.
func (mr *MockpvzRegisteringMockRecorder) RegisterPVZ(ctx, city any) *MockpvzRegisteringRegisterPVZCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterPVZ", reflect.TypeOf((*MockpvzRegistering)(nil).RegisterPVZ), ctx, city)
	return &MockpvzRegisteringRegisterPVZCall{Call: call}
}

// MockpvzRegisteringRegisterPVZCall wrap *gomock.Call
type MockpvzRegisteringRegisterPVZCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockpvzRegisteringRegisterPVZCall) Return(arg0 model.PVZ, arg1 error) *MockpvzRegisteringRegisterPVZCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockpvzRegisteringRegisterPVZCall) Do(f func(context.Context, string) (model.PVZ, error)) *MockpvzRegisteringRegisterPVZCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockpvzRegisteringRegisterPVZCall) DoAndReturn(f func(context.Context, string) (model.PVZ, error)) *MockpvzRegisteringRegisterPVZCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
