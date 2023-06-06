package update_service_rules

import (
	dto_resource "github.com/scanyourkube/cronjob/dto/service/resource"

	log "github.com/sirupsen/logrus"
)

type VulnerablePackageHasAFixVersionUpdateRule struct {
	name string
}

func (rule VulnerablePackageHasAFixVersionUpdateRule) GetName() string {
	return rule.name
}

func NewVulnerablePackageHasAFixVersionUpdateRule() VulnerablePackageHasAFixVersionUpdateRule {
	return VulnerablePackageHasAFixVersionUpdateRule{
		name: "Check if vulnerable package has a fixed version",
	}
}

func (VulnerablePackageHasAFixVersionUpdateRule) ShouldUpdateImage(vulnerability dto_resource.Vulnerability) bool {
	log.Debugf("Check if vulnerable package %s has a fixed version: %s", vulnerability.VulnerabilityName, vulnerability.FixVersion)
	return vulnerability.FixVersion != ""
}
