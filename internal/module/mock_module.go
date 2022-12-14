// Code generated by MockGen. DO NOT EDIT.
// Source: module.go

// Package module is a generated GoMock package.
package module

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1beta1 "github.com/kubernetes-sigs/kernel-module-management/api/v1beta1"
	v1alpha1 "github.com/mresvanis/he-sample-operator/api/v1alpha1"
)

// MockReconciler is a mock of Reconciler interface.
type MockReconciler struct {
	ctrl     *gomock.Controller
	recorder *MockReconcilerMockRecorder
}

// MockReconcilerMockRecorder is the mock recorder for MockReconciler.
type MockReconcilerMockRecorder struct {
	mock *MockReconciler
}

// NewMockReconciler creates a new mock instance.
func NewMockReconciler(ctrl *gomock.Controller) *MockReconciler {
	mock := &MockReconciler{ctrl: ctrl}
	mock.recorder = &MockReconcilerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReconciler) EXPECT() *MockReconcilerMockRecorder {
	return m.recorder
}

// DeleteModule mocks base method.
func (m *MockReconciler) DeleteModule(ctx context.Context, dc *v1alpha1.DeviceConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteModule", ctx, dc)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteModule indicates an expected call of DeleteModule.
func (mr *MockReconcilerMockRecorder) DeleteModule(ctx, dc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteModule", reflect.TypeOf((*MockReconciler)(nil).DeleteModule), ctx, dc)
}

// ReconcileModule mocks base method.
func (m *MockReconciler) ReconcileModule(ctx context.Context, dc *v1alpha1.DeviceConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReconcileModule", ctx, dc)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReconcileModule indicates an expected call of ReconcileModule.
func (mr *MockReconcilerMockRecorder) ReconcileModule(ctx, dc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReconcileModule", reflect.TypeOf((*MockReconciler)(nil).ReconcileModule), ctx, dc)
}

// SetDesiredModule mocks base method.
func (m_2 *MockReconciler) SetDesiredModule(m *v1beta1.Module, cr *v1alpha1.DeviceConfig) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SetDesiredModule", m, cr)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDesiredModule indicates an expected call of SetDesiredModule.
func (mr *MockReconcilerMockRecorder) SetDesiredModule(m, cr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDesiredModule", reflect.TypeOf((*MockReconciler)(nil).SetDesiredModule), m, cr)
}
