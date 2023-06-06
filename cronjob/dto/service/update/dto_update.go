package dto_service_update

type UpdateServiceImageDto struct {
	Id           string `json:"id,omitempty"`
	ResourceName string `json:"resourceName,omitempty"`
	ResourceHash string `json:"resourceHash,omitempty"`
}
