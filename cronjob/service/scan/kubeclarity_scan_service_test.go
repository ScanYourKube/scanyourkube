package service_scan

import (
	"errors"
	"github.com/scanyourkube/cronjob/dto/api/scan"
	kubeclarity_api_runtime "github.com/scanyourkube/cronjob/testing/mocks/api/kubeclarity"
	kubernetes_resource_service "github.com/scanyourkube/cronjob/testing/mocks/service/kubernetes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSunnyCaseForService(t *testing.T) {
	//arrange
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockKubeclarityClient := kubeclarity_api_runtime.NewMockIKubeClarityApi(mockCtrl)
	mockKubernetesResourceService := kubernetes_resource_service.NewMockIKubernetesResourceService(mockCtrl)

	var service = NewScanService(mockKubeclarityClient, mockKubernetesResourceService)
	mockKubernetesResourceService.EXPECT().GetNamespaceNames().Return(make([]string, 0)).Times(1)
	mockKubeclarityClient.EXPECT().StartScan(gomock.Any()).Return(nil).Times(1)
	mockKubeclarityClient.EXPECT().GetScanProgress().Return(&scan.ScanProgress{Scanned: 100}, nil).Times(1)

	//act
	service.ScanNamespaces()

	//assert
	mockCtrl.Finish()
}

func TestSlowScanProgress(t *testing.T) {
	//arrange
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockKubeclarityClient := kubeclarity_api_runtime.NewMockIKubeClarityApi(mockCtrl)
	mockKubernetesResourceService := kubernetes_resource_service.NewMockIKubernetesResourceService(mockCtrl)

	var service = NewScanService(mockKubeclarityClient, mockKubernetesResourceService)
	mockKubernetesResourceService.EXPECT().GetNamespaceNames().Return(make([]string, 0)).Times(1)
	mockKubeclarityClient.EXPECT().StartScan(gomock.Any()).Return(nil).Times(1)
	mockKubeclarityClient.EXPECT().GetScanProgress().Return(&scan.ScanProgress{Scanned: 10}, nil).Times(1)
	mockKubeclarityClient.EXPECT().GetScanProgress().Return(&scan.ScanProgress{Scanned: 100}, nil).Times(1)

	//act
	err := service.ScanNamespaces()

	//assert
	assert.Nil(t, err)
	mockCtrl.Finish()
}

func TestWhenNoScanProgressIsReturned(t *testing.T) {
	//arrange
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockKubeclarityClient := kubeclarity_api_runtime.NewMockIKubeClarityApi(mockCtrl)
	mockKubernetesResourceService := kubernetes_resource_service.NewMockIKubernetesResourceService(mockCtrl)

	var service = NewScanService(mockKubeclarityClient, mockKubernetesResourceService)
	mockKubernetesResourceService.EXPECT().GetNamespaceNames().Return(make([]string, 0)).Times(1)
	mockKubeclarityClient.EXPECT().StartScan(gomock.Any()).Return(nil).Times(1)
	mockKubeclarityClient.EXPECT().GetScanProgress().Return(nil, errors.New("HTTP error")).Times(1)

	//act
	err := service.ScanNamespaces()

	//assert
	assert.NotNil(t, err)
	mockCtrl.Finish()
}
