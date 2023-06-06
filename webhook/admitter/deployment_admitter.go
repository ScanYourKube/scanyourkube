// Package admission handles kubernetes admissions,
// it takes admission requests and returns admission reviews;
// for example, to mutate or validate deployments.
package admission

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wI2L/jsondiff"

	annotator "github.com/scanyourkube/webhook/annotator"
	labeler "github.com/scanyourkube/webhook/labeler"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
)

// DeploymentAdmitter is a container for admission business
type DeploymentAdmitter struct {
}

// MutateAdmissionReview takes an admission request and mutates the deployment within,
// it returns an admission review with mutations as a json patch (if any)
func (a DeploymentAdmitter) MutateAdmissionReview(
	request *admissionv1.AdmissionRequest,
) (*admissionv1.AdmissionReview, error) {

	deployment, err := a.deployment(request)
	if err != nil {
		e := fmt.Sprintf("could not parse deployment in admission review request: %v", err)
		return ReviewResponse(request.UID, false, http.StatusBadRequest, e), err
	}
	mutatedDeployment := deployment.DeepCopy()

	imageVersion, err := GetVersionFromImageName(deployment.Spec.Template.Spec.Containers[0].Image)

	if err != nil {
		e := fmt.Sprintf("could not mutate deployment: %v", err)
		return ReviewResponse(request.UID, false, http.StatusBadRequest, e), err
	}

	annotators := []annotator.IAnnotator{
		annotator.KeelPolicyAnnotator{
			ImageTag: *imageVersion,
		},
	}

	mutatedDeployment.Spec.Template.Annotations = Annotate(mutatedDeployment.Annotations, annotators)

	podLabelers := []labeler.ILabeler{
		labeler.ScanYourKubeOwnerLabeler{
			Owner:     request.UserInfo.Username,
			Operation: request.Operation,
		},
		labeler.ScanYourKubePodNameLabeler{
			PodName: deployment.Namespace + "-" + deployment.Name,
		},
	}

	mutatedDeployment.Spec.Template.Labels = Label(mutatedDeployment.Spec.Template.Labels, podLabelers)
	// generate json patch
	patch, err := jsondiff.Compare(deployment, mutatedDeployment)
	if err != nil {
		e := fmt.Sprintf("could not mutate deployment: %v", err)
		return ReviewResponse(request.UID, false, http.StatusBadRequest, e), err
	}

	patchb, err := json.Marshal(patch)
	if err != nil {
		e := fmt.Sprintf("could not mutate deployment: %v", err)
		return ReviewResponse(request.UID, false, http.StatusBadRequest, e), err
	}

	return PatchReviewResponse(request.UID, patchb)
}

// Deployment extracts a deployment of an admission request
func (a DeploymentAdmitter) deployment(request *admissionv1.AdmissionRequest) (*appsv1.Deployment, error) {
	if request.Kind.Kind != "Deployment" {
		return nil, fmt.Errorf("only deployments are supported here")
	}

	d := appsv1.Deployment{}
	if err := json.Unmarshal(request.Object.Raw, &d); err != nil {
		return nil, err
	}

	return &d, nil
}
