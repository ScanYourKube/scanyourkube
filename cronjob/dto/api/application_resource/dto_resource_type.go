package dto_application_resource

type ResourceType string

// List of ResourceType
const (
	IMAGE     ResourceType = "IMAGE"
	DIRECTORY ResourceType = "DIRECTORY"
	FILE      ResourceType = "FILE"
	ROOTFS    ResourceType = "ROOTFS"
)
