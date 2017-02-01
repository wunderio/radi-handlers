package rancher

import (
	rancher_client "github.com/rancher/go-rancher/client"
)

// Settings that can come fomr the project configuration.
// @NOTE Not sure what we want to put in here
type RancherSettings struct {
}

// All of the settings needed to get a Rancher client connection : see github.com/rancher/go-rancher/client/ClientOpts
type RancherClientSettings struct {
	Url       string
	AccessKey string
	SecretKey string
}

// Convert RancherSettings to the rancher library client settings object natively.
func (settings *RancherClientSettings) rancher_client_ClientOpts() rancher_client.ClientOpts {
	return rancher_client.ClientOpts{
		Url:       settings.Url,
		AccessKey: settings.AccessKey,
		SecretKey: settings.SecretKey,
	}
}

// All settings related to the environment used for the project
type RancherEnvironmentSettings struct {
	Id string
}
