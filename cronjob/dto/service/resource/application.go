package resource_service

import (
	api_dto_application "github.com/scanyourkube/cronjob/dto/api/application"
)

type Application struct {
	Id                   string                         `json:"id,omitempty"`
	ApplicationName      string                         `json:"applicationName,omitempty"`
	Environments         []string                       `json:"environments,omitempty"`
	Labels               []string                       `json:"labels,omitempty"`
	ApplicationResources map[string]ApplicationResource `json:"applicationResources,omitempty"`
}

func (dto Application) Map(dtoApiApplication api_dto_application.Application) Application {

	return Application{
		Id:              dtoApiApplication.Id,
		ApplicationName: dtoApiApplication.ApplicationName,
		Environments:    dtoApiApplication.Environments,
		Labels:          dtoApiApplication.Labels,
	}
}
