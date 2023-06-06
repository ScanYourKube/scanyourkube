package annotator

import (
	"github.com/Masterminds/semver"
)

type KeelPolicyAnnotator struct {
	ImageTag string
}

func (a KeelPolicyAnnotator) Annotate(annotations map[string]string) map[string]string {
	_, err := semver.NewVersion(a.ImageTag)
	if err != nil {
		annotations["keel.sh/policy"] = "force"
		annotations["keel.sh/match-tag"] = "true"
		return annotations
	}
	annotations["keel.sh/policy"] = "minor"
	return annotations
}
