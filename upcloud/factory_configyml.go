package upcloud

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"

	upcloud_client "github.com/Jalle19/upcloud-go-sdk/upcloud/client"
	upcloud_service "github.com/Jalle19/upcloud-go-sdk/upcloud/service"

	api_config "github.com/james-nesbitt/kraut-api/operation/config"
)

/**
 * This file provides a ConfigWrapper based tool for reading
 * and creating an UpCloud client/service pair for a project
 */

/**
 * Interpreting build config as yml
 */

// Constructor for UpcloudFactoryConfigWrapperYaml
func New_UpcloudFactoryConfigWrapperYaml(configWrapper api_config.ConfigWrapper) UpcloudFactory {
	return UpcloudFactory(&UpcloudFactoryConfigWrapperYaml{
		configWrapper: configWrapper,
		ymlFactory: Yml_UpcloudFactory{},
	})
}


// A BuilderSettingsConfigWrapper, that interprets build config as yml
type UpcloudFactoryConfigWrapperYaml struct {
	configWrapper api_config.ConfigWrapper
	ymlFactory Yml_UpcloudFactory
}

func (configFactory *UpcloudFactoryConfigWrapperYaml) DefaultScope() string {
	/**
	 * @TODO come up with better scopes, but it has to match local conf path keys
	 * @SEE configconnect/settings_yml.go which has the same issue
	 */
	return "project"
}

func (configFactory *UpcloudFactoryConfigWrapperYaml) safe() {
	if configFactory.ymlFactory.Empty() {
		if err := configFactory.Load(); err != nil {
			log.WithError(err).Error("Could not load build configuration")			
		}
	}
}

// Convert this YML struct into a Client
func (configFactory *UpcloudFactoryConfigWrapperYaml) Client() *upcloud_client.Client {
	configFactory.safe()
	return configFactory.ymlFactory.MakeClient()
}
// Convert this YML struct into a Service
func (configFactory *UpcloudFactoryConfigWrapperYaml) Service() *upcloud_service.Service {
	configFactory.safe()
	return configFactory.ymlFactory.MakeService()
}
// Convert this YML struct into a ServiceWrapper
func (configFactory *UpcloudFactoryConfigWrapperYaml) ServiceWrapper() *UpcloudServiceWrapper {
	configFactory.safe()
	return configFactory.ymlFactory.MakeServiceWrapper()
}

// Retrieve values by parsing bytes from the wrapper
func (configFactory *UpcloudFactoryConfigWrapperYaml) Load() error {
	configFactory.ymlFactory = Yml_UpcloudFactory{} // reset stored settings so that we can repopulate it.

	log.Debug("Loading UpCloud config")

	if sources, err := configFactory.configWrapper.Get(CONFIG_KEY_UPCLOUD); err == nil {
		for _, scope := range sources.Order() {
			scopedSource, _ := sources.Get(scope)

			scopedValues := Yml_UpcloudFactory{} // temporarily hold all settings for a specific scope in this
			if err := yaml.Unmarshal(scopedSource, &scopedValues); err == nil {
				log.WithFields(log.Fields{"YmlUser": scopedValues.User, "scope": scope}).Debug("First UpCloud settings yml")
				configFactory.ymlFactory = scopedValues
			} else {
				log.WithError(err).WithFields(log.Fields{"scope": scope}).Error("Couldn't marshall yml scope for upcloud settings")
			}
		}
		return nil
	} else {
		log.WithError(err).Error("Error loading Upcloud config")
		return err
	}
}

// Save the current values to the wrapper
func (configFactory *UpcloudFactoryConfigWrapperYaml) Save() error {
	/**
	 * @TODO THIS
     */
	return errors.New("UpcloudFactoryConfigWrapperYaml Set operation not yet written.")
}

// A temporary holder of BuildSettings, just for yml parsing (probably not needed even)
type Yml_UpcloudFactory struct {
	scope string

	User string 				`yaml:"User"`
	Password string 	        `yaml:"Password"`
}
// Is this struct populated?
func (ymlFactory *Yml_UpcloudFactory) Empty() bool {
	return ymlFactory.User == ""
}
// Convert this YML struct into a Client
func (ymlFactory *Yml_UpcloudFactory) MakeClient() *upcloud_client.Client {
	return New_UpcloudClientSettings(ymlFactory.User, ymlFactory.Password).Client()
}
// Convert this YML struct into a Service
func (ymlFactory *Yml_UpcloudFactory) MakeService() *upcloud_service.Service {
	client := ymlFactory.MakeClient()
	return New_UpcloudServiceSettings(*client).Service()
}
// Convert this YML struct into a Service
func (ymlFactory *Yml_UpcloudFactory) MakeServiceWrapper() *UpcloudServiceWrapper {
	client := ymlFactory.MakeClient()
	return New_UpcloudServiceSettings(*client).ServiceWrapper()
}
