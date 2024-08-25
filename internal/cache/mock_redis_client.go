// Code generated by MockGen. DO NOT EDIT.
// Source: internal/cache/redis_client.go

// Package mocks is a generated GoMock package.
package cache

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCacheInterface is a mock of CacheInterface interface.
type MockCacheInterface struct {
	ctrl     *gomock.Controller
	recorder *MockCacheInterfaceMockRecorder
}

// MockCacheInterfaceMockRecorder is the mock recorder for MockCacheInterface.
type MockCacheInterfaceMockRecorder struct {
	mock *MockCacheInterface
}

// NewMockCacheInterface creates a new mock instance.
func NewMockCacheInterface(ctrl *gomock.Controller) *MockCacheInterface {
	mock := &MockCacheInterface{ctrl: ctrl}
	mock.recorder = &MockCacheInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCacheInterface) EXPECT() *MockCacheInterfaceMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockCacheInterface) Get(ctx context.Context, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCacheInterfaceMockRecorder) Get(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCacheInterface)(nil).Get), ctx, key)
}

// Set mocks base method.
func (m *MockCacheInterface) Set(ctx context.Context, key, value string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockCacheInterfaceMockRecorder) Set(ctx, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCacheInterface)(nil).Set), ctx, key, value)
}