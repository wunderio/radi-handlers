package rancher

import (
	rancher_client "github.com/rancher/go-rancher/client"
)

const (
	CONFIG_KEY_RANCHER = "rancher"
)

// Define a source of config for the rancher handler
type RancherConfigSource interface {
	RancherClient() *rancher_client.RancherClient

	RancherClientSettings() RancherClientSettings
	RancherEnvironmentSettings() RancherEnvironmentSettings
}

// Create a new RancherClient from ClientOpts (small wrapper to simplify imports)
func MakeRancherClientFromSettings(settings RancherClientSettings) (*rancher_client.RancherClient, error) {
	opts := settings.rancher_client_ClientOpts()
	client, err := rancher_client.NewRancherClient(&opts)
	return client, err
}
