// Code generated by MockGen. DO NOT EDIT.
// Source: deps.go
//
// Generated by this command:
//
//	mockgen -source deps.go -package reception_close -typed -destination mock_deps_test.go
//

// Package reception_close is a generated GoMock package.
package reception_close

import (
	context "context"
	reflect "reflect"

	model "github.com/inna-maikut/avito-pvz/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockreceptionClosing is a mock of receptionClosing interface.
type MockreceptionClosing struct {
	ctrl     *gomock.Controller
	recorder *MockreceptionClosingMockRecorder
	isgomock struct{}
}

// MockreceptionClosingMockRecorder is the mock recorder for MockreceptionClosing.
type MockreceptionClosingMockRecorder struct {
	mock *MockreceptionClosing
}

// NewMockreceptionClosing creates a new mock instance.
func NewMockreceptionClosing(ctrl *gomock.Controller) *MockreceptionClosing {
	mock := &MockreceptionClosing{ctrl: ctrl}
	mock.recorder = &MockreceptionClosingMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockreceptionClosing) EXPECT() *MockreceptionClosingMockRecorder {
	return m.recorder
}

// CloseReception mocks base method.
func (m *MockreceptionClosing) CloseReception(ctx context.Context, pvzID model.PVZID) (model.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseReception", ctx, pvzID)
	ret0, _ := ret[0].(model.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloseReception indicates an expected call of CloseReception.
func (mr *MockreceptionClosingMockRecorder) CloseReception(ctx, pvzID any) *MockreceptionClosingCloseReceptionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseReception", reflect.TypeOf((*MockreceptionClosing)(nil).CloseReception), ctx, pvzID)
	return &MockreceptionClosingCloseReceptionCall{Call: call}
}

// MockreceptionClosingCloseReceptionCall wrap *gomock.Call
type MockreceptionClosingCloseReceptionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockreceptionClosingCloseReceptionCall) Return(arg0 model.Reception, arg1 error) *MockreceptionClosingCloseReceptionCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockreceptionClosingCloseReceptionCall) Do(f func(context.Context, model.PVZID) (model.Reception, error)) *MockreceptionClosingCloseReceptionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockreceptionClosingCloseReceptionCall) DoAndReturn(f func(context.Context, model.PVZID) (model.Reception, error)) *MockreceptionClosingCloseReceptionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
