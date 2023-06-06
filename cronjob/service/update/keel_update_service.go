package update_service

import (
	api_keel_webhook "github.com/scanyourkube/cronjob/api/keel"
	update_keel "github.com/scanyourkube/cronjob/dto/api/update"
	registry_image "github.com/scanyourkube/cronjob/service/update/registryImage"

	log "github.com/sirupsen/logrus"
)

type KeelUpdateService struct {
	UpdateWebHookApi api_keel_webhook.IKeelWebhookApi
}

type IKeelUpdateService interface {
	Update(registryImage registry_image.RegistryImage)
}

func NewKeelUpdateService(UpdateWebHookApi api_keel_webhook.IKeelWebhookApi) KeelUpdateService {
	return KeelUpdateService{
		UpdateWebHookApi: UpdateWebHookApi,
	}
}

func (service KeelUpdateService) Update(registryImage registry_image.RegistryImage) {
	log.Debugf("Updating %s to version %s \n", registryImage.ImageName, registryImage.Tag)
	service.UpdateWebHookApi.UpdateImage(update_keel.KeelWebhookDto{
		Name: registryImage.ImageName,
		Tag:  registryImage.Tag,
	})
}
