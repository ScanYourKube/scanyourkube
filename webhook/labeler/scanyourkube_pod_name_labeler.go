package labeler

import (
	log "github.com/sirupsen/logrus"
)

type ScanYourKubePodNameLabeler struct {
	PodName string
}

func (labeler ScanYourKubePodNameLabeler) Label(labels map[string]string) map[string]string {
	log.Infof("Labelling scanyourkube.io/podName with %s", labeler.PodName)

	labels["scanyourkube.io/podName"] = labeler.PodName
	return labels
}
