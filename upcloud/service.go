package upcloud

import (
	upcloud_client "github.com/Jalle19/upcloud-go-sdk/upcloud/client"
	upcloud_service "github.com/Jalle19/upcloud-go-sdk/upcloud/service"
)

/**
 * UpCloud SDK service wrapper
 */

// Constructor for UpcloudServiceSettings
func New_UpcloudServiceSettings(client upcloud_client.Client, hosts []string) *UpcloudServiceSettings {
	return &UpcloudServiceSettings{
		client: client,
		hosts: hosts,
	}
}

// Settings for the 
type UpcloudServiceSettings struct {
	client upcloud_client.Client

	hosts []string
}
// Get an Upcloud service from these settings
func (serviceSettings UpcloudServiceSettings) Service() *upcloud_service.Service {
	return New_UpcloudServiceFromClient(serviceSettings.client)
}

// Get an Upcloud service from these settings
func (serviceSettings UpcloudServiceSettings) ServiceWrapper() *UpcloudServiceWrapper {
	service := serviceSettings.Service()
	return New_UpcloudServiceWrapper(*service)
}


// Constructor for upcloud Service from a client
func New_UpcloudServiceFromClient(client upcloud_client.Client) *upcloud_service.Service {
	service := upcloud_service.New(&client)
	return service
}


// Constructor for UpcloudServiceWrapper
func New_UpcloudServiceWrapper(service upcloud_service.Service) *UpcloudServiceWrapper {
	return &UpcloudServiceWrapper{
		Service: service,
	}
}

// Wrapper for the upcloud service, so that we can limit operations
type UpcloudServiceWrapper struct {
	upcloud_service.Service
}
