package update_service

import (
	"testing"

	update_keel "github.com/scanyourkube/cronjob/dto/api/update"
	registry_image "github.com/scanyourkube/cronjob/service/update/registryImage"
	"github.com/stretchr/testify/assert"
)

type MockKeelWebhookApi struct {
	UpdateImageFunc func(dto update_keel.KeelWebhookDto) error
}

func (m *MockKeelWebhookApi) UpdateImage(dto update_keel.KeelWebhookDto) error {
	m.UpdateImageFunc(dto)
	return nil
}

func TestNewKeelUpdateService(t *testing.T) {
	// Create a mock KeelWebhookApi
	mockWebhookApi := &MockKeelWebhookApi{}

	// Create a new instance of KeelUpdateService
	service := NewKeelUpdateService(mockWebhookApi)

	// Verify that the UpdateWebHookApi field is set correctly
	assert.Equal(t, mockWebhookApi, service.UpdateWebHookApi)
}

func TestKeelUpdateServiceUpdate(t *testing.T) {
	// Create a mock KeelWebhookApi
	mockWebhookApi := &MockKeelWebhookApi{}

	// Create an instance of KeelUpdateService with the mock KeelWebhookApi
	service := KeelUpdateService{
		UpdateWebHookApi: mockWebhookApi,
	}

	// Create a sample registry image
	image := registry_image.RegistryImage{
		ImageName: "my-app",
		Tag:       "v1.0.0",
	}

	// Define a flag to check if the UpdateImageFunc is called
	isUpdateImageCalled := false

	// Set the UpdateImageFunc on the mock KeelWebhookApi
	mockWebhookApi.UpdateImageFunc = func(dto update_keel.KeelWebhookDto) error {
		// Verify that the UpdateImageFunc is called with the expected KeelWebhookDto
		assert.Equal(t, update_keel.KeelWebhookDto{
			Name: "my-app",
			Tag:  "v1.0.0",
		}, dto)

		isUpdateImageCalled = true
		return nil
	}

	// Call the Update method
	service.Update(image)

	// Verify that the UpdateImageFunc was called
	assert.True(t, isUpdateImageCalled)
}
