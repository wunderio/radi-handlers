package rancher

import (
	"time"
)

/**
 * Handler builder and settings to create rancher handlers
 */

type RancherBuilder struct {
	settings RancherSettings

	parent   api_api.API
	handlers api_handler.Handlers
}

// All of the settings needed to get a Rancher client connection : see github.com/rancher/go-rancher/client/ClientOpts
type RancherSettings struct {
	Url       string        `yaml:"Url"`
	AccessKey string        `yaml:"AccessKey"`
	SecretKey string        `yaml:"SecretKey"`
	Timeout   time.Duration `yaml:"Timeout"`
}
