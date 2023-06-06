package resource_service

import (
	api "github.com/scanyourkube/cronjob/api/kubeclarity"
	resource_service_dto "github.com/scanyourkube/cronjob/dto/service/resource"

	log "github.com/sirupsen/logrus"
)

type ResourceService struct {
	ScanApi api.IKubeClarityApi
}

type IResourceService interface {
	GetApplicationsFromLastScan() ([]resource_service_dto.Application, error)
}

func NewResourceService(ScanApi api.IKubeClarityApi) ResourceService {
	return ResourceService{ScanApi: ScanApi}
}

func (s ResourceService) GetApplicationsFromLastScan() ([]resource_service_dto.Application, error) {
	log.Debug("Starting to get applications with vulnerabilities of last scan")

	apiApplications, err := s.ScanApi.GetApplications()

	log.Debugf("Got %d applications", len(apiApplications))
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	var applications []resource_service_dto.Application
	for _, apiApplication := range apiApplications {

		applicationResources, err := s.getApplicationResourcesByApplicationId(apiApplication.Id)
		if err != nil {
			return nil, err
		}
		application := resource_service_dto.Application{}.Map(apiApplication)
		application.ApplicationResources = applicationResources
		applications = append(applications, application)
	}

	return applications, nil
}

func (s ResourceService) getApplicationResourcesByApplicationId(applicationId string) (map[string]resource_service_dto.ApplicationResource, error) {
	apiApplicationResources, err := s.ScanApi.GetApplicationResourceByApplicationId(applicationId)

	if err != nil {
		return nil, err
	}

	applicationResources := map[string]resource_service_dto.ApplicationResource{}

	for _, apiApplicationResource := range apiApplicationResources {
		vulnerabilityDetails, err := s.getVulnerabilitiesOfApplicationResource(apiApplicationResource.Id)
		if err != nil {
			return nil, err
		}

		applicationResource := resource_service_dto.ApplicationResource{}.Map(apiApplicationResource)
		applicationResource.Vulnerabilities = vulnerabilityDetails

		applicationResources[applicationResource.Id] = applicationResource
	}

	return applicationResources, nil
}

func (s ResourceService) getVulnerabilitiesOfApplicationResource(applicationResourceId string) (map[string]resource_service_dto.Vulnerability, error) {
	vulnerabilites, err := s.ScanApi.GetVulnerabilitiesOfApplicationResource(applicationResourceId)
	if err != nil {
		return nil, err
	}

	log.Debugf("Found %d vulnerabilites for resource with id: %s ", len(vulnerabilites), applicationResourceId)

	vulnerabilityDetails := map[string]resource_service_dto.Vulnerability{}
	for _, vulnerability := range vulnerabilites {
		apiVulnerabilityDetail, err := s.ScanApi.GetVulnerabilityDetailsWithPackage(vulnerability.VulnerabilityID, vulnerability.PackageID)
		if err != nil {
			return nil, err
		}
		vulnerabilityDetail := resource_service_dto.Vulnerability{}.Map(*apiVulnerabilityDetail)
		vulnerabilityDetails[vulnerabilityDetail.Id] = vulnerabilityDetail
	}

	return vulnerabilityDetails, nil
}
