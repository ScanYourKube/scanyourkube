package resource_service

import (
	api_dto_applicationResource "github.com/scanyourkube/cronjob/dto/api/application_resource"
)

type ApplicationResource struct {
	Id              string                   `json:"id,omitempty"`
	ResourceName    string                   `json:"resourceName,omitempty"`
	ResourceHash    string                   `json:"resourceHash,omitempty"`
	Vulnerabilities map[string]Vulnerability `json:"vulnerabilities,omitempty"`
}

func (dto ApplicationResource) Map(apiApplicationResource api_dto_applicationResource.ApplicationResource) ApplicationResource {
	return ApplicationResource{
		Id:           apiApplicationResource.Id,
		ResourceName: apiApplicationResource.ResourceName,
		ResourceHash: apiApplicationResource.ResourceHash,
	}
}
