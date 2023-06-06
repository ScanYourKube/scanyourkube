// Code generated by MockGen. DO NOT EDIT.
// Source: ../service/scan/kubeclarity_scan_service.go

// Package mock_service_scan is a generated GoMock package.
package mock_service_scan

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIScanService is a mock of IScanService interface.
type MockIScanService struct {
	ctrl     *gomock.Controller
	recorder *MockIScanServiceMockRecorder
}

// MockIScanServiceMockRecorder is the mock recorder for MockIScanService.
type MockIScanServiceMockRecorder struct {
	mock *MockIScanService
}

// NewMockIScanService creates a new mock instance.
func NewMockIScanService(ctrl *gomock.Controller) *MockIScanService {
	mock := &MockIScanService{ctrl: ctrl}
	mock.recorder = &MockIScanServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIScanService) EXPECT() *MockIScanServiceMockRecorder {
	return m.recorder
}

// GetScanProgressUntilFinished mocks base method.
func (m *MockIScanService) GetScanProgressUntilFinished() (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScanProgressUntilFinished")
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScanProgressUntilFinished indicates an expected call of GetScanProgressUntilFinished.
func (mr *MockIScanServiceMockRecorder) GetScanProgressUntilFinished() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScanProgressUntilFinished", reflect.TypeOf((*MockIScanService)(nil).GetScanProgressUntilFinished))
}

// ScanNamespaces mocks base method.
func (m *MockIScanService) ScanNamespaces() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScanNamespaces")
	ret0, _ := ret[0].(error)
	return ret0
}

// ScanNamespaces indicates an expected call of ScanNamespaces.
func (mr *MockIScanServiceMockRecorder) ScanNamespaces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScanNamespaces", reflect.TypeOf((*MockIScanService)(nil).ScanNamespaces))
}