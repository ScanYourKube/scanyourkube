package update_service_rules

import (
	dto_resource "github.com/scanyourkube/cronjob/dto/service/resource"

	log "github.com/sirupsen/logrus"
)

type CVSSUpdateRule struct {
	name string
}

func (rule CVSSUpdateRule) GetName() string {
	return rule.name
}

func NewCVSSUpdateRule() CVSSUpdateRule {
	return CVSSUpdateRule{
		name: "Check if CVS score is over defined margin update rule",
	}
}

func (CVSSUpdateRule) ShouldUpdateImage(vulnerability dto_resource.Vulnerability) bool {
	if vulnerability.CvssSeverity == nil {
		return false
	}
	log.Debugf("Check if CVS score is over defined margin update rule cvs_score: %f, severity: %v", vulnerability.CvssBaseScore, vulnerability.CvssSeverity)
	return (*vulnerability.AttackVector == dto_resource.NETWORK || *vulnerability.AttackVector == dto_resource.ADJACENT) && (*vulnerability.CvssSeverity == dto_resource.HIGH || *vulnerability.CvssSeverity == dto_resource.CRITICAL)
}
