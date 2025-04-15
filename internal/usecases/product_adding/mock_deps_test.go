// Code generated by MockGen. DO NOT EDIT.
// Source: deps.go
//
// Generated by this command:
//
//	mockgen -source deps.go -package product_adding -typed -destination mock_deps_test.go
//

// Package product_adding is a generated GoMock package.
package product_adding

import (
	context "context"
	reflect "reflect"

	model "github.com/inna-maikut/avito-pvz/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MocktrManager is a mock of trManager interface.
type MocktrManager struct {
	ctrl     *gomock.Controller
	recorder *MocktrManagerMockRecorder
	isgomock struct{}
}

// MocktrManagerMockRecorder is the mock recorder for MocktrManager.
type MocktrManagerMockRecorder struct {
	mock *MocktrManager
}

// NewMocktrManager creates a new mock instance.
func NewMocktrManager(ctrl *gomock.Controller) *MocktrManager {
	mock := &MocktrManager{ctrl: ctrl}
	mock.recorder = &MocktrManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocktrManager) EXPECT() *MocktrManagerMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MocktrManager) Do(ctx context.Context, fn func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// Do indicates an expected call of Do.
func (mr *MocktrManagerMockRecorder) Do(ctx, fn any) *MocktrManagerDoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MocktrManager)(nil).Do), ctx, fn)
	return &MocktrManagerDoCall{Call: call}
}

// MocktrManagerDoCall wrap *gomock.Call
type MocktrManagerDoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MocktrManagerDoCall) Return(err error) *MocktrManagerDoCall {
	c.Call = c.Call.Return(err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MocktrManagerDoCall) Do(f func(context.Context, func(context.Context) error) error) *MocktrManagerDoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MocktrManagerDoCall) DoAndReturn(f func(context.Context, func(context.Context) error) error) *MocktrManagerDoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockreceptionRepo is a mock of receptionRepo interface.
type MockreceptionRepo struct {
	ctrl     *gomock.Controller
	recorder *MockreceptionRepoMockRecorder
	isgomock struct{}
}

// MockreceptionRepoMockRecorder is the mock recorder for MockreceptionRepo.
type MockreceptionRepoMockRecorder struct {
	mock *MockreceptionRepo
}

// NewMockreceptionRepo creates a new mock instance.
func NewMockreceptionRepo(ctrl *gomock.Controller) *MockreceptionRepo {
	mock := &MockreceptionRepo{ctrl: ctrl}
	mock.recorder = &MockreceptionRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockreceptionRepo) EXPECT() *MockreceptionRepoMockRecorder {
	return m.recorder
}

// GetInProgress mocks base method.
func (m *MockreceptionRepo) GetInProgress(ctx context.Context, pvzID model.PVZID) (model.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInProgress", ctx, pvzID)
	ret0, _ := ret[0].(model.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInProgress indicates an expected call of GetInProgress.
func (mr *MockreceptionRepoMockRecorder) GetInProgress(ctx, pvzID any) *MockreceptionRepoGetInProgressCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInProgress", reflect.TypeOf((*MockreceptionRepo)(nil).GetInProgress), ctx, pvzID)
	return &MockreceptionRepoGetInProgressCall{Call: call}
}

// MockreceptionRepoGetInProgressCall wrap *gomock.Call
type MockreceptionRepoGetInProgressCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockreceptionRepoGetInProgressCall) Return(arg0 model.Reception, arg1 error) *MockreceptionRepoGetInProgressCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockreceptionRepoGetInProgressCall) Do(f func(context.Context, model.PVZID) (model.Reception, error)) *MockreceptionRepoGetInProgressCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockreceptionRepoGetInProgressCall) DoAndReturn(f func(context.Context, model.PVZID) (model.Reception, error)) *MockreceptionRepoGetInProgressCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockproductRepo is a mock of productRepo interface.
type MockproductRepo struct {
	ctrl     *gomock.Controller
	recorder *MockproductRepoMockRecorder
	isgomock struct{}
}

// MockproductRepoMockRecorder is the mock recorder for MockproductRepo.
type MockproductRepoMockRecorder struct {
	mock *MockproductRepo
}

// NewMockproductRepo creates a new mock instance.
func NewMockproductRepo(ctrl *gomock.Controller) *MockproductRepo {
	mock := &MockproductRepo{ctrl: ctrl}
	mock.recorder = &MockproductRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockproductRepo) EXPECT() *MockproductRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockproductRepo) Create(ctx context.Context, receptionID model.ReceptionID, category model.ProductCategory) (model.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, receptionID, category)
	ret0, _ := ret[0].(model.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockproductRepoMockRecorder) Create(ctx, receptionID, category any) *MockproductRepoCreateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockproductRepo)(nil).Create), ctx, receptionID, category)
	return &MockproductRepoCreateCall{Call: call}
}

// MockproductRepoCreateCall wrap *gomock.Call
type MockproductRepoCreateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockproductRepoCreateCall) Return(arg0 model.Product, arg1 error) *MockproductRepoCreateCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockproductRepoCreateCall) Do(f func(context.Context, model.ReceptionID, model.ProductCategory) (model.Product, error)) *MockproductRepoCreateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockproductRepoCreateCall) DoAndReturn(f func(context.Context, model.ReceptionID, model.ProductCategory) (model.Product, error)) *MockproductRepoCreateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockpvzLocker is a mock of pvzLocker interface.
type MockpvzLocker struct {
	ctrl     *gomock.Controller
	recorder *MockpvzLockerMockRecorder
	isgomock struct{}
}

// MockpvzLockerMockRecorder is the mock recorder for MockpvzLocker.
type MockpvzLockerMockRecorder struct {
	mock *MockpvzLocker
}

// NewMockpvzLocker creates a new mock instance.
func NewMockpvzLocker(ctrl *gomock.Controller) *MockpvzLocker {
	mock := &MockpvzLocker{ctrl: ctrl}
	mock.recorder = &MockpvzLockerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockpvzLocker) EXPECT() *MockpvzLockerMockRecorder {
	return m.recorder
}

// Lock mocks base method.
func (m *MockpvzLocker) Lock(ctx context.Context, pvzID model.PVZID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Lock", ctx, pvzID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Lock indicates an expected call of Lock.
func (mr *MockpvzLockerMockRecorder) Lock(ctx, pvzID any) *MockpvzLockerLockCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Lock", reflect.TypeOf((*MockpvzLocker)(nil).Lock), ctx, pvzID)
	return &MockpvzLockerLockCall{Call: call}
}

// MockpvzLockerLockCall wrap *gomock.Call
type MockpvzLockerLockCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockpvzLockerLockCall) Return(arg0 error) *MockpvzLockerLockCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockpvzLockerLockCall) Do(f func(context.Context, model.PVZID) error) *MockpvzLockerLockCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockpvzLockerLockCall) DoAndReturn(f func(context.Context, model.PVZID) error) *MockpvzLockerLockCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Mockmetrics is a mock of metrics interface.
type Mockmetrics struct {
	ctrl     *gomock.Controller
	recorder *MockmetricsMockRecorder
	isgomock struct{}
}

// MockmetricsMockRecorder is the mock recorder for Mockmetrics.
type MockmetricsMockRecorder struct {
	mock *Mockmetrics
}

// NewMockmetrics creates a new mock instance.
func NewMockmetrics(ctrl *gomock.Controller) *Mockmetrics {
	mock := &Mockmetrics{ctrl: ctrl}
	mock.recorder = &MockmetricsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockmetrics) EXPECT() *MockmetricsMockRecorder {
	return m.recorder
}

// ProductAddedCountInc mocks base method.
func (m *Mockmetrics) ProductAddedCountInc() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ProductAddedCountInc")
}

// ProductAddedCountInc indicates an expected call of ProductAddedCountInc.
func (mr *MockmetricsMockRecorder) ProductAddedCountInc() *MockmetricsProductAddedCountIncCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProductAddedCountInc", reflect.TypeOf((*Mockmetrics)(nil).ProductAddedCountInc))
	return &MockmetricsProductAddedCountIncCall{Call: call}
}

// MockmetricsProductAddedCountIncCall wrap *gomock.Call
type MockmetricsProductAddedCountIncCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockmetricsProductAddedCountIncCall) Return() *MockmetricsProductAddedCountIncCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockmetricsProductAddedCountIncCall) Do(f func()) *MockmetricsProductAddedCountIncCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockmetricsProductAddedCountIncCall) DoAndReturn(f func()) *MockmetricsProductAddedCountIncCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
