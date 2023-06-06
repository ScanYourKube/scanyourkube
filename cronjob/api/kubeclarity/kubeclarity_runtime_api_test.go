package kubeclarity_api_runtime_test

import (
	"encoding/json"
	"fmt"
	"github.com/scanyourkube/cronjob/api/http_client"
	kubeclarity_api_runtime "github.com/scanyourkube/cronjob/api/kubeclarity"
	"github.com/scanyourkube/cronjob/dto/api/scan"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
)

func TestGetScanProgress(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.Rng.Seed(1024)
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("returns a valid ScanProgress", prop.ForAll(
		func(input scan.ScanProgress) bool {
			responseBody, _ := json.Marshal(input)
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.Write([]byte(responseBody))
			}))
			defer server.Close()

			myHttpClient := http_client.NewMyHttpClient(server.Client())

			api := kubeclarity_api_runtime.NewKubeClarityApi(myHttpClient, server.URL, make([]string, 0))
			scanProgress, err := api.GetScanProgress()

			return assert.NoError(t, err) && assert.NotNil(t, scanProgress)
		},
		genScanProgress(),
	))

	properties.Property("returns an error when the HTTP GET request fails", prop.ForAll(
		func(statusCode int) bool {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(statusCode)
			}))
			defer server.Close()

			myHttpClient := http_client.NewMyHttpClient(server.Client())

			api := kubeclarity_api_runtime.NewKubeClarityApi(myHttpClient, server.URL, make([]string, 0))
			_, err := api.GetScanProgress()
			return assert.Error(t, err)
		},
		gen.IntRange(400, 599),
	))

	properties.Property("returns an error when the response body cannot be read", prop.ForAll(
		func(responseBody string) bool {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.Header().Set("Content-Length", "1")
				rw.Write([]byte(responseBody))
			}))
			defer server.Close()

			myHttpClient := http_client.NewMyHttpClient(server.Client())

			api := kubeclarity_api_runtime.NewKubeClarityApi(myHttpClient, server.URL, make([]string, 0))
			_, err := api.GetScanProgress()
			return assert.Error(t, err)
		},
		gen.AlphaString(),
	))

	properties.Property("returns an error when the response body is not valid JSON", prop.ForAll(
		func(responseBody string) bool {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.Write([]byte(responseBody))
			}))
			defer server.Close()

			myHttpClient := http_client.NewMyHttpClient(server.Client())

			api := kubeclarity_api_runtime.NewKubeClarityApi(myHttpClient, server.URL, make([]string, 0))
			_, err := api.GetScanProgress()
			return assert.Error(t, err)
		},
		gen.AlphaString().SuchThat(func(s string) bool { return !json.Valid([]byte(s)) }),
	))

	properties.Property("returns an error when the client has not a valid URL specified", prop.ForAll(
		func() bool {
			myHttpClient := http_client.NewMyHttpClient(&http.Client{})
			api := kubeclarity_api_runtime.NewKubeClarityApi(myHttpClient, "", make([]string, 0))
			_, err := api.GetScanProgress()
			return assert.Error(t, err)
		},
	))
	properties.TestingRun(t)
}

func TestStartScan(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("returns no error when the HTTP PUT request succeeds", prop.ForAll(
		func(scanConfig scan.RuntimeScanConfig) bool {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, http.MethodPut, req.Method)
				assert.Equal(t, "/api/runtime/scan/start", req.URL.Path)
			}))
			defer server.Close()

			myHttpClient := http_client.NewMyHttpClient(server.Client())

			api := kubeclarity_api_runtime.NewKubeClarityApi(myHttpClient, server.URL, make([]string, 0))
			err := api.StartScan(scanConfig)
			return assert.NoError(t, err)
		},
		genRuntimeScanConfig(),
	))

	properties.Property("returns an error when the HTTP PUT request fails", prop.ForAll(
		func(statusCode int) bool {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(statusCode)
			}))
			defer server.Close()
			myHttpClient := http_client.NewMyHttpClient(server.Client())

			api := kubeclarity_api_runtime.NewKubeClarityApi(myHttpClient, server.URL, make([]string, 0))
			err := api.StartScan(scan.RuntimeScanConfig{})
			return assert.Error(t, err)
		},
		gen.IntRange(400, 599),
	))

	properties.TestingRun(t)
}

func TestGetVulnerabilitesOfApplicationResource(t *testing.T) {
	cases := []struct {
		input struct {
			ignoredCVEs []string
		}
		expected int
	}{
		{
			input: struct {
				ignoredCVEs []string
			}{
				ignoredCVEs: []string{"CVE-0000-00000"},
			},
			expected: 0,
		},
		{
			input: struct {
				ignoredCVEs []string
			}{
				ignoredCVEs: []string{"CVE-0000-00001"},
			},
			expected: 1,
		},
		{
			input: struct {
				ignoredCVEs []string
			}{
				ignoredCVEs: []string{},
			},
			expected: 1,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("vulnerabilityName[isNot]") == "CVE-0000-00000" ||
			req.URL.Query().Get("page") == "2" {
			rw.Write([]byte(`{
				"total": 0,
				"items": []
			  }`))
		} else {

			rw.Write([]byte(`{
			"total": 1,
			"items": [
			  {
				"id": "string",
				"vulnerabilityID": "string",
				"vulnerabilityName": "CVE-0000-00000",
				"packageName": "string",
				"packageVersion": "string",
				"packageID": "string",
				"severity": "CRITICAL",
				"cvssSeverity": "CRITICAL",
				"cvssBaseScore": 0,
				"applicationResources": 0,
				"applications": 0,
				"fixVersion": "string",
				"reportingScanners": [
				  "string"
				],
				"source": "CICD"
			  }
			]
		  }`))
		}
	}))
	defer server.Close()
	myHttpClient := http_client.NewMyHttpClient(server.Client())

	for _, c := range cases {

		api := kubeclarity_api_runtime.NewKubeClarityApi(myHttpClient, server.URL, c.input.ignoredCVEs)
		vulnerabilites, err := api.GetVulnerabilitiesOfApplicationResource("DUMMY")
		assert.NoError(t, err)
		assert.Equal(t, c.expected, len(vulnerabilites))
	}

}

func TestGetVulnerabilityDetailsWithPackage(t *testing.T) {
	cases := []struct {
		input struct {
			vulnerabilityId string
			packageId       string
		}
		expected bool
	}{
		{
			input: struct {
				vulnerabilityId string
				packageId       string
			}{
				vulnerabilityId: "1",
				packageId:       "1",
			},
			expected: true,
		},
		{
			input: struct {
				vulnerabilityId string
				packageId       string
			}{
				vulnerabilityId: "1",
				packageId:       "2",
			},
			expected: false,
		},
		{
			input: struct {
				vulnerabilityId string
				packageId       string
			}{
				vulnerabilityId: "2",
				packageId:       "1",
			},
			expected: false,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/vulnerabilities/1/1" {
			rw.WriteHeader(http.StatusNotFound)
		} else {

			rw.Write([]byte(`{
				"vulnerability": {
				  "id": "1",
				  "vulnerabilityID": "string",
				  "vulnerabilityName": "string",
				  "packageName": "string",
				  "packageVersion": "string",
				  "packageID": "1",
				  "severity": "CRITICAL",
				  "cvssSeverity": "CRITICAL",
				  "cvssBaseScore": 0,
				  "applicationResources": 0,
				  "applications": 0,
				  "fixVersion": "string",
				  "reportingScanners": [
					"string"
				  ],
				  "source": "CICD"
				},
				"scanDate": "2023-06-03T12:00:16.188Z",
				"description": "string",
				"links": [
				  "string"
				],
				"cvss": {
				  "cvssV3Vector": {
					"vector": "string",
					"attackVector": "NETWORK",
					"attackComplexity": "LOW",
					"privilegesRequired": "NONE",
					"userInteraction": "NONE",
					"scope": "UNCHANGED",
					"confidentiality": "NONE",
					"integrity": "NONE",
					"availability": "NONE"
				  },
				  "cvssV3Metrics": {
					"severity": "CRITICAL",
					"baseScore": 0,
					"exploitabilityScore": 0,
					"impactScore": 0
				  }
				}
			  }`))
		}
	}))
	defer server.Close()
	myHttpClient := http_client.NewMyHttpClient(server.Client())

	for _, c := range cases {

		api := kubeclarity_api_runtime.NewKubeClarityApi(myHttpClient, server.URL, make([]string, 0))
		vulnerability, err := api.GetVulnerabilityDetailsWithPackage(c.input.vulnerabilityId, c.input.packageId)

		if c.expected {
			assert.NoError(t, err)
			assert.NotNil(t, vulnerability)
		} else {
			assert.Error(t, err)
			assert.Nil(t, vulnerability)
		}
	}

}

func TestGetApplicationResourceByApplicationId(t *testing.T) {
	cases := []struct {
		input struct {
			applicationId string
		}
		expected bool
	}{
		{
			input: struct {
				applicationId string
			}{
				applicationId: "1",
			},
			expected: true,
		},
		{
			input: struct {
				applicationId string
			}{
				applicationId: "2",
			},
			expected: false,
		},
		{
			input: struct {
				applicationId string
			}{
				applicationId: "",
			},
			expected: false,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("applicationID") != "1" ||
			req.URL.Query().Get("page") == "2" {
			rw.Write([]byte(`{
				"total": 0,
				"items": [
				  
				]
			  }`))
		} else {

			rw.Write([]byte(`{
				"total": 1,
				"items": [
				  {
					"id": "string",
					"resourceName": "string",
					"resourceHash": "string",
					"resourceType": "IMAGE",
					"vulnerabilities": [
					  {
						"severity": "CRITICAL",
						"count": 0
					  }
					],
					"cisDockerBenchmarkResults": [
					  {
						"level": "INFO",
						"count": 0
					  }
					],
					"applications": 0,
					"packages": 0,
					"reportingSBOMAnalyzers": [
					  "string"
					]
				  }
				]
			  }`))
		}
	}))
	defer server.Close()
	myHttpClient := http_client.NewMyHttpClient(server.Client())

	for _, c := range cases {

		api := kubeclarity_api_runtime.NewKubeClarityApi(myHttpClient, server.URL, make([]string, 0))
		applicationResources, err := api.GetApplicationResourceByApplicationId(c.input.applicationId)

		if c.expected {
			assert.NoError(t, err)
			assert.NotNil(t, applicationResources)
			assert.Equal(t, 1, len(applicationResources))
		} else {
			assert.Nil(t, applicationResources)
		}
	}

}

func TestGetApplications(t *testing.T) {
	cases := []struct {
		input struct {
			pages int
		}
		expected int
	}{
		{
			input: struct {
				pages int
			}{
				pages: 2,
			},
			expected: 1,
		},
		{
			input: struct {
				pages int
			}{
				pages: 3,
			},
			expected: 2,
		},
		{
			input: struct {
				pages int
			}{
				pages: 5,
			},
			expected: 4,
		},
	}

	for _, c := range cases {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.URL.Query().Get("page") == fmt.Sprintf("%d", c.input.pages) {
				rw.Write([]byte(`{
				"total": 0,
				"items": [
				  
				]
			  }`))
			} else {

				rw.Write([]byte(`{
					"total": 1,
					"items": [
					  {
						"id": "string",
						"applicationName": "string",
						"applicationType": "POD",
						"labels": [
						  "string"
						],
						"vulnerabilities": [
						  {
							"severity": "CRITICAL",
							"count": 0
						  }
						],
						"cisDockerBenchmarkResults": [
						  {
							"level": "INFO",
							"count": 0
						  }
						],
						"environments": [
						  "string"
						],
						"applicationResources": 0,
						"packages": 0
					  }
					]
				  }`))
			}
		}))
		defer server.Close()
		myHttpClient := http_client.NewMyHttpClient(server.Client())

		api := kubeclarity_api_runtime.NewKubeClarityApi(myHttpClient, server.URL, make([]string, 0))
		applicationResources, err := api.GetApplications()

		assert.Nil(t, err)
		assert.Equal(t, c.expected, len(applicationResources))
	}

}

func genScanProgress() gopter.Gen {
	return gopter.CombineGens(
		gen.Int32Range(0, 100),
		gen.SliceOf(gen.AlphaString()),
		gen.Time(),
	).Map(func(values []interface{}) scan.ScanProgress {
		return scan.ScanProgress{
			Scanned:           values[0].(int32),
			ScannedNamespaces: values[1].([]string),
			StartTime:         values[2].(time.Time),
		}
	})
}

func genRuntimeScanConfig() gopter.Gen {
	return gen.SliceOf(gen.AlphaString()).Map(
		func(namespaces []string) scan.RuntimeScanConfig {
			return scan.RuntimeScanConfig{
				Namespaces: namespaces,
			}
		})
}
