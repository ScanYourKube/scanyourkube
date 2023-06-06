package di

import (
	"net/http"

	http_client "github.com/scanyourkube/cronjob/api/http_client"
	api_keel_webhook "github.com/scanyourkube/cronjob/api/keel"
	kubeclarity_api "github.com/scanyourkube/cronjob/api/kubeclarity"
	kubernetes_api "github.com/scanyourkube/cronjob/api/kubernetes"
	resource_service_kubernetes "github.com/scanyourkube/cronjob/service/kubernetes"
	notification_service "github.com/scanyourkube/cronjob/service/notification"
	resource_service "github.com/scanyourkube/cronjob/service/resource"
	scanservice "github.com/scanyourkube/cronjob/service/scan"
	update_service "github.com/scanyourkube/cronjob/service/update"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func InitializeScanService(
	kubeclarityApiUrl string,
	kubernetesService resource_service_kubernetes.IKubernetesResourceService,
	vulnerabilitesToIgnore []string,
) scanservice.ScanService {
	return scanservice.NewScanService(kubeclarity_api.NewKubeClarityApi(
		http_client.NewMyHttpClient(&http.Client{}),
		kubeclarityApiUrl,
		vulnerabilitesToIgnore),
		kubernetesService)
}

func InitializeResourceService(kubeclarityApiUrl string, vulnerabilitesToIgnore []string) resource_service.ResourceService {
	return resource_service.NewResourceService(kubeclarity_api.NewKubeClarityApi(
		http_client.NewMyHttpClient(&http.Client{}),
		kubeclarityApiUrl,
		vulnerabilitesToIgnore))
}

func InitializeKeelUpdateService(webhookUrl string) update_service.KeelUpdateService {
	return update_service.NewKeelUpdateService(api_keel_webhook.NewKeelWebHookApi(&http.Client{}, webhookUrl))
}

func InitializeAutomaticUpdateService(
	keelUpdateService update_service.KeelUpdateService,
	scanService scanservice.ScanService,
	kubernetesService resource_service_kubernetes.IKubernetesResourceService,
	notificationService notification_service.IEmailNotificationService,
	kubeClarityResourceService resource_service.ResourceService) update_service.AutomaticUpdateService {
	return update_service.NewAutomaticUpdateService(keelUpdateService, scanService, kubernetesService, notificationService, kubeClarityResourceService)
}

func InizializeKubernetesService(namespacesToIgnore []string) resource_service_kubernetes.IKubernetesResourceService {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the kubernetes clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// creates the dynamic ClientSet
	dynamicClientSet, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return resource_service_kubernetes.NewKubernetesResourceService(kubernetes_api.NewKubernetesApi(clientSet, dynamicClientSet), namespacesToIgnore)
}

func InitializeNotificationService(
	senderEmailAddress string,
	smtpServer string,
	smtpPort string,
) notification_service.IEmailNotificationService {
	return notification_service.NewEmailNotificationService(senderEmailAddress, smtpServer, smtpPort)
}
