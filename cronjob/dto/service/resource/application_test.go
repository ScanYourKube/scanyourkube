package resource_service

import (
	"fmt"
	dto_application "github.com/scanyourkube/cronjob/dto/api/application"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Generate arbitrary values for DtoVulnerability
func genDtoApplication() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&dto_application.Application{}), map[string]gopter.Gen{
		"Id":              gen.Identifier(),
		"ApplicationName": gen.AnyString(),
		"Environments":    gen.SliceOf(gen.AnyString()),
		"Labels":          gen.SliceOf(gen.AnyString()),
	})
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Test that the Map method preserves the fields of the input DtoVulnerabilityDetails
func TestApplicationMapping(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.Rng.Seed(1024)
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Map preserves fields", prop.ForAll(
		func(dtoApiApplication dto_application.Application) bool {
			dto := Application{}
			application := dto.Map(dtoApiApplication)

			fmt.Printf("dtoApiApplication: %+v\n", dtoApiApplication)
			fmt.Printf("application: %+v\n", application)
			return application.Id == dtoApiApplication.Id &&
				application.ApplicationName == dtoApiApplication.ApplicationName &&
				stringSlicesEqual(application.Environments, dtoApiApplication.Environments) &&
				stringSlicesEqual(application.Labels, dtoApiApplication.Labels)
		},
		genDtoApplication(),
	))

	properties.TestingRun(t)
}
