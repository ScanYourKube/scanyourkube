package resource_service

import (
	"reflect"
	"testing"

	dto_application_api "github.com/scanyourkube/cronjob/dto/api/application"
	dto_application_resource_api "github.com/scanyourkube/cronjob/dto/api/application_resource"
	dto_vulnerability_api "github.com/scanyourkube/cronjob/dto/api/vulnerability"
	"github.com/scanyourkube/cronjob/dto/api/vulnerability/cvssv3vector"
	api "github.com/scanyourkube/cronjob/testing/mocks/api/kubeclarity"

	"github.com/golang/mock/gomock"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
)

// Generate arbitrary values for VulnerabilitySeverity
func genVulnerabilitySeverity() gopter.Gen {
	return gen.OneConstOf(dto_vulnerability_api.CRITICAL, dto_vulnerability_api.HIGH, dto_vulnerability_api.MEDIUM, dto_vulnerability_api.LOW, dto_vulnerability_api.NEGLIGIBLE)
}

// Generate arbitrary values for AttackVector
func genAttackVector() gopter.Gen {
	return gen.OneConstOf(cvssv3vector.NETWORK, cvssv3vector.ADJACENT, cvssv3vector.LOCAL, cvssv3vector.PHYSICAL)
}

func genPrivilegesRequired() gopter.Gen {
	return gen.OneConstOf(cvssv3vector.NONE, cvssv3vector.LOW, cvssv3vector.HIGH)
}

func genApplicationType() gopter.Gen {
	return gen.OneConstOf(dto_application_api.POD, dto_application_api.DIRECTORY, dto_application_api.LAMBDA)
}

func genResourceType() gopter.Gen {
	return gen.OneConstOf(dto_application_resource_api.IMAGE, dto_application_resource_api.DIRECTORY, dto_application_resource_api.FILE, dto_application_resource_api.ROOTFS)
}

// Generate arbitrary values for DtoVulnerability
func genDtoVulnerability() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&dto_vulnerability_api.DtoVulnerability{}), map[string]gopter.Gen{
		"Id":                gen.Identifier(),
		"VulnerabilityID":   gen.Identifier(),
		"VulnerabilityName": gen.AnyString(),
		"PackageName":       gen.AnyString(),
		"PackageVersion":    gen.AnyString(),
		"FixVersion":        gen.AnyString(),
		"Severity":          gen.PtrOf(genVulnerabilitySeverity()),
		"CvssSeverity":      gen.PtrOf(genVulnerabilitySeverity()),
		"CvssBaseScore":     gen.Float32Range(0.0, 10.0),
	})
}

// Generate arbitrary values for DtoCvss
func genDtoCvss() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&dto_vulnerability_api.DtoCvss{}), map[string]gopter.Gen{
		"CvssV3Vector": gen.Struct(reflect.TypeOf(&cvssv3vector.Cvssv3Vector{}), map[string]gopter.Gen{
			"AttackVector":       gen.PtrOf(genAttackVector()),
			"PrivilegesRequired": gen.PtrOf(genPrivilegesRequired()),
			"Vector":             gen.AnyString(),
		}),
	})
}

// Generate arbitrary values for DtoVulnerabilityDetails
func genDtoVulnerabilityDetails() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&dto_vulnerability_api.DtoVulnerabilityDetails{}), map[string]gopter.Gen{
		"Description":   gen.AnyString(),
		"Links":         gen.AnyString(),
		"Vulnerability": genDtoVulnerability(),
		"Cvss":          genDtoCvss(),
	})
}

func genDtoApplication() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&dto_application_api.Application{}), map[string]gopter.Gen{
		"Id":                   gen.Identifier(),
		"ApplicationName":      gen.Identifier(),
		"ApplicationType":      gen.PtrOf(genApplicationType()),
		"Labels":               gen.SliceOf(gen.AnyString()),
		"Vulnerabilities":      gen.SliceOf(genVulnerabilityCount()),
		"Environments":         gen.SliceOf(gen.AnyString()),
		"ApplicationResources": gen.Int32(),
		"Packages":             gen.Int32(),
	})
}

func genDtoApplicationResource() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&dto_application_resource_api.ApplicationResource{}), map[string]gopter.Gen{
		"Id":                     gen.Identifier(),
		"ResourceName":           gen.Identifier(),
		"ResourceHash":           gen.Identifier(),
		"ResourceType":           gen.PtrOf(genResourceType()),
		"ReportingSBOMAnalyzers": gen.SliceOf(gen.AnyString()),
		"Vulnerabilities":        gen.Int32(),
		"Packages":               gen.Int32(),
		"Applications":           gen.Int32(),
	})
}

func genVulnerabilityCount() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&dto_vulnerability_api.VulnerabilityCount{}), map[string]gopter.Gen{
		"Severity": gen.PtrOf(genVulnerabilitySeverity()),
		"Count":    gen.Int32(),
	})
}

func TestGetApplicationsFromLastScan(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.Rng.Seed(1024)
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock implementation of api.IKubeClarityApi
	mockScanApi := api.NewMockIKubeClarityApi(ctrl)

	properties.Property("API returns all applications", prop.ForAll(
		func(applications []dto_application_api.Application,
			applicationResources []dto_application_resource_api.ApplicationResource,
			vulnerabilities []dto_vulnerability_api.DtoVulnerability,
			vulnerability dto_vulnerability_api.DtoVulnerabilityDetails) bool {
			mockScanApi.EXPECT().GetApplications().Return(applications, nil).MaxTimes(1)
			mockScanApi.EXPECT().GetApplicationResourceByApplicationId(gomock.Any()).Return(applicationResources, nil).Times(len(applications))
			mockScanApi.EXPECT().GetVulnerabilitiesOfApplicationResource(gomock.Any()).Return(vulnerabilities, nil).Times(len(applicationResources) * len(applications))
			mockScanApi.EXPECT().GetVulnerabilityDetailsWithPackage(gomock.Any(), gomock.Any()).Return(&dto_vulnerability_api.DtoVulnerabilityDetails{
				Description:   vulnerability.Description,
				Links:         vulnerability.Links,
				Vulnerability: vulnerability.Vulnerability,
			}, nil).Times(len(applicationResources) * len(applications) * len(vulnerabilities))
			// Create the ResourceService with the mockScanApi
			resourceService := NewResourceService(mockScanApi)
			result, err := resourceService.GetApplicationsFromLastScan()
			return assert.NoError(t, err) &&
				(len(applications) == 0 || assert.NotNil(t, result)) &&
				assert.Equal(t, len(applications), len(result))
		},
		gen.SliceOf(genDtoApplication()), gen.SliceOf(genDtoApplicationResource()), gen.SliceOf(genDtoVulnerability()), genDtoVulnerabilityDetails(),
	))

	properties.TestingRun(t)
}
