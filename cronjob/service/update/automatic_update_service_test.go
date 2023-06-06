package update_service

import (
	"github.com/scanyourkube/cronjob/dto/service/notification"
	dto_service_resource "github.com/scanyourkube/cronjob/dto/service/resource"
	mock_service_kubernetes "github.com/scanyourkube/cronjob/testing/mocks/service/kubernetes"
	mock_service_notification "github.com/scanyourkube/cronjob/testing/mocks/service/notification"
	mock_service_resource "github.com/scanyourkube/cronjob/testing/mocks/service/resource"
	mock_service_scan "github.com/scanyourkube/cronjob/testing/mocks/service/scan"
	mock_service_update "github.com/scanyourkube/cronjob/testing/mocks/service/update"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func vul(vul dto_service_resource.VulnerabilitySeverity) *dto_service_resource.VulnerabilitySeverity {
	return &vul
}
func attackVector(attackVector dto_service_resource.AttackVector) *dto_service_resource.AttackVector {
	return &attackVector
}

func TestAutomaticUpdateService_StartAutomaticUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScanService := mock_service_scan.NewMockIScanService(ctrl)
	mockKubernetesResourceService := mock_service_kubernetes.NewMockIKubernetesResourceService(ctrl)
	mockKubeClarityResourceService := mock_service_resource.NewMockIResourceService(ctrl)
	mockKeelUpdateService := mock_service_update.NewMockIKeelUpdateService(ctrl)
	mockNotificationService := mock_service_notification.NewMockIEmailNotificationService(ctrl)

	cases := []struct {
		input       []dto_service_resource.Application
		expectedErr bool
	}{
		{
			input: []dto_service_resource.Application{
				{
					ApplicationResources: map[string]dto_service_resource.ApplicationResource{
						"id": {
							Vulnerabilities: map[string]dto_service_resource.Vulnerability{
								"id": {
									CvssSeverity: vul(dto_service_resource.CRITICAL),
									AttackVector: attackVector(dto_service_resource.NETWORK),
									FixVersion:   "1.0.0",
								},
							},
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			input: []dto_service_resource.Application{
				{
					Labels: []string{"scanyourkube.io/owner=fluflo", "scanyourkube.io/podName=fluflo"},
					ApplicationResources: map[string]dto_service_resource.ApplicationResource{
						"id": {
							ResourceName: "trittale/trittale-nginx:latest",
							ResourceHash: "6c3e2af067405c14b199902b7541d6534b0bcf471f76565408327c67c723b6f5",
							Vulnerabilities: map[string]dto_service_resource.Vulnerability{
								"id": {
									CvssSeverity: vul(dto_service_resource.CRITICAL),
									AttackVector: attackVector(dto_service_resource.NETWORK),
									FixVersion:   "1.0.0",
								},
							},
						},
					},
				},
			},
			expectedErr: false,
		},
		{
			input: []dto_service_resource.Application{
				{
					ApplicationResources: map[string]dto_service_resource.ApplicationResource{
						"id": {
							Vulnerabilities: map[string]dto_service_resource.Vulnerability{},
						},
					},
				},
			},
			expectedErr: false,
		},
		{
			input: []dto_service_resource.Application{
				{
					ApplicationResources: map[string]dto_service_resource.ApplicationResource{},
				},
			},
			expectedErr: false,
		},
	}

	for _, c := range cases {
		// Set up expectations for the mocked dependencies
		mockScanService.EXPECT().ScanNamespaces().Return(nil)
		mockKubeClarityResourceService.EXPECT().GetApplicationsFromLastScan().Return(c.input, nil)
		// Create an instance of AutomaticUpdateService with mocked dependencies
		service := NewAutomaticUpdateService(
			mockKeelUpdateService,
			mockScanService,
			mockKubernetesResourceService,
			mockNotificationService,
			mockKubeClarityResourceService,
		)
		// Call the method being tested
		err := service.StartAutomaticUpdate()

		// Assert the result
		if c.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

}

func TestAutomaticUpdateService_updateApplicationResourceAndNotifyOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKubernetesResourceService := mock_service_kubernetes.NewMockIKubernetesResourceService(ctrl)
	mockKubeClarityResourceService := mock_service_resource.NewMockIResourceService(ctrl)
	mockKeelUpdateService := mock_service_update.NewMockIKeelUpdateService(ctrl)
	mockNotificationService := mock_service_notification.NewMockIEmailNotificationService(ctrl)

	// Create an instance of AutomaticUpdateService with mocked dependencies
	service := NewAutomaticUpdateService(
		mockKeelUpdateService,
		nil,
		mockKubernetesResourceService,
		mockNotificationService,
		mockKubeClarityResourceService,
	)

	// Call the method being tested
	err := service.updateApplicationResourceAndNotifyOwner(dto_service_resource.Application{}, dto_service_resource.ApplicationResource{
		ResourceName: "trittale/trittale-nginx:latest",
		ResourceHash: "6c3e2af067405c14b199902b7541d6534b0bcf471f76565408327c67c723b6f5",
		Vulnerabilities: map[string]dto_service_resource.Vulnerability{
			"id": {
				CvssSeverity: vul(dto_service_resource.CRITICAL),
				AttackVector: attackVector(dto_service_resource.NETWORK),
				FixVersion:   "1.0.0",
			},
		},
	})

	// Assert the result
	assert.NoError(t, err)
}

func TestAutomaticUpdateService_updateApplicationResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKeelUpdateService := mock_service_update.NewMockIKeelUpdateService(ctrl)

	// Set up expectations for the mocked dependency
	mockKeelUpdateService.EXPECT().Update(gomock.Any())
	// Create an instance of AutomaticUpdateService with mocked dependencies
	service := NewAutomaticUpdateService(
		mockKeelUpdateService,
		nil,
		nil,
		nil,
		nil,
	)

	// Call the method being tested
	_, err := service.updateApplicationResource(dto_service_resource.ApplicationResource{
		ResourceName: "trittale/trittale-nginx:latest",
		ResourceHash: "DUMMY_VALUE",
		Vulnerabilities: map[string]dto_service_resource.Vulnerability{
			"id": {
				CvssSeverity: vul(dto_service_resource.CRITICAL),
				AttackVector: attackVector(dto_service_resource.NETWORK),
				FixVersion:   "1.0.0",
			},
		},
	})

	// Assert the result
	assert.NoError(t, err)
}

func TestAutomaticUpdateService_sendNotificationToApplicationOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKubernetesResourceService := mock_service_kubernetes.NewMockIKubernetesResourceService(ctrl)
	mockNotificationService := mock_service_notification.NewMockIEmailNotificationService(ctrl)

	// Set up expectations for the mocked dependencies
	mockKubernetesResourceService.EXPECT().GetOwnerEmailAddressesByPodLabels(gomock.Any(), gomock.Any()).Return("dummy@dummy.com", nil).Times(1)
	mockNotificationService.EXPECT().SendEmail(gomock.Any()).Return(nil).Times(1)

	// Create an instance of AutomaticUpdateService with mocked dependencies
	service := NewAutomaticUpdateService(
		nil,
		nil,
		mockKubernetesResourceService,
		mockNotificationService,
		nil,
	)

	// Call the method being tested
	err := service.sendNotificationToApplicationOwner(dto_service_resource.Application{
		Environments: []string{"dev"},
		Labels:       []string{"scanyourkube.io/owner=fluflo", "scanyourkube.io/podName=fluflo"},
	}, notification.OwnerNotification{
		ApplicationName:          "fluflo",
		ApplicationResourceName:  "trittale/trittale-nginx:latest",
		VulnerabilityName:        "CVE-2020-1234",
		VulnerabilityDescription: "This is a description",
		NewVersion:               "1.0.0",
		ShouldUpdate:             true,
	})

	// Assert the result
	assert.NoError(t, err)
}

func TestAutomaticUpdateService_shouldApplicationResourceUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create an instance of AutomaticUpdateService with mocked dependencies
	service := NewAutomaticUpdateService(
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	// Call the method being tested
	_, shouldUpdate := service.shouldApplicationResourceUpdate(dto_service_resource.ApplicationResource{
		Vulnerabilities: map[string]dto_service_resource.Vulnerability{
			"id": {
				CvssSeverity: vul(dto_service_resource.CRITICAL),
				AttackVector: attackVector(dto_service_resource.NETWORK),
				FixVersion:   "1.0.0",
			},
		},
	})

	// Assert the result
	assert.True(t, shouldUpdate)
}
