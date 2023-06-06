package update_service_rules

import dto_resource "github.com/scanyourkube/cronjob/dto/service/resource"

type IAutomaticUpdateRule interface {
	ShouldUpdateImage(vulnerability dto_resource.Vulnerability) bool
	GetName() string
}
