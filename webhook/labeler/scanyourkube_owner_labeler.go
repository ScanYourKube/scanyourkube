package labeler

import (
	"regexp"

	log "github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
)

type ScanYourKubeOwnerLabeler struct {
	Owner     string
	Operation admissionv1.Operation
}

func (labeler ScanYourKubeOwnerLabeler) Label(labels map[string]string) map[string]string {
	isUser, err := regexp.MatchString(`^u-[a-z,0-9]*$`, labeler.Owner)
	log.Debugf("Modifying user is a rancher user: %t", isUser)
	if err != nil || labeler.Operation == admissionv1.Update && !isUser {
		return labels
	}

	log.Infof("Labelling scanyourkube.io/owner with %s", labeler.Owner)

	labels["scanyourkube.io/owner"] = labeler.Owner
	return labels
}
