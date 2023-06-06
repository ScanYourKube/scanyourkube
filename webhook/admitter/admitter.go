// Package admission handles kubernetes admissions,
// it takes admission requests and returns admission reviews;
// for example, to mutate or validate deployments.
package admission

import (
	"errors"
	"strings"

	annotator "github.com/scanyourkube/webhook/annotator"
	labeler "github.com/scanyourkube/webhook/labeler"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
)

type IAdmitter interface {
	MutateAdmissionReview(Request *admissionv1.AdmissionRequest) (*admissionv1.AdmissionReview, error)
}

func Annotate(annotations map[string]string, annotators []annotator.IAnnotator) map[string]string {
	if annotations == nil {
		annotations = map[string]string{}
	}

	for _, annotator := range annotators {
		annotations = annotator.Annotate(annotations)
	}

	return annotations
}

func Label(labels map[string]string, labelers []labeler.ILabeler) map[string]string {
	if labels == nil {
		labels = map[string]string{}
	}

	for _, labeler := range labelers {
		labels = labeler.Label(labels)
	}

	return labels
}

func GetVersionFromImageName(name string) (*string, error) {
	parts := strings.Split(name, ":")
	if len(parts) > 1 {
		return &parts[1], nil
	}

	return nil, errors.New("missing version in image tag")
}

// reviewResponse TODO: godoc
func ReviewResponse(uid types.UID, allowed bool, httpCode int32,
	reason string) *admissionv1.AdmissionReview {
	return &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:     uid,
			Allowed: allowed,
			Result: &metav1.Status{
				Code:    httpCode,
				Message: reason,
			},
		},
	}
}

// patchReviewResponse builds an admission review with given json patch
func PatchReviewResponse(uid types.UID, patch []byte) (*admissionv1.AdmissionReview, error) {
	patchType := admissionv1.PatchTypeJSONPatch

	return &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:       uid,
			Allowed:   true,
			PatchType: &patchType,
			Patch:     patch,
		},
	}, nil
}
