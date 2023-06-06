package kubernetes_resource_service

import (
	"errors"
	"fmt"
	"strings"

	api_kubernetes "github.com/scanyourkube/cronjob/api/kubernetes"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/strings/slices"
)

type IKubernetesResourceService interface {
	GetPods(namespace string) []v1.Pod
	GetNamespaceNames() []string
	GetOwnerEmailAddressFromAnnotations(annotations map[string]string) (string, error)
	GetOwnerEmailAddressesByPodLabels(labels []string, environment string) (string, error)
}

type KubernetesResourceService struct {
	KubernetesApi      api_kubernetes.IKubernetesApi
	namespacesToIgnore []string
}

func NewKubernetesResourceService(
	KubernetesApi api_kubernetes.IKubernetesApi,
	namespacesToIgnore []string,
) KubernetesResourceService {
	return KubernetesResourceService{
		KubernetesApi:      KubernetesApi,
		namespacesToIgnore: namespacesToIgnore}
}

func (s KubernetesResourceService) GetPods(namespace string) []v1.Pod {
	pods, err := s.KubernetesApi.GetPods(namespace)
	if err != nil {
		return nil
	}

	return pods
}

func (s KubernetesResourceService) GetDeploymentOfReplicaSet(
	name string,
	namespace string,
) (*appsv1.Deployment, error) {
	replicaSetLog := log.WithFields(log.Fields{"replicaSetName": name, "method": "GetDeploymentOfReplicaSet"})
	replicaSet, err := s.KubernetesApi.GetReplicaSet(name, namespace)
	if err != nil {
		replicaSetLog.Errorf("Couldn't get replica set with error: %v", err)
		return nil, err
	}
	deploymentReference := &metav1.OwnerReference{}
	for _, ownerReference := range replicaSet.OwnerReferences {
		if ownerReference.DeepCopy().Kind == "Deployment" {
			replicaSetLog.Debugf("Found owner reference in replica set %s for deployment %v", replicaSet.Name, ownerReference)
			deploymentReference = ownerReference.DeepCopy()
		}
	}

	if deploymentReference.Name == "" {
		replicaSetLog.Errorf("Missing owner reference for deployment in replica set %s", replicaSet.Name)
		return nil, errors.New("missing owner reference for deployment in replica set")
	}

	deployment, err := s.KubernetesApi.GetDeployment(deploymentReference.Name, namespace)
	if err != nil {
		replicaSetLog.Errorf("Couldn't get deployment with error: %v", err)
		return nil, err
	}
	replicaSetLog.Debugf("Got deployment with name: %s", deployment.Name)
	return deployment, nil
}

func (s KubernetesResourceService) GetOwnerEmailAddressFromAnnotations(
	annotations map[string]string) (string, error) {
	if ownerId, ok := annotations["scanyourkube.io/owner"]; ok {
		if ownerId == "" {
			return "", errors.New("couldn't find owner email address")
		}
		return s.GetOwnerEmailAddressFromUserId(ownerId)
	}

	return "", errors.New("couldn't find owner email address")
}

func (s KubernetesResourceService) GetOwnerEmailAddressFromUserId(
	userId string) (string, error) {
	ownerAttributes, err := s.KubernetesApi.GetUserAttributes(userId)
	if err != nil {
		return "", err
	}

	if ownerEmailAddresses, ok := ownerAttributes.ExtraByProvider["activedirectory"]["username"]; ok {
		if len(ownerEmailAddresses) > 0 {
			log.Debugf("Got owner email address: %s", ownerEmailAddresses[0])
			return ownerEmailAddresses[0], nil
		}
	}
	return "", errors.New("couldn't find owner email address")
}

func (s KubernetesResourceService) GetOwnerEmailAddressesByPodLabels(
	labels []string,
	environment string,
) (string, error) {
	ownerLabel, _ := s.getKubernetesLabelByName(labels, "scanyourkube.io/owner")
	if ownerLabel != "" {

		labelSplit := strings.Split(ownerLabel, "=")
		if len(labelSplit) == 1 {
			return "", fmt.Errorf("error while splitting ownerLabel %s with seperator '='", ownerLabel)
		}
		return s.GetOwnerEmailAddressFromUserId(labelSplit[1])
	}

	labelAppName, err := s.getKubernetesLabelByName(labels, api_kubernetes.KubernetesAppNameLabel)
	if err != nil {

		return "", err
	}

	podOwnerReference, err := s.getOwnerReferenceOfPod(labelAppName, environment)
	if err != nil {
		return "", err
	}

	return s.GetOwnerFromPodOwnerReference(*podOwnerReference, environment)
}

func (s KubernetesResourceService) GetOwnerFromPodOwnerReference(podOwnerReference metav1.OwnerReference, environment string) (string, error) {

	switch podOwnerReference.Kind {
	case "ReplicaSet":

		deployment, err := s.GetDeploymentOfReplicaSet(podOwnerReference.Name, environment)
		if err != nil {
			log.Errorf("Couldn't get deployment with error: %v", err)
			return "", err
		}

		log.Debugf("Was able to receive deployment of application %s", deployment.Name)
		owner, err := s.GetOwnerEmailAddressFromAnnotations(deployment.Annotations)

		if err != nil {
			log.Errorf("Failed to get owner email address of deployment %s: %v", deployment.Name, err)
			return "", err
		}

		return owner, nil

	case "StatefulSet":
		statefulSet, err := s.KubernetesApi.GetStatefulSet(podOwnerReference.Name, environment)
		if err != nil {
			log.Errorf("Couldn't get statefulset with error: %v", err)
			return "", err
		}

		owner, err := s.GetOwnerEmailAddressFromAnnotations(statefulSet.Annotations)
		if err != nil {
			log.Errorf("Failed to get owner email address of statefulSet %s: %v", statefulSet.Name, err)
			return "", err
		}
		return owner, nil

	default:
		log.Errorf("Not implemented case for Kind %s", podOwnerReference.Kind)
		return "", fmt.Errorf("not implemented case for Kind %s", podOwnerReference.Kind)
	}

}

func (s KubernetesResourceService) GetNamespaceNames() []string {
	activeNamespaces := s.KubernetesApi.GetNamespaces()
	namespaces := make([]string, 0)
	for _, namespace := range activeNamespaces {
		if !slices.Contains(s.namespacesToIgnore, namespace.Name) {
			namespaces = append(namespaces, namespace.Name)
		}
	}
	return namespaces
}

func (KubernetesResourceService) getKubernetesLabelByName(labels []string, labelName string) (string, error) {
	for _, label := range labels {
		if strings.Contains(label, labelName) {
			return label, nil
		}
	}

	return "", fmt.Errorf("could not find label %s in labels %v", labelName, labels)
}

func (KubernetesResourceService) getGenericOwnerReference(
	podLog *log.Entry,
	podReferences []metav1.OwnerReference) (*metav1.OwnerReference, error) {
	for _, podOwnerReference := range podReferences {
		switch podOwnerReference.Kind {
		case "ReplicaSet", "StatefulSet":
			podLog.Debugf("Found owner reference in pod for %s %v", podOwnerReference.Kind, podOwnerReference)
			return podOwnerReference.DeepCopy(), nil
		default:
			continue
		}
	}

	return nil, errors.New("missing owner reference for replica set in pod")
}

func (s KubernetesResourceService) getOwnerReferenceOfPod(
	podLabel string,
	namespace string,
) (*metav1.OwnerReference, error) {
	podLog := log.WithFields(log.Fields{"podLabel": podLabel, "method": "getOwnerReferenceOfPod"})
	pod, err := s.KubernetesApi.GetPodByLabel(podLabel, namespace)
	if err != nil {
		podLog.Errorf("Couldn't get pod with error: %v", err)
		return nil, err
	}

	podLog.Debugf("Got pod with name: %s in namespace: %s", pod.Name, pod.Namespace)

	ownerReference, err := s.getGenericOwnerReference(podLog, pod.OwnerReferences)

	if err != nil {
		podLog.Errorf("Missing owner reference for replica set in pod %s", pod.Name)
		return nil, err
	}

	return ownerReference, nil
}
