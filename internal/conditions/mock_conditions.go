// Code generated by MockGen. DO NOT EDIT.
// Source: conditions.go

// Package conditions is a generated GoMock package.
package conditions

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/mresvanis/he-sample-operator/api/v1alpha1"
)

// MockUpdater is a mock of Updater interface.
type MockUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockUpdaterMockRecorder
}

// MockUpdaterMockRecorder is the mock recorder for MockUpdater.
type MockUpdaterMockRecorder struct {
	mock *MockUpdater
}

// NewMockUpdater creates a new mock instance.
func NewMockUpdater(ctrl *gomock.Controller) *MockUpdater {
	mock := &MockUpdater{ctrl: ctrl}
	mock.recorder = &MockUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUpdater) EXPECT() *MockUpdaterMockRecorder {
	return m.recorder
}

// SetConditionsErrored mocks base method.
func (m *MockUpdater) SetConditionsErrored(ctx context.Context, cr *v1alpha1.DeviceConfig, reason, message string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetConditionsErrored", ctx, cr, reason, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetConditionsErrored indicates an expected call of SetConditionsErrored.
func (mr *MockUpdaterMockRecorder) SetConditionsErrored(ctx, cr, reason, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConditionsErrored", reflect.TypeOf((*MockUpdater)(nil).SetConditionsErrored), ctx, cr, reason, message)
}

// SetConditionsReady mocks base method.
func (m *MockUpdater) SetConditionsReady(ctx context.Context, cr *v1alpha1.DeviceConfig, reason, message string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetConditionsReady", ctx, cr, reason, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetConditionsReady indicates an expected call of SetConditionsReady.
func (mr *MockUpdaterMockRecorder) SetConditionsReady(ctx, cr, reason, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConditionsReady", reflect.TypeOf((*MockUpdater)(nil).SetConditionsReady), ctx, cr, reason, message)
}
