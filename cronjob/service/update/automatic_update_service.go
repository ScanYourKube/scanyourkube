package update_service

import (
	"time"

	"github.com/scanyourkube/cronjob/dto/service/notification"
	dto_service_resource "github.com/scanyourkube/cronjob/dto/service/resource"
	dto_service_update "github.com/scanyourkube/cronjob/dto/service/update"
	kubernetes_resource_service "github.com/scanyourkube/cronjob/service/kubernetes"
	notification_service "github.com/scanyourkube/cronjob/service/notification"
	kubeclarity_resource_service "github.com/scanyourkube/cronjob/service/resource"
	service_scan "github.com/scanyourkube/cronjob/service/scan"
	registry_image "github.com/scanyourkube/cronjob/service/update/registryImage"
	update_service_rules "github.com/scanyourkube/cronjob/service/update/rules"

	log "github.com/sirupsen/logrus"
)

type AutomaticUpdateService struct {
	keelUpdateService          IKeelUpdateService
	scanService                service_scan.IScanService
	kubernetesResourceService  kubernetes_resource_service.IKubernetesResourceService
	notificationService        notification_service.IEmailNotificationService
	kubeClarityResourceService kubeclarity_resource_service.IResourceService
}

func NewAutomaticUpdateService(
	keelUpdateService IKeelUpdateService,
	scanService service_scan.IScanService,
	kubernetesResourceService kubernetes_resource_service.IKubernetesResourceService,
	notificationService notification_service.IEmailNotificationService,
	kubeClarityResourceService kubeclarity_resource_service.IResourceService) AutomaticUpdateService {

	return AutomaticUpdateService{
		keelUpdateService:          keelUpdateService,
		scanService:                scanService,
		kubernetesResourceService:  kubernetesResourceService,
		notificationService:        notificationService,
		kubeClarityResourceService: kubeClarityResourceService,
	}
}

func (service AutomaticUpdateService) StartAutomaticUpdate() error {
	err := service.scanService.ScanNamespaces()
	log.Debugf("Scanned namespaces with error: %v", err)
	if err != nil {
		log.Debugf("Now returning as error is :%v", err)
		return err
	}

	// as the kubeclarity scan is async, we need to wait a bit to get the results
	time.Sleep(5 * time.Second)

	applications, err := service.kubeClarityResourceService.GetApplicationsFromLastScan()
	if err != nil {
		return err
	}

	for _, application := range applications {
		for _, applicationResource := range application.ApplicationResources {
			log.Debugf("Application resource %s in application %s has %d vulnerabilites",
				applicationResource.ResourceName,
				application.ApplicationName,
				len(applicationResource.Vulnerabilities))
			err := service.updateApplicationResourceAndNotifyOwner(application, applicationResource)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (service AutomaticUpdateService) updateApplicationResourceAndNotifyOwner(
	application dto_service_resource.Application,
	applicationResource dto_service_resource.ApplicationResource,
) error {
	if triggeringVulnerability,
		shouldUpdate := service.shouldApplicationResourceUpdate(applicationResource); shouldUpdate {
		registryImageStatus, err := service.updateApplicationResource(applicationResource)
		if err != nil {
			return err
		}
		log.Debugf("Registry image status: %v", registryImageStatus)
		emailNotification := notification.NewOwnerNotification(
			application.ApplicationName,
			applicationResource.ResourceName,
			triggeringVulnerability.VulnerabilityName,
			triggeringVulnerability.Description,
			registryImageStatus.RegistryImage.Tag,
			registryImageStatus.NewVersionAvailable,
		)
		err = service.sendNotificationToApplicationOwner(application, emailNotification)
		if err != nil {
			log.Errorf("There was an error sending notifications to the application owner: %v", err)
		}
	}
	return nil
}

func (service AutomaticUpdateService) updateApplicationResource(
	applicationResource dto_service_resource.ApplicationResource) (registry_image.ImageRegistryStatus, error) {
	updateImage := dto_service_update.UpdateServiceImageDto{
		Id:           applicationResource.Id,
		ResourceName: applicationResource.ResourceName,
		ResourceHash: applicationResource.ResourceHash,
	}

	registryImageStatus, err := registry_image.GetImageRegistryStatus(updateImage)

	if err != nil {
		return registryImageStatus, err
	}

	log.Debugf("Registry image status: %v", registryImageStatus)
	if registryImageStatus.NewVersionAvailable {
		service.keelUpdateService.Update(registryImageStatus.RegistryImage)
	}

	return registryImageStatus, nil
}

func (service AutomaticUpdateService) sendNotificationToApplicationOwner(
	application dto_service_resource.Application,
	emailNotification notification.OwnerNotification) error {
	log.Debugf("sending notification for the following environments %v", application.Environments)
	for _, environment := range application.Environments {
		log.Debugf("sending notification for the following environment %s", environment)
		log.Debugf("Getting owner email address of %s with labels %v", application.ApplicationName, application.Labels)
		owner, err := service.kubernetesResourceService.GetOwnerEmailAddressesByPodLabels(application.Labels, environment)

		if err != nil {
			log.Errorf("Failed to get owner email addresses of application %s: %v", application.ApplicationName, err)
			return err
		}

		log.Debugf("Generating email for owner %s", owner)
		log.Debugf("Creating html body for email")
		err = service.notificationService.SendEmail(notification_service.EmailNotification{
			Body:    emailNotification.GetBody(),
			Subject: emailNotification.GetSubject(),
			ToEmail: owner,
		})

		if err != nil {
			log.Errorf("Failed to send email to owner email address: %s of application %s: %v",
				owner,
				application.ApplicationName,
				err,
			)

			return err
		}
	}

	return nil
}

func (service AutomaticUpdateService) shouldApplicationResourceUpdate(
	applicationResource dto_service_resource.ApplicationResource) (dto_service_resource.Vulnerability, bool) {

	rules := []update_service_rules.IAutomaticUpdateRule{
		update_service_rules.NewCVSSUpdateRule(),
		update_service_rules.NewVulnerablePackageHasAFixVersionUpdateRule(),
	}

	for _, vulnerability := range applicationResource.Vulnerabilities {
		shouldUpdate := true
		for _, rule := range rules {
			log.Debugf("Checking if vulnerability %s matches rule %s", vulnerability.VulnerabilityName, rule)
			if !rule.ShouldUpdateImage(vulnerability) {
				log.Debugf("Rule %s did not match for current vulnerability %s in package %s",
					rule.GetName(),
					vulnerability.VulnerabilityName,
					vulnerability.PackageName)
				shouldUpdate = false
				continue
			}

			log.Infof("Rule %s did match for current vulnerability %s in package %s",
				rule.GetName(),
				vulnerability.VulnerabilityName,
				vulnerability.PackageName)
		}

		if shouldUpdate {
			return vulnerability, true
		}
	}

	return dto_service_resource.Vulnerability{}, false
}
