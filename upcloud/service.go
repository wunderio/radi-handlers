package upcloud

import (
	log "github.com/Sirupsen/logrus"

	upcloud_client "github.com/Jalle19/upcloud-go-sdk/upcloud/client"
	upcloud_service "github.com/Jalle19/upcloud-go-sdk/upcloud/service"
)

/**
 * UpCloud SDK service wrapper
 */

// Constructor for UpcloudServiceSettings
func New_UpcloudServiceSettings(client upcloud_client.Client) *UpcloudServiceSettings {
	return &UpcloudServiceSettings{
		client: client,
	}
}

// Settings for the 
type UpcloudServiceSettings struct {
	client upcloud_client.Client
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


// Define some values that can be used by the ServiceWrapper to limit and configure it
type UpcloudBuilderSettings struct {
	Hosts []string  `yml:"Hosts"`
}
// Merge settings
func (settings *UpcloudBuilderSettings) Merge(merge UpcloudBuilderSettings) {
	// merge hosts
	for _, host := range merge.Hosts {
		exists := false
		for _, existing := range settings.Hosts {
			if existing == host {
				exists = true
				break
			}
		}
		if !exists {
			settings.Hosts = append(settings.Hosts, host)
		}
	}

	log.WithFields(log.Fields{"settings": settings}).Debug("Merged UpCloud settings")
}
// It doesn't want to automatically marshal, so do it manually @TODO why isn't it unmarshalling automatically?
func (settings *UpcloudBuilderSettings) UnmarshalYAML(unmarshal func(interface{}) error) error {
	placeholder := map[string][]string{}
	if err := unmarshal(&placeholder); err != nil {
		return err
	}

	if hosts, defined := placeholder["Hosts"]; defined {
		for _, host := range hosts {
			exists := false
			for _, existing := range settings.Hosts {
				if existing == host {
					exists = true
					break
				}
			}
			if !exists {
				settings.Hosts = append(settings.Hosts, host)
			}
		}
	}
	return nil
}

// Wrapper for the upcloud service, so that we can limit operations
type UpcloudServiceWrapper struct {
	upcloud_service.Service

}
