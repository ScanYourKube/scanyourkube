// Package admission handles kubernetes admissions,
// it takes admission requests and returns admission reviews;
// for example, to mutate or validate deployments.
package admission

import (
	"encoding/json"
	"fmt"
	"net/http"

	annotator "github.com/scanyourkube/webhook/annotator"
	labeler "github.com/scanyourkube/webhook/labeler"

	"github.com/wI2L/jsondiff"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
)

// Admitter is a container for admission business
type StatefulSetAdmitter struct {
}

// MutateAdmissionReview takes an admission request and mutates the statefulset within,
// it returns an admission review with mutations as a json patch (if any)
func (a StatefulSetAdmitter) MutateAdmissionReview(
	request *admissionv1.AdmissionRequest,
) (*admissionv1.AdmissionReview, error) {

	statefulSet, err := a.statefulSet(request)
	if err != nil {
		e := fmt.Sprintf("could not parse statefulset in admission review request: %v", err)
		return ReviewResponse(request.UID, false, http.StatusBadRequest, e), err
	}
	mutatedStatefulset := statefulSet.DeepCopy()

	imageVersion, err := GetVersionFromImageName(statefulSet.Spec.Template.Spec.Containers[0].Image)

	if err != nil {
		e := fmt.Sprintf("could not mutate deployment: %v", err)
		return ReviewResponse(request.UID, false, http.StatusBadRequest, e), err
	}

	annotators := []annotator.IAnnotator{
		annotator.KeelPolicyAnnotator{
			ImageTag: *imageVersion,
		},
	}

	mutatedStatefulset.Spec.Template.Annotations = Annotate(mutatedStatefulset.Annotations, annotators)

	podLabelers := []labeler.ILabeler{
		labeler.ScanYourKubeOwnerLabeler{
			Owner:     request.UserInfo.Username,
			Operation: request.Operation,
		},
		labeler.ScanYourKubePodNameLabeler{
			PodName: statefulSet.Namespace + "-" + statefulSet.Name,
		},
	}

	mutatedStatefulset.Spec.Template.Labels = Label(mutatedStatefulset.Spec.Template.Labels, podLabelers)

	// generate json patch
	patch, err := jsondiff.Compare(statefulSet, mutatedStatefulset)
	if err != nil {
		e := fmt.Sprintf("could not mutate statefulset: %v", err)
		return ReviewResponse(request.UID, false, http.StatusBadRequest, e), err
	}

	patchb, err := json.Marshal(patch)
	if err != nil {
		e := fmt.Sprintf("could not mutate statefulset: %v", err)
		return ReviewResponse(request.UID, false, http.StatusBadRequest, e), err
	}

	return PatchReviewResponse(request.UID, patchb)
}

// statefulSet extracts a StatefulSet of an admission request
func (a StatefulSetAdmitter) statefulSet(request *admissionv1.AdmissionRequest) (*appsv1.StatefulSet, error) {
	if request.Kind.Kind != "StatefulSet" {
		return nil, fmt.Errorf("only statefulset are supported here")
	}

	d := appsv1.StatefulSet{}
	if err := json.Unmarshal(request.Object.Raw, &d); err != nil {
		return nil, err
	}

	return &d, nil
}
