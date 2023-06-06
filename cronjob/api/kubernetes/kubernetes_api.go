package api_kubernetes

import (
	"context"
	"encoding/json"
	"fmt"

	rancherv3 "github.com/rancher/rancher/pkg/client/generated/management/v3"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type IKubernetesApi interface {
	GetNamespaces() []v1.Namespace
	GetPods(namespace string) ([]v1.Pod, error)
	GetPodByLabel(label string, namespace string) (*v1.Pod, error)
	GetDeployment(deploymentName string, namespace string) (*appsv1.Deployment, error)
	GetStatefulSet(statefulSetName string, namespace string) (*appsv1.StatefulSet, error)
	GetUserAttributes(userName string) (*rancherv3.UserAttribute, error)
	GetReplicaSet(replicaSetName string, namespace string) (*appsv1.ReplicaSet, error)
}

type KubernetesApi struct {
	ClientSet        kubernetes.Interface
	DynamicClientSet dynamic.Interface
}

func NewKubernetesApi(clientSet *kubernetes.Clientset, dynamicClientSet *dynamic.DynamicClient) IKubernetesApi {
	return KubernetesApi{
		ClientSet:        clientSet,
		DynamicClientSet: dynamicClientSet,
	}
}

const KubernetesAppNameLabel = "scanyourkube.io/podName"

func (a KubernetesApi) GetNamespaces() []v1.Namespace {

	namespaces, err := a.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	log.Debugf("There are %d namespaces in the cluster\n", len(namespaces.Items))

	return namespaces.Items
}

func (a KubernetesApi) GetPods(namespace string) ([]v1.Pod, error) {

	pods, err := a.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return make([]v1.Pod, 0), err
	}
	log.Debugf("There are %d pods in the namespace %s\n", len(pods.Items), namespace)

	return pods.Items, nil
}

func (a KubernetesApi) GetPodByLabel(label string, namespace string) (*v1.Pod, error) {

	pods, err := a.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}

	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("pod with label %s not found in namespace %s", label, namespace)
	}

	pod := &pods.Items[0]
	log.Debugf("Found %s pod in namespace %s\n", pod.Name, pod.Namespace)

	return pod, nil
}

func (a KubernetesApi) GetDeployment(deploymentName string, namespace string) (*appsv1.Deployment, error) {

	deployment, err := a.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	log.Debugf("Found %s deployment in namespace %s\n", deployment.Name, deployment.Namespace)

	return deployment, nil
}

func (a KubernetesApi) GetStatefulSet(statefulSetName string, namespace string) (*appsv1.StatefulSet, error) {
	statefulSet, err := a.ClientSet.AppsV1().StatefulSets(namespace).Get(
		context.TODO(),
		statefulSetName,
		metav1.GetOptions{},
	)
	if err != nil {
		return nil, err
	}
	log.Debugf("Found %s statefulSet in namespace %s\n", statefulSet.Name, statefulSet.Namespace)

	return statefulSet, nil
}

func (a KubernetesApi) GetReplicaSet(replicaSetName string, namespace string) (*appsv1.ReplicaSet, error) {
	replicaSet, err := a.ClientSet.AppsV1().ReplicaSets(namespace).Get(context.TODO(), replicaSetName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	log.Debugf("Found %s replicaSet in namespace %s\n", replicaSet.Name, replicaSet.Namespace)

	return replicaSet, nil
}

func (a KubernetesApi) GetUserAttributes(userName string) (*rancherv3.UserAttribute, error) {

	// Get resource interface for user attribute
	userAttributeGVR := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "userattributes",
	}
	attributeClient := a.DynamicClientSet.Resource(userAttributeGVR)

	unstructuredObject, err := attributeClient.Get(context.TODO(), userName, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	unstructured := unstructuredObject.UnstructuredContent()

	var userAttribute rancherv3.UserAttribute
	unstructuredjson, err := json.Marshal(unstructured)

	if err != nil {
		log.Errorf("error marshalling unstructured object: %v", err)
		return nil, err
	}
	log.Debugf("The json of the unstructured object is: %s", string(unstructuredjson))

	// Convert the jsonObject into userAttribute struct
	err = json.Unmarshal(unstructuredjson, &userAttribute)
	if err != nil {
		log.Debugf("error %s, converting unstructured to UserAttribute type with json", err.Error())
	}

	log.Debugf("Found UserAttribute for user %v with username %s \n", userAttribute, userName)
	log.Infof("Found user with username %s", userName)
	return &userAttribute, nil
}
