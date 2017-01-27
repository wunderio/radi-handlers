package rancher

import (
	"time"

	rancher_client "github.com/rancher/go-rancher/client"	
)

// All of the settings needed to get a Rancher client connection : see github.com/rancher/go-rancher/client/ClientOpts
type RancherSettings struct {
	Url       string	`yaml:"Url"`
	AccessKey string	`yaml:"AccessKey"`
	SecretKey string	`yaml:"SecretKey"`
	Timeout   time.Duration  `yaml:"Timeout"`
}

// Convert RancherSettings to the rancher library client settings object natively.
func (settings *RancherSettings) rancher_client_ClientOpts() rancher_client.ClientOpts {
	return rancher_client.ClientOpts{
		Url: settings.Url,
		AccessKey: settings.AccessKey,
		SecretKey: settings.SecretKey,
		Timeout: settings.Timeout,
	}
}
// Is this struct empty
func (settings *RancherSettings) Empty() bool {
	return settings.Url == ""
}
