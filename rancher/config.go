package rancher

import (
	rancher_client "github.com/rancher/go-rancher/client"
)

const (
	CONFIG_KEY_RANCHER = "rancher"
)

type RancherConfigSource interface {
	Client() *rancher_client.RancherClient
	Settings() RancherSettings
}

// Create a new RancherClient from ClientOpts (small wrapper to simplify imports)
func MakeRancherClientFromSettings(settings *RancherSettings) *rancher_client.RancherClient {
	opts := settings.rancher_client_ClientOpts()
	client, _ := rancher_client.NewRancherClient(&opts)
	return client
}