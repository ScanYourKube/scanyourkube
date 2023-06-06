package kubeclarity_api_runtime

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	http_client "github.com/scanyourkube/cronjob/api/http_client"
	dto_application "github.com/scanyourkube/cronjob/dto/api/application"
	dto_application_resource "github.com/scanyourkube/cronjob/dto/api/application_resource"
	dto_kubeclarity_api "github.com/scanyourkube/cronjob/dto/api/kubeclarity"
	dto_scan "github.com/scanyourkube/cronjob/dto/api/scan"
	dto_vulnerability "github.com/scanyourkube/cronjob/dto/api/vulnerability"

	log "github.com/sirupsen/logrus"
)

type IKubeClarityApi interface {
	GetScanProgress() (*dto_scan.ScanProgress, error)
	StartScan(scanConfig dto_scan.RuntimeScanConfig) error
	GetApplications() ([]dto_application.Application, error)
	GetApplicationResourceByApplicationId(applicationId string) ([]dto_application_resource.ApplicationResource, error)
	GetVulnerabilitiesOfApplicationResource(applicationResourceId string) ([]dto_vulnerability.DtoVulnerability, error)
	GetVulnerabilityDetailsWithPackage(vulnerabilityId string, packageId string) (*dto_vulnerability.DtoVulnerabilityDetails, error)
}

type KubeClarityApi struct {
	baseUrl                string
	client                 http_client.IMyHttpClient
	vulnerabilitesToIgnore []string
}

func NewKubeClarityApi(client http_client.IMyHttpClient, baseUrl string, vulnerabilitesToIgnore []string) KubeClarityApi {

	return KubeClarityApi{
		client:                 client,
		baseUrl:                baseUrl,
		vulnerabilitesToIgnore: vulnerabilitesToIgnore,
	}
}

func (a KubeClarityApi) GetScanProgress() (*dto_scan.ScanProgress, error) {
	response, err := a.client.Get(a.baseUrl + "/api/runtime/scan/progress")
	if err != nil {
		logError(err)
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		logError(err)
		return nil, err
	}
	var scanProgress dto_scan.ScanProgress
	err = json.Unmarshal(responseData, &scanProgress)

	if err != nil {
		logError(err)
		return nil, err
	}

	return &scanProgress, nil
}

func (a KubeClarityApi) StartScan(scanConfig dto_scan.RuntimeScanConfig) error {
	json, _ := json.Marshal(scanConfig)

	response, err := a.client.Put(a.baseUrl+"/api/runtime/scan/start", json)
	if err != nil {
		logError(err)
		return err
	}

	if response.StatusCode >= 400 && response.StatusCode <= 599 {
		return fmt.Errorf("returned response code didn't show a successful state status code: %d", response.StatusCode)
	}

	return nil
}

func (a KubeClarityApi) GetVulnerabilitiesOfApplicationResource(applicationResourceId string) ([]dto_vulnerability.DtoVulnerability, error) {
	var page int = 1
	const pageSize int = 50
	var vulnerabilites []dto_vulnerability.DtoVulnerability
	for {
		log.Debugf("Getting vulnerabilites from page %d", page)
		queryParams := map[string]string{
			"page":                           strconv.Itoa(page),
			"pageSize":                       strconv.Itoa(pageSize),
			"sortKey":                        "vulnerabilityName",
			"sortDir":                        "ASC",
			"hasFixVersion%5Bis%5D":          "true",
			"vulnerabilitySource":            "RUNTIME",
			"vulnerabilitySeverity%5Bgte%5D": "HIGH",
			"applicationResourceID":          applicationResourceId,
			"vulnerabilityName%5BisNot%5D":   strings.Join(a.vulnerabilitesToIgnore, ","),
		}
		receivedVulnerabilites, err := a.getListOfVulnerabilitesFromApi(a.baseUrl + "/api/vulnerabilities?" + urlBuilder(queryParams))
		if err != nil {
			logError(err)
			return nil, err
		}

		if len(receivedVulnerabilites) == 0 {
			log.Debugf("Got %d vulnerabilites\n", len(vulnerabilites))
			return vulnerabilites, nil
		}

		log.Debugf("Got vulnerabilites from page %d\n", page)
		vulnerabilites = append(vulnerabilites, receivedVulnerabilites...)
		page++
	}
}

func (a KubeClarityApi) GetVulnerabilityDetailsWithPackage(vulnerabilityId string, packageId string) (*dto_vulnerability.DtoVulnerabilityDetails, error) {
	log.Debugf("Getting vulnerability details for vulnerability: %s and package %s", vulnerabilityId, packageId)

	vulnerability, err := a.getVulnerabilityDetailsFromApi(a.baseUrl + fmt.Sprintf("/api/vulnerabilities/%s/%s", vulnerabilityId, packageId))
	if err != nil {
		logError(err)
		return nil, err
	}

	log.Debugf("Got vulnerability details for vulnerability: %s and package %s", vulnerabilityId, packageId)
	return vulnerability, nil
}

func (a KubeClarityApi) GetApplicationResourceByApplicationId(applicationId string) ([]dto_application_resource.ApplicationResource, error) {
	var page int = 1
	const pageSize int = 50
	var applicationResources []dto_application_resource.ApplicationResource
	for {
		log.Debugf("Getting application resources from page %d", page)
		queryParams := map[string]string{
			"page":          strconv.Itoa(page),
			"pageSize":      strconv.Itoa(pageSize),
			"sortKey":       "resourceName",
			"sortDir":       "ASC",
			"applicationID": applicationId,
		}
		receivedApplicationResources, err := a.getListOfApplicationResourcesFromApi(a.baseUrl + "/api/applicationResources?" + urlBuilder(queryParams))
		if err != nil {
			logError(err)
			return nil, err
		}

		if len(receivedApplicationResources) == 0 {
			log.Debugf("Got %d application resources\n", len(applicationResources))
			return applicationResources, nil
		}

		log.Debugf("Got application resources from page %d\n", page)
		applicationResources = append(applicationResources, receivedApplicationResources...)
		page++
	}
}

func (a KubeClarityApi) GetApplications() ([]dto_application.Application, error) {
	var page int = 1
	const pageSize int = 50
	var applications []dto_application.Application
	for {
		log.Debugf("Getting applications from page %d", page)
		queryParams := map[string]string{
			"page":     strconv.Itoa(page),
			"pageSize": strconv.Itoa(pageSize),
			"sortKey":  "applicationName",
			"sortDir":  "ASC",
		}
		receivedApplications, err := a.getListOfApplicationsFromApi(a.baseUrl + "/api/applications?" + urlBuilder(queryParams))
		if err != nil {
			logError(err)
			return nil, err
		}

		if len(receivedApplications) == 0 {
			log.Debugf("Got %d applications\n", len(applications))
			return applications, nil
		}

		log.Debugf("Got applications from page %d\n", page)
		applications = append(applications, receivedApplications...)
		page++
	}
}

func (a KubeClarityApi) getVulnerabilityDetailsFromApi(url string) (*dto_vulnerability.DtoVulnerabilityDetails, error) {
	log.Debugf("Get vulnerability with url %s", url)
	response, err := a.client.Get(url)
	if err != nil {
		logError(err)
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		logError(err)
		return nil, err
	}

	if response.StatusCode >= 400 && response.StatusCode <= 599 {
		return nil, fmt.Errorf("returned response code didn't show a successful state status code: %d with body: %s", response.StatusCode, responseData)
	}

	var apiResponse dto_vulnerability.DtoVulnerabilityDetails
	json.Unmarshal(responseData, &apiResponse)

	if err != nil {
		logError(err)
		return nil, err
	}

	return &apiResponse, nil
}

func (a KubeClarityApi) getListOfVulnerabilitesFromApi(url string) ([]dto_vulnerability.DtoVulnerability, error) {
	log.Debugf("Get list of vulnerabilities with url %s", url)
	response, err := a.client.Get(url)
	if err != nil {
		logError(err)
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		logError(err)
		return nil, err
	}
	var apiResponse dto_kubeclarity_api.DtoKubeClarityApiResponse[dto_vulnerability.DtoVulnerability]
	json.Unmarshal(responseData, &apiResponse)

	if err != nil {
		logError(err)
		return nil, err
	}

	return apiResponse.Items, nil
}

func (a KubeClarityApi) getListOfApplicationResourcesFromApi(url string) ([]dto_application_resource.ApplicationResource, error) {
	log.Debugf("Get list of application resources with url %s\n", url)
	response, err := a.client.Get(url)
	if err != nil {
		logError(err)
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		logError(err)
		return nil, err
	}
	var apiResponse dto_kubeclarity_api.DtoKubeClarityApiResponse[dto_application_resource.ApplicationResource]
	json.Unmarshal(responseData, &apiResponse)

	if err != nil {
		logError(err)
		return nil, err
	}

	return apiResponse.Items, nil
}

func (a KubeClarityApi) getListOfApplicationsFromApi(url string) ([]dto_application.Application, error) {
	log.Debugf("Get list of applications with url %s\n", url)
	response, err := a.client.Get(url)
	if err != nil {
		logError(err)
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		logError(err)
		return nil, err
	}
	var apiResponse dto_kubeclarity_api.DtoKubeClarityApiResponse[dto_application.Application]
	json.Unmarshal(responseData, &apiResponse)

	if err != nil {
		logError(err)
		return nil, err
	}

	return apiResponse.Items, nil
}

func logError(err error) {
	log.Errorf("Got an error %v\n", err)
}

func urlBuilder(params map[string]string) string {
	var urlBuilder string
	for key, value := range params {
		urlBuilder = urlBuilder + key + "=" + value + "&"
	}
	return strings.TrimSuffix(urlBuilder, "&")
}
