package update_service_rules

import (
	dto_service_scan "github.com/scanyourkube/cronjob/dto/service/resource"
	"testing"
)

// arg1 means argument 1 and the expected stands for the 'result we expect'
type updateRuleTest struct {
	arg1     dto_service_scan.Vulnerability
	expected bool
}

func vul(vul dto_service_scan.VulnerabilitySeverity) *dto_service_scan.VulnerabilitySeverity {
	return &vul
}

func attackVector(attackVector dto_service_scan.AttackVector) *dto_service_scan.AttackVector {
	return &attackVector
}

var updateRuleTests = []updateRuleTest{
	{
		dto_service_scan.Vulnerability{
			CvssSeverity: vul(dto_service_scan.CRITICAL),
			AttackVector: attackVector(dto_service_scan.NETWORK),
		},
		true,
	},
	{
		dto_service_scan.Vulnerability{
			CvssSeverity: vul(dto_service_scan.CRITICAL),
			AttackVector: attackVector(dto_service_scan.ADJACENT),
		},
		true,
	},
	{
		dto_service_scan.Vulnerability{
			CvssSeverity: vul(dto_service_scan.HIGH),
			AttackVector: attackVector(dto_service_scan.NETWORK),
		},
		true,
	},
	{
		dto_service_scan.Vulnerability{
			CvssSeverity: vul(dto_service_scan.CRITICAL),
			AttackVector: attackVector(dto_service_scan.ADJACENT),
		},
		true,
	},
	{
		dto_service_scan.Vulnerability{
			CvssSeverity: vul(dto_service_scan.HIGH),
			AttackVector: attackVector(dto_service_scan.NETWORK),
		},
		true,
	},
	{
		dto_service_scan.Vulnerability{
			CvssSeverity: vul(dto_service_scan.LOW),
			AttackVector: attackVector(dto_service_scan.NETWORK),
		},
		false,
	},
	{
		dto_service_scan.Vulnerability{
			CvssSeverity: vul(dto_service_scan.HIGH),
			AttackVector: attackVector(dto_service_scan.PHYSICAL),
		},
		false,
	},
	{
		dto_service_scan.Vulnerability{
			CvssSeverity: vul(dto_service_scan.CRITICAL),
			AttackVector: attackVector(dto_service_scan.LOCAL),
		},
		false,
	},
	{
		dto_service_scan.Vulnerability{
			CvssSeverity: nil,
			AttackVector: attackVector(dto_service_scan.NETWORK),
		},
		false,
	},
}

func TestUpdateRule(t *testing.T) {

	for _, test := range updateRuleTests {
		updateRule := NewCVSSUpdateRule()
		if output := updateRule.ShouldUpdateImage(test.arg1); output != test.expected {
			t.Errorf("Output %t not equal to expected %t, input was %v", output, test.expected, test.arg1)
		}
	}
}
