package api_kubernetes

import (
	"testing"

	fakedynamic "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestKubernetesApi_GetNamespaces(t *testing.T) {
	namespaces := []v1.Namespace{{
		ObjectMeta: metav1.ObjectMeta{
			Name: "DUMMY",
		},
	}}
	clientset := fake.NewSimpleClientset(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "DUMMY",
		},
	})

	api := KubernetesApi{
		ClientSet: clientset,
	}

	result := api.GetNamespaces()
	assert.NotNil(t, result)
	assert.Equal(t, namespaces, result)
}

func TestKubernetesApi_GetPods(t *testing.T) {
	pods := []v1.Pod{{

		ObjectMeta: metav1.ObjectMeta{
			Name:      "DUMMY",
			Namespace: "DUMMY",
		},
	}}
	clientset := fake.NewSimpleClientset(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "DUMMY",
			Namespace: "DUMMY",
		},
	})

	api := KubernetesApi{
		ClientSet: clientset,
	}

	result, err := api.GetPods("DUMMY")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, pods, result)
}

func TestKubernetesApiGetPodsByLabel(t *testing.T) {
	pod := &v1.Pod{

		ObjectMeta: metav1.ObjectMeta{
			Name:      "DUMMY",
			Namespace: "DUMMY",
			Labels: map[string]string{
				"scanyourkube.io/podName": "DUMMY",
			},
		},
	}
	clientset := fake.NewSimpleClientset(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "DUMMY",
			Namespace: "DUMMY",
			Labels: map[string]string{
				"scanyourkube.io/podName": "DUMMY",
			},
		},
	})

	api := KubernetesApi{
		ClientSet: clientset,
	}

	result, err := api.GetPodByLabel("scanyourkube.io/podName", "DUMMY")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, pod, result)
}

func TestKubernetesApiGetNoPodsByWrongLabel(t *testing.T) {
	clientset := fake.NewSimpleClientset(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "DUMMY",
			Namespace: "DUMMY",
			Labels: map[string]string{
				"NAME": "VALUE",
			},
		},
	})

	api := KubernetesApi{
		ClientSet: clientset,
	}

	result, err := api.GetPodByLabel("DUMMY", "DUMMY")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestKubernetesApi_GetDeployments(t *testing.T) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "DUMMY",
			Namespace: "DUMMY",
		},
	}
	clientset := fake.NewSimpleClientset(deployment)

	api := KubernetesApi{
		ClientSet: clientset,
	}

	result, err := api.GetDeployment("DUMMY", "DUMMY")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, deployment, result)
}

func TestKubernetesApi_GetStatefulSet(t *testing.T) {
	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "DUMMY",
			Namespace: "DUMMY",
		},
	}
	clientset := fake.NewSimpleClientset(statefulset)

	api := KubernetesApi{
		ClientSet: clientset,
	}

	result, err := api.GetStatefulSet("DUMMY", "DUMMY")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, statefulset, result)
}

func TestKubernetesApi_GetReplicaSet(t *testing.T) {
	replicaSet := &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "DUMMY",
			Namespace: "DUMMY",
		},
	}
	clientset := fake.NewSimpleClientset(replicaSet)

	api := KubernetesApi{
		ClientSet: clientset,
	}

	result, err := api.GetReplicaSet("DUMMY", "DUMMY")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, replicaSet, result)
}

func newUnstructured(apiVersion, kind, name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": apiVersion,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name": name,
			},
		},
	}
}

func TestKubernetesApi_GetUserAttributes(t *testing.T) {

	scheme := runtime.NewScheme()

	clientset := fakedynamic.NewSimpleDynamicClient(scheme, newUnstructured("management.cattle.io/v3", "UserAttribute", "DUMMY"))

	api := KubernetesApi{
		DynamicClientSet: clientset,
	}

	result, err := api.GetUserAttributes("DUMMY")
	assert.Nil(t, err)
	assert.NotNil(t, result)
}
