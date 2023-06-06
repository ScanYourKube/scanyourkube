package update_service_rules

import (
	dto_service_scan "github.com/scanyourkube/cronjob/dto/service/resource"
	"testing"
)

// arg1 means argument 1 and the expected stands for the 'result we expect'
type vulnerablePackageHasAFixVersionUpdateRuleTest struct {
	arg1     dto_service_scan.Vulnerability
	expected bool
}

var vulnerablePackagesHasFixVersionTests = []vulnerablePackageHasAFixVersionUpdateRuleTest{
	{
		dto_service_scan.Vulnerability{
			VulnerabilityName: "DUMMY_VALUE",
			FixVersion:        "",
		},
		false,
	},
	{
		dto_service_scan.Vulnerability{
			VulnerabilityName: "DUMMY_VALUE",
			FixVersion:        "DUMMY_VALUE",
		},
		true,
	},
}

func TestVulnerablePackageHasAFixVersionUpdateRule(t *testing.T) {

	for _, test := range vulnerablePackagesHasFixVersionTests {
		updateRule := NewVulnerablePackageHasAFixVersionUpdateRule()
		if output := updateRule.ShouldUpdateImage(test.arg1); output != test.expected {
			t.Errorf("Output %t not equal to expected %t, input was %v", output, test.expected, test.arg1)
		}
	}
}
