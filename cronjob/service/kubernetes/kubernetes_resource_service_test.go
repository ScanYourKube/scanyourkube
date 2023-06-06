package kubernetes_resource_service

import (
	"errors"
	"reflect"

	kubernetes_api "github.com/scanyourkube/cronjob/testing/mocks/api/kubernetes"
	"testing"

	"github.com/golang/mock/gomock"
	rancherv3 "github.com/rancher/rancher/pkg/client/generated/management/v3"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetNamespaceNames(t *testing.T) {
	cases := []struct {
		input struct {
			ignoredNamespaces []string
			allNamespaces     []v1.Namespace
		}
		expected []string
	}{
		{
			input: struct {
				ignoredNamespaces []string
				allNamespaces     []v1.Namespace
			}{
				ignoredNamespaces: []string{"foo", "bar", "baz"},
				allNamespaces: []v1.Namespace{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "foo",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "bar",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "baz",
						},
					},
				},
			},
			expected: []string{},
		},
		{
			input: struct {
				ignoredNamespaces []string
				allNamespaces     []v1.Namespace
			}{
				ignoredNamespaces: []string{},
				allNamespaces: []v1.Namespace{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "foo",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "bar",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "baz",
						},
					},
				},
			},
			expected: []string{"foo", "bar", "baz"},
		},
		{
			input: struct {
				ignoredNamespaces []string
				allNamespaces     []v1.Namespace
			}{
				ignoredNamespaces: []string{"foo", "baz"},
				allNamespaces: []v1.Namespace{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "foo",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "bar",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "baz",
						},
					},
				},
			},
			expected: []string{"bar"},
		},
	}
	for _, c := range cases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockKubernetesApi := kubernetes_api.NewMockIKubernetesApi(mockCtrl)
		resourceService := NewKubernetesResourceService(mockKubernetesApi, c.input.ignoredNamespaces)
		mockKubernetesApi.EXPECT().GetNamespaces().Return(c.input.allNamespaces).Times(1)
		got := resourceService.GetNamespaceNames()

		if !reflect.DeepEqual(got, c.expected) {
			t.Errorf("removeIgnoredNamespaces(%v, %v) = %v, want %v",
				c.input.ignoredNamespaces,
				c.input.allNamespaces,
				got,
				c.expected,
			)
		}
	}
}

func TestGetOwnerEmailAddressOfDeployment(t *testing.T) {
	cases := []struct {
		input struct {
			annotations        map[string]string
			userAttribute      rancherv3.UserAttribute
			deploymentCalls    int
			userAttributeCalls int
		}
		expected    string
		expectedErr bool
	}{
		{
			input: struct {
				annotations        map[string]string
				userAttribute      rancherv3.UserAttribute
				deploymentCalls    int
				userAttributeCalls int
			}{
				annotations: map[string]string{
					"scanyourkube.io/owner": "foo",
				},

				userAttribute: rancherv3.UserAttribute{
					ExtraByProvider: map[string]map[string][]string{
						"activedirectory": {
							"username": []string{"test@test.com"},
						},
					},
				},
				deploymentCalls:    1,
				userAttributeCalls: 1,
			},
			expected:    "test@test.com",
			expectedErr: false,
		},
		{
			input: struct {
				annotations        map[string]string
				userAttribute      rancherv3.UserAttribute
				deploymentCalls    int
				userAttributeCalls int
			}{
				annotations: map[string]string{
					"scanyourkube.io/owner": "wrongUserName",
				},
				userAttribute: rancherv3.UserAttribute{},

				deploymentCalls:    1,
				userAttributeCalls: 1,
			},
			expected:    "",
			expectedErr: true,
		},
		{
			input: struct {
				annotations        map[string]string
				userAttribute      rancherv3.UserAttribute
				deploymentCalls    int
				userAttributeCalls int
			}{
				annotations: map[string]string{},
				userAttribute: rancherv3.UserAttribute{
					ExtraByProvider: map[string]map[string][]string{
						"activedirectory": {
							"username": []string{"test@test.com"},
						},
					},
				},

				deploymentCalls:    1,
				userAttributeCalls: 0,
			},
			expected:    "",
			expectedErr: true,
		},
	}
	for _, c := range cases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockKubernetesApi := kubernetes_api.NewMockIKubernetesApi(mockCtrl)
		resourceService := NewKubernetesResourceService(mockKubernetesApi, []string{})
		mockKubernetesApi.EXPECT().GetUserAttributes(gomock.Any()).Return(
			&c.input.userAttribute,
			nil,
		).Times(c.input.userAttributeCalls)
		email, err := resourceService.GetOwnerEmailAddressFromAnnotations(c.input.annotations)
		if err != nil && !c.expectedErr {
			t.Errorf("GetOwnerEmailAddressFromAnnotations(%v) threw error %v", c.input.annotations, err)
		}
		if !c.expectedErr {
			got := email
			if !reflect.DeepEqual(got, c.expected) {
				t.Errorf("GetOwnerEmailAddressFromAnnotations(%v) = %v, want %v",
					c.input.annotations,
					got,
					c.expected,
				)
			}
		}

		if c.expectedErr && email != "" {
			got := email
			t.Errorf("GetOwnerEmailAddressFromAnnotations(%v) = %v, want %v",
				c.input.annotations,
				got,
				"nil",
			)
		}

	}
}

func TestGetOwnerEmailAddressesByPodLabels(t *testing.T) {
	cases := []struct {
		input struct {
			labels             []string
			userAttribute      rancherv3.UserAttribute
			deploymentCalls    int
			userAttributeCalls int
		}
		expected    string
		expectedErr bool
	}{
		{
			input: struct {
				labels             []string
				userAttribute      rancherv3.UserAttribute
				deploymentCalls    int
				userAttributeCalls int
			}{
				labels: []string{
					"scanyourkube.io/owner=foo",
				},

				userAttribute: rancherv3.UserAttribute{
					ExtraByProvider: map[string]map[string][]string{
						"activedirectory": {
							"username": []string{"test@test.com"},
						},
					},
				},
				deploymentCalls:    1,
				userAttributeCalls: 1,
			},
			expected:    "test@test.com",
			expectedErr: false,
		},
		{
			input: struct {
				labels             []string
				userAttribute      rancherv3.UserAttribute
				deploymentCalls    int
				userAttributeCalls int
			}{
				labels: []string{
					"scanyourkube.io/owner=wrongUserName",
				},
				userAttribute: rancherv3.UserAttribute{},

				deploymentCalls:    1,
				userAttributeCalls: 1,
			},
			expected:    "",
			expectedErr: true,
		},
		{
			input: struct {
				labels             []string
				userAttribute      rancherv3.UserAttribute
				deploymentCalls    int
				userAttributeCalls int
			}{
				labels: []string{},
				userAttribute: rancherv3.UserAttribute{
					ExtraByProvider: map[string]map[string][]string{
						"activedirectory": {
							"username": []string{"test@test.com"},
						},
					},
				},

				deploymentCalls:    1,
				userAttributeCalls: 0,
			},
			expected:    "",
			expectedErr: true,
		},
	}
	for _, c := range cases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockKubernetesApi := kubernetes_api.NewMockIKubernetesApi(mockCtrl)
		resourceService := NewKubernetesResourceService(mockKubernetesApi, []string{})
		mockKubernetesApi.EXPECT().GetUserAttributes(gomock.Any()).Return(
			&c.input.userAttribute,
			nil,
		).Times(c.input.userAttributeCalls)
		email, err := resourceService.GetOwnerEmailAddressesByPodLabels(c.input.labels, "DUMMY")
		if err != nil && !c.expectedErr {
			t.Errorf("GetOwnerEmailAddressFromAnnotations(%v) threw error %v", c.input.labels, err)
		}
		if !c.expectedErr {
			got := email
			if !reflect.DeepEqual(got, c.expected) {
				t.Errorf("GetOwnerEmailAddressFromAnnotations(%v) = %v, want %v",
					c.input.labels,
					got,
					c.expected,
				)
			}
		}

		if c.expectedErr && email != "" {
			got := email
			t.Errorf("GetOwnerEmailAddressFromAnnotations(%v) = %v, want %v",
				c.input.labels,
				got,
				"nil",
			)
		}

	}
}

func TestGetOwnerEmailAddressesByPodLabelsFromPod(t *testing.T) {
	cases := []struct {
		input struct {
			labels             []string
			userAttribute      rancherv3.UserAttribute
			deploymentCalls    int
			userAttributeCalls int
			deploymentName     string
			annotations        map[string]string
		}
		expected    string
		expectedErr bool
	}{
		{
			input: struct {
				labels             []string
				userAttribute      rancherv3.UserAttribute
				deploymentCalls    int
				userAttributeCalls int
				deploymentName     string
				annotations        map[string]string
			}{
				labels: []string{
					"scanyourkube.io/podName=foo",
				},

				userAttribute: rancherv3.UserAttribute{
					ExtraByProvider: map[string]map[string][]string{
						"activedirectory": {
							"username": []string{"test@test.com"},
						},
					},
				},
				deploymentCalls:    1,
				userAttributeCalls: 1,
				deploymentName:     "foo",
				annotations:        map[string]string{"scanyourkube.io/owner": "foo"},
			},
			expected:    "test@test.com",
			expectedErr: false,
		},
		{
			input: struct {
				labels             []string
				userAttribute      rancherv3.UserAttribute
				deploymentCalls    int
				userAttributeCalls int
				deploymentName     string
				annotations        map[string]string
			}{
				labels: []string{
					"scanyourkube.io/podName=wrongUserName",
				},
				userAttribute: rancherv3.UserAttribute{},

				deploymentCalls:    1,
				userAttributeCalls: 1,
				deploymentName:     "foo",
				annotations:        map[string]string{"scanyourkube.io/owner": "foo"},
			},
			expected:    "",
			expectedErr: true,
		},
		{
			input: struct {
				labels             []string
				userAttribute      rancherv3.UserAttribute
				deploymentCalls    int
				userAttributeCalls int
				deploymentName     string
				annotations        map[string]string
			}{
				labels: []string{"scanyourkube.io/podName=foo"},
				userAttribute: rancherv3.UserAttribute{
					ExtraByProvider: map[string]map[string][]string{
						"activedirectory": {
							"username": []string{"test@test.com"},
						},
					},
				},

				deploymentCalls:    1,
				userAttributeCalls: 0,
				deploymentName:     "foo",
				annotations:        map[string]string{"scanyourkube.io/owner": ""},
			},
			expected:    "",
			expectedErr: true,
		},
	}
	for _, c := range cases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockKubernetesApi := kubernetes_api.NewMockIKubernetesApi(mockCtrl)
		resourceService := NewKubernetesResourceService(mockKubernetesApi, []string{})
		mockKubernetesApi.EXPECT().GetPodByLabel(gomock.Any(), gomock.Any()).Return(
			&v1.Pod{ObjectMeta: metav1.ObjectMeta{OwnerReferences: []metav1.OwnerReference{{Name: c.input.deploymentName, Kind: "ReplicaSet"}}}},
			nil,
		).Times(c.input.deploymentCalls)
		mockKubernetesApi.EXPECT().GetReplicaSet(c.input.deploymentName, gomock.Any()).Return(
			&appsv1.ReplicaSet{ObjectMeta: metav1.ObjectMeta{OwnerReferences: []metav1.OwnerReference{{Name: c.input.deploymentName, Kind: "Deployment"}}}},
			nil,
		).Times(c.input.deploymentCalls)
		mockKubernetesApi.EXPECT().GetDeployment(c.input.deploymentName, gomock.Any()).Return(
			&appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Annotations: c.input.annotations}},
			nil,
		).Times(c.input.deploymentCalls)
		mockKubernetesApi.EXPECT().GetUserAttributes(gomock.Any()).Return(
			&c.input.userAttribute,
			nil,
		).Times(c.input.userAttributeCalls)
		email, err := resourceService.GetOwnerEmailAddressesByPodLabels(c.input.labels, "DUMMY")
		if err != nil && !c.expectedErr {
			t.Errorf("GetOwnerEmailAddressFromAnnotations(%v) threw error %v", c.input.labels, err)
		}
		if !c.expectedErr {
			got := email
			if !reflect.DeepEqual(got, c.expected) {
				t.Errorf("GetOwnerEmailAddressFromAnnotations(%v) = %v, want %v",
					c.input.labels,
					got,
					c.expected,
				)
			}
		}

		if c.expectedErr && email != "" {
			got := email
			t.Errorf("GetOwnerEmailAddressFromAnnotations(%v) = %v, want %v",
				c.input.labels,
				got,
				"nil",
			)
		}

	}
}

func TestGetOwnerEmailAddressesByPodLabelsFromPodInStatefulSet(t *testing.T) {
	cases := []struct {
		input struct {
			labels             []string
			userAttribute      rancherv3.UserAttribute
			statefulSetCalls   int
			userAttributeCalls int
			statefulSetName    string
			annotations        map[string]string
		}
		expected    string
		expectedErr bool
	}{
		{
			input: struct {
				labels             []string
				userAttribute      rancherv3.UserAttribute
				statefulSetCalls   int
				userAttributeCalls int
				statefulSetName    string
				annotations        map[string]string
			}{
				labels: []string{
					"scanyourkube.io/podName=foo",
				},

				userAttribute: rancherv3.UserAttribute{
					ExtraByProvider: map[string]map[string][]string{
						"activedirectory": {
							"username": []string{"test@test.com"},
						},
					},
				},
				statefulSetCalls:   1,
				userAttributeCalls: 1,
				statefulSetName:    "foo",
				annotations:        map[string]string{"scanyourkube.io/owner": "foo"},
			},
			expected:    "test@test.com",
			expectedErr: false,
		},
		{
			input: struct {
				labels             []string
				userAttribute      rancherv3.UserAttribute
				statefulSetCalls   int
				userAttributeCalls int
				statefulSetName    string
				annotations        map[string]string
			}{
				labels: []string{
					"scanyourkube.io/podName=wrongUserName",
				},
				userAttribute: rancherv3.UserAttribute{},

				userAttributeCalls: 1,
				statefulSetCalls:   1,
				statefulSetName:    "foo",
				annotations:        map[string]string{"scanyourkube.io/owner": "foo"},
			},
			expected:    "",
			expectedErr: true,
		},
		{
			input: struct {
				labels             []string
				userAttribute      rancherv3.UserAttribute
				statefulSetCalls   int
				userAttributeCalls int
				statefulSetName    string
				annotations        map[string]string
			}{
				labels: []string{"scanyourkube.io/podName=foo"},
				userAttribute: rancherv3.UserAttribute{
					ExtraByProvider: map[string]map[string][]string{
						"activedirectory": {
							"username": []string{"test@test.com"},
						},
					},
				},

				statefulSetCalls:   1,
				userAttributeCalls: 0,
				statefulSetName:    "foo",
				annotations:        map[string]string{"scanyourkube.io/owner": ""},
			},
			expected:    "",
			expectedErr: true,
		},
	}
	for _, c := range cases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockKubernetesApi := kubernetes_api.NewMockIKubernetesApi(mockCtrl)
		resourceService := NewKubernetesResourceService(mockKubernetesApi, []string{})
		mockKubernetesApi.EXPECT().GetPodByLabel(gomock.Any(), gomock.Any()).Return(
			&v1.Pod{ObjectMeta: metav1.ObjectMeta{OwnerReferences: []metav1.OwnerReference{{Name: c.input.statefulSetName, Kind: "StatefulSet"}}}},
			nil,
		).Times(c.input.statefulSetCalls)
		mockKubernetesApi.EXPECT().GetStatefulSet(c.input.statefulSetName, gomock.Any()).Return(
			&appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Annotations: c.input.annotations}},
			nil,
		).Times(c.input.statefulSetCalls)
		mockKubernetesApi.EXPECT().GetUserAttributes(gomock.Any()).Return(
			&c.input.userAttribute,
			nil,
		).Times(c.input.userAttributeCalls)
		email, err := resourceService.GetOwnerEmailAddressesByPodLabels(c.input.labels, "DUMMY")
		if err != nil && !c.expectedErr {
			t.Errorf("GetOwnerEmailAddressFromAnnotations(%v) threw error %v", c.input.labels, err)
		}
		if !c.expectedErr {
			got := email
			if !reflect.DeepEqual(got, c.expected) {
				t.Errorf("GetOwnerEmailAddressFromAnnotations(%v) = %v, want %v",
					c.input.labels,
					got,
					c.expected,
				)
			}
		}

		if c.expectedErr && email != "" {
			got := email
			t.Errorf("GetOwnerEmailAddressFromAnnotations(%v) = %v, want %v",
				c.input.labels,
				got,
				"nil",
			)
		}

	}
}

func TestFailingGetOwnerEmailAddressFromPod(t *testing.T) {
	cases := []struct {
		input struct {
			labels            []string
			failPod           bool
			failDeployment    bool
			failStatefulSet   bool
			podDeploymentType string
		}
		expectedErr bool
	}{
		{
			input: struct {
				labels            []string
				failPod           bool
				failDeployment    bool
				failStatefulSet   bool
				podDeploymentType string
			}{
				labels:            []string{"scanyourkube.io/podName=foo"},
				failPod:           true,
				failDeployment:    false,
				failStatefulSet:   false,
				podDeploymentType: "",
			},
			expectedErr: true,
		},
		{
			input: struct {
				labels            []string
				failPod           bool
				failDeployment    bool
				failStatefulSet   bool
				podDeploymentType string
			}{
				labels:            []string{"scanyourkube.io/podName=foo"},
				failPod:           false,
				failDeployment:    true,
				failStatefulSet:   false,
				podDeploymentType: "Deployment",
			},
			expectedErr: true,
		},
		{
			input: struct {
				labels            []string
				failPod           bool
				failDeployment    bool
				failStatefulSet   bool
				podDeploymentType string
			}{
				labels:            []string{"scanyourkube.io/podName=foo"},
				failPod:           false,
				failDeployment:    false,
				failStatefulSet:   true,
				podDeploymentType: "StatefulSet",
			},
			expectedErr: true,
		},
	}
	for _, c := range cases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockKubernetesApi := kubernetes_api.NewMockIKubernetesApi(mockCtrl)
		resourceService := NewKubernetesResourceService(mockKubernetesApi, []string{})
		if c.input.failPod {
			mockKubernetesApi.EXPECT().GetPodByLabel(gomock.Any(), gomock.Any()).Return(
				nil,
				errors.New("error in pod"),
			).AnyTimes()
		} else {
			mockKubernetesApi.EXPECT().GetPodByLabel(gomock.Any(), gomock.Any()).Return(
				&v1.Pod{ObjectMeta: metav1.ObjectMeta{OwnerReferences: []metav1.OwnerReference{{Name: "foo", Kind: c.input.podDeploymentType}}}},
				nil,
			).AnyTimes()
		}

		if c.input.failDeployment {
			mockKubernetesApi.EXPECT().GetDeployment("foo", gomock.Any()).Return(
				nil,
				errors.New("error in statefulset"),
			).AnyTimes()
		} else {
			mockKubernetesApi.EXPECT().GetDeployment("foo", gomock.Any()).Return(
				&appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"scanyourkube.io/owner": "foo"}}},
				nil,
			).AnyTimes()
		}

		if c.input.failStatefulSet {
			mockKubernetesApi.EXPECT().GetStatefulSet("foo", gomock.Any()).Return(
				nil,
				errors.New("error in statefulset"),
			).AnyTimes()
		} else {
			mockKubernetesApi.EXPECT().GetStatefulSet("foo", gomock.Any()).Return(
				&appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"scanyourkube.io/owner": "foo"}}},
				nil,
			).AnyTimes()
		}
		mockKubernetesApi.EXPECT().GetUserAttributes(gomock.Any()).Times(0)
		email, err := resourceService.GetOwnerEmailAddressesByPodLabels(c.input.labels, "DUMMY")
		if err != nil && !c.expectedErr {
			t.Errorf("GetOwnerEmailAddressFromAnnotations(%v) threw error %v", c.input.labels, err)
		}

		if c.expectedErr && email != "" {
			got := email
			t.Errorf("GetOwnerEmailAddressFromAnnotations(%v) = %v, want %v",
				c.input.labels,
				got,
				"nil",
			)
		}

	}
}
