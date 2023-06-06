package service_scan

import (
	"time"

	api "github.com/scanyourkube/cronjob/api/kubeclarity"
	dto_scan "github.com/scanyourkube/cronjob/dto/api/scan"
	resource_service "github.com/scanyourkube/cronjob/service/kubernetes"

	log "github.com/sirupsen/logrus"
)

type ScanService struct {
	ScanApi           api.IKubeClarityApi
	KubernetesService resource_service.IKubernetesResourceService
}

type IScanService interface {
	ScanNamespaces() error
	GetScanProgressUntilFinished() (int32, error)
}

func NewScanService(ScanApi api.IKubeClarityApi, KubernetesService resource_service.IKubernetesResourceService) ScanService {
	return ScanService{
		ScanApi:           ScanApi,
		KubernetesService: KubernetesService,
	}
}

func (s ScanService) ScanNamespaces() error {
	scanconfig := dto_scan.RuntimeScanConfig{
		Namespaces: s.KubernetesService.GetNamespaceNames(),
	}

	s.ScanApi.StartScan(scanconfig)

	for {
		currentProgress, err := s.GetScanProgressUntilFinished()
		if err != nil {
			return err
		}

		log.Infof("Got new progress: %d", currentProgress)
		if currentProgress == 100 {
			log.Info("Scan finished")
			break
		}
		log.Infof("Scanned %d \n", currentProgress)
		time.Sleep(1 * time.Second)
	}

	log.Debug("Scanned all namespaces")

	log.Debug("Returning nil as there is no error")
	return nil
}

func (s ScanService) GetScanProgressUntilFinished() (int32, error) {
	scanProgress, err := s.ScanApi.GetScanProgress()
	log.Debug("Scanning")
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}

	percentage := scanProgress.Scanned
	return percentage, nil
}
