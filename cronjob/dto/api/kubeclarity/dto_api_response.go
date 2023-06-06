package dto_kubeclarity_api

type DtoKubeClarityApiResponse[T any] struct {
	// Unique ID for the vulnerabilityID and packageID combination.
	Total int `json:"total,omitempty"`
	Items []T `json:"items,omitempty"`
}
