// Code generated by MockGen. DO NOT EDIT.
// Source: deps.go
//
// Generated by this command:
//
//	mockgen -source deps.go -package dummy_authenticating -typed -destination mock_deps_test.go
//

// Package dummy_authenticating is a generated GoMock package.
package dummy_authenticating

import (
	reflect "reflect"

	model "github.com/inna-maikut/avito-pvz/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MocktokenProvider is a mock of tokenProvider interface.
type MocktokenProvider struct {
	ctrl     *gomock.Controller
	recorder *MocktokenProviderMockRecorder
	isgomock struct{}
}

// MocktokenProviderMockRecorder is the mock recorder for MocktokenProvider.
type MocktokenProviderMockRecorder struct {
	mock *MocktokenProvider
}

// NewMocktokenProvider creates a new mock instance.
func NewMocktokenProvider(ctrl *gomock.Controller) *MocktokenProvider {
	mock := &MocktokenProvider{ctrl: ctrl}
	mock.recorder = &MocktokenProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocktokenProvider) EXPECT() *MocktokenProviderMockRecorder {
	return m.recorder
}

// CreateToken mocks base method.
func (m *MocktokenProvider) CreateToken(email string, userID int64, role model.UserRole) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateToken", email, userID, role)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateToken indicates an expected call of CreateToken.
func (mr *MocktokenProviderMockRecorder) CreateToken(email, userID, role any) *MocktokenProviderCreateTokenCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateToken", reflect.TypeOf((*MocktokenProvider)(nil).CreateToken), email, userID, role)
	return &MocktokenProviderCreateTokenCall{Call: call}
}

// MocktokenProviderCreateTokenCall wrap *gomock.Call
type MocktokenProviderCreateTokenCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MocktokenProviderCreateTokenCall) Return(arg0 string, arg1 error) *MocktokenProviderCreateTokenCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MocktokenProviderCreateTokenCall) Do(f func(string, int64, model.UserRole) (string, error)) *MocktokenProviderCreateTokenCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MocktokenProviderCreateTokenCall) DoAndReturn(f func(string, int64, model.UserRole) (string, error)) *MocktokenProviderCreateTokenCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
