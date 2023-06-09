/*
 * KubeClarity APIs
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package scan

import (
	"time"
)

type ScanProgress struct {
	Status *RuntimeScanStatus `json:"status,omitempty"`
	// Percentage of scanned images from total images (0-100)
	Scanned           int32     `json:"scanned,omitempty"`
	ScannedNamespaces []string  `json:"scannedNamespaces,omitempty"`
	ScanType          *ScanType `json:"scanType,omitempty"`
	StartTime         time.Time `json:"startTime,omitempty"`
}
