package rancher

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"

	rancher_client "github.com/rancher/go-rancher/client"

	api_config "github.com/wunderkraut/radi-api/operation/config"
)

// A BuilderSettingsConfigWrapper, that interprets build config as yml
type RancherConfigSourceYaml struct {
	configWrapper api_config.ConfigWrapper
	scope         string
	settings      *RancherSettings_yml
	client        *rancher_client.RancherClient
}

// Constructor for RancherConfigSourceYaml
func New_RancherConfigSourceYaml(configWrapper api_config.ConfigWrapper) *RancherConfigSourceYaml {
	return &RancherConfigSourceYaml{
		configWrapper: configWrapper,
	}
}

func (source *RancherConfigSourceYaml) RancherClient() *rancher_client.RancherClient {
	if source.client == nil {
		source.safe()

		if client, err := MakeRancherClientFromSettings(source.RancherClientSettings()); err != nil {
			log.WithError(err).WithFields(log.Fields{"settings": source.RancherClientSettings(), "client": source.client}).Error("rancher-configsourceyml: failed to create new client")
		} else {
			source.client = client
		}
	}
	return source.client
}

func (source *RancherConfigSourceYaml) RancherClientSettings() RancherClientSettings {
	source.safe()
	return source.settings.RancherClientSettings()
}

func (source *RancherConfigSourceYaml) RancherEnvironmentSettings() RancherEnvironmentSettings {
	source.safe()
	return source.settings.RancherEnvironmentSettings()
}

func (source *RancherConfigSourceYaml) DefaultScope() string {
	/**
	 * @TODO come up with better scopes, but it has to match local conf path keys
	 * @SEE configconnect/settings_yml.go which has the same issue
	 */
	return "project"
}

func (source *RancherConfigSourceYaml) safe() {
	if source.settings == nil {
		if err := source.Load(); err != nil {
			log.WithError(err).Error("Could not load rancher configuration")
		}
	}
}

// Retrieve values by parsing bytes from the ConfigWrapper
func (source *RancherConfigSourceYaml) Load() error {
	log.Debug("Loading Rancher config")

	source.settings = &RancherSettings_yml{}

	if sources, err := source.configWrapper.Get(CONFIG_KEY_RANCHER); err == nil {
		for _, scope := range sources.Order() {
			scopedSource, _ := sources.Get(scope)

			// empty out this oobject
			source.scope = scope
			source.settings = &RancherSettings_yml{}

			if err := yaml.Unmarshal(scopedSource, source.settings); err == nil {
				log.WithFields(log.Fields{"settings": source.settings, "scope": source.scope}).Debug("Rancher settings parsed from config yml")
				break
			} else {
				log.WithError(err).WithFields(log.Fields{"scope": scope}).Error("Couldn't marshall yml scope for rancher settings for scope.")
			}
		}
		return nil
	} else {
		log.WithError(err).Error("Error loading Rancher config.")
		return err
	}
}

// Save the current values to the wrapper
func (source *RancherConfigSourceYaml) Save() error {
	/**
	 * @TODO THIS
	 */
	return errors.New("RancherConfigSourceYaml Set operation not yet written.")
}

/**
 * YML Settings holder
 */
type RancherSettings_yml struct {
	Access RancherSettings_Access_yml `yaml:"Access"`
}

type RancherSettings_Access_yml struct {
	Url       string        `yaml:"Url"`
	AccessKey string        `yaml:"AccessKey"`
	SecretKey string        `yaml:"SecretKey"`
	Timeout   time.Duration `yaml:"Timeout"`
}

// convert this to a proper rancher settings struct
func (settings *RancherSettings_yml) RancherClientSettings() RancherClientSettings {
	clientSettings := RancherClientSettings{
		Url:       "http://127.0.0.1:8080/v2-beta",
		AccessKey: "75AAAFB2BCFA7DD83BE2",
		SecretKey: "XMJM2EgZia7ohYKpyfXPZxdq63C3UJEGe8qBYnt4",
		Timeout:   0,
	}
	log.WithFields(log.Fields{"settings": clientSettings}).Debug("OVERRIDE RANCHER CLIENT SETTINGS WITH SOME MANUAL VALUE")
	return clientSettings
}

// convert this to a proper rancher settings struct
func (settings *RancherSettings_yml) RancherEnvironmentSettings() RancherEnvironmentSettings {
	environmentSettings := RancherEnvironmentSettings{
		Id: "default",
	}
	log.WithFields(log.Fields{"settings": environmentSettings}).Debug("OVERRIDE RANCHER ENVIRONMENT SETTINGS WITH SOME MANUAL VALUE")
	return environmentSettings
}
