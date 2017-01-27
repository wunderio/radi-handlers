package rancher

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"

	rancher_client "github.com/rancher/go-rancher/client"

	api_config "github.com/wunderkraut/radi-api/operation/config"
)

// A BuilderSettingsConfigWrapper, that interprets build config as yml
type RancherConfigSourceYaml struct {
	configWrapper api_config.ConfigWrapper
	scope string
	settings *RancherSettings
	client *rancher_client.RancherClient
}

// Constructor for RancherConfigSourceYaml
func New_RancherConfigSourceYaml(configWrapper api_config.ConfigWrapper) *RancherConfigSourceYaml {
	return &RancherConfigSourceYaml{
		configWrapper: configWrapper,
	}
}

func (source *RancherConfigSourceYaml) Client() *rancher_client.RancherClient {
	if source.client == nil {
		source.safe()
		source.client = MakeRancherClientFromSettings(source.settings)		
	}
	return source.client
}

func (source *RancherConfigSourceYaml) Settings() RancherSettings {
	if source.settings == nil {
		source.safe()	
	}
	return *source.settings
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
			log.WithError(err).Error("Could not load build configuration")
		}
	}
}

// Retrieve values by parsing bytes from the ConfigWrapper
func (source *RancherConfigSourceYaml) Load() error {
	log.Debug("Loading Rancher config")

	source.settings = &RancherSettings{}

	if sources, err := source.configWrapper.Get(CONFIG_KEY_RANCHER); err == nil {
		for _, scope := range sources.Order() {
			scopedSource, _ := sources.Get(scope)

			// empty out this oobject
			source.scope = scope
			source.settings = &RancherSettings{}

			if err := yaml.Unmarshal(scopedSource, source.settings); err == nil {
				log.WithFields(log.Fields{"settings": source.Settings, "scope": source.scope}).Debug("Rancher settings parsed from config yml")
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