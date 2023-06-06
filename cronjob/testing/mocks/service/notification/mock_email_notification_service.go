// Code generated by MockGen. DO NOT EDIT.
// Source: ../service/notification/email_notification_service.go

// Package mock_service_notification is a generated GoMock package.
package mock_service_notification

import (
	reflect "reflect"
	service_notification "github.com/scanyourkube/cronjob/service/notification"

	gomock "github.com/golang/mock/gomock"
)

// MockIEmailNotificationService is a mock of IEmailNotificationService interface.
type MockIEmailNotificationService struct {
	ctrl     *gomock.Controller
	recorder *MockIEmailNotificationServiceMockRecorder
}

// MockIEmailNotificationServiceMockRecorder is the mock recorder for MockIEmailNotificationService.
type MockIEmailNotificationServiceMockRecorder struct {
	mock *MockIEmailNotificationService
}

// NewMockIEmailNotificationService creates a new mock instance.
func NewMockIEmailNotificationService(ctrl *gomock.Controller) *MockIEmailNotificationService {
	mock := &MockIEmailNotificationService{ctrl: ctrl}
	mock.recorder = &MockIEmailNotificationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIEmailNotificationService) EXPECT() *MockIEmailNotificationServiceMockRecorder {
	return m.recorder
}

// SendEmail mocks base method.
func (m *MockIEmailNotificationService) SendEmail(email service_notification.EmailNotification) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmail", email)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockIEmailNotificationServiceMockRecorder) SendEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockIEmailNotificationService)(nil).SendEmail), email)
}