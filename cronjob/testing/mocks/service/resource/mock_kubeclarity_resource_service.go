// Code generated by MockGen. DO NOT EDIT.
// Source: ../service/resource/kubeclarity_resource_service.go

// Package mock_resource_service is a generated GoMock package.
package mock_resource_service

import (
	reflect "reflect"
	resource_service "github.com/scanyourkube/cronjob/dto/service/resource"

	gomock "github.com/golang/mock/gomock"
)

// MockIResourceService is a mock of IResourceService interface.
type MockIResourceService struct {
	ctrl     *gomock.Controller
	recorder *MockIResourceServiceMockRecorder
}

// MockIResourceServiceMockRecorder is the mock recorder for MockIResourceService.
type MockIResourceServiceMockRecorder struct {
	mock *MockIResourceService
}

// NewMockIResourceService creates a new mock instance.
func NewMockIResourceService(ctrl *gomock.Controller) *MockIResourceService {
	mock := &MockIResourceService{ctrl: ctrl}
	mock.recorder = &MockIResourceServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIResourceService) EXPECT() *MockIResourceServiceMockRecorder {
	return m.recorder
}

// GetApplicationsFromLastScan mocks base method.
func (m *MockIResourceService) GetApplicationsFromLastScan() ([]resource_service.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApplicationsFromLastScan")
	ret0, _ := ret[0].([]resource_service.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApplicationsFromLastScan indicates an expected call of GetApplicationsFromLastScan.
func (mr *MockIResourceServiceMockRecorder) GetApplicationsFromLastScan() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApplicationsFromLastScan", reflect.TypeOf((*MockIResourceService)(nil).GetApplicationsFromLastScan))
}
