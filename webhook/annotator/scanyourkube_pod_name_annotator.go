package annotator

import (
	log "github.com/sirupsen/logrus"
)

type ScanYourKubePodNameAnnotator struct {
	PodName string
}

func (annotator ScanYourKubePodNameAnnotator) Annotate(annotations map[string]string) map[string]string {
	log.Infof("Annotating scanyourkube.io/podName with %s", annotator.PodName)

	annotations["scanyourkube.io/podName"] = annotator.PodName
	return annotations
}
