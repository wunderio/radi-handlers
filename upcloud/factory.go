package upcloud

const (
	CONFIG_KEY_UPCLOUD = "upcloud"
)

// A backend for pulling Config for UpCloud configuration
type UpcloudFactory interface {
	ServiceWrapper() *UpcloudServiceWrapper
}

