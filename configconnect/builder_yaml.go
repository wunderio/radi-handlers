package configconnect

import (
	"errors"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"

	api_builder "github.com/james-nesbitt/kraut-api/builder"
	api_config "github.com/james-nesbitt/kraut-api/operation/config"
)

/**
 * Interpreting build config as yml
 */

// Constructor for BuilderSettingsConfigWrapperYaml
func New_BuilderSettingsConfigWrapperYaml(configWrapper api_config.ConfigWrapper) api_builder.BuilderConfigWrapper {
	return api_builder.BuilderConfigWrapper(&BuilderSettingsConfigWrapperYaml{
		configWrapper: configWrapper,
		buildSettings: api_builder.BuildSettings{},
	})
}


// A BuilderSettingsConfigWRapper, that interprets build config as yml
type BuilderSettingsConfigWrapperYaml struct {
	configWrapper api_config.ConfigWrapper
	buildSettings api_builder.BuildSettings
}

func (buildSettings *BuilderSettingsConfigWrapperYaml) DefaultScope() string {
	/**
	 * @TODO come up with better scopes, but it has to match local conf path keys
	 * @SEE configconnect/settings_yml.go which has the same issue
	 */
	return "project"
}

func (buildSettings *BuilderSettingsConfigWrapperYaml) safe() {
	if &buildSettings.buildSettings == nil {
		buildSettings.buildSettings = api_builder.BuildSettings{}
	}
	if buildSettings.buildSettings.Empty() {
		if err := buildSettings.Load(); err != nil {
			log.WithError(err).Error("Could not load build configuration")			
		}
	}
}
func (buildSettings *BuilderSettingsConfigWrapperYaml) Get(key string) (api_builder.BuildSetting, bool) {
	buildSettings.safe()
	builder, found := buildSettings.buildSettings.Get(key)
	return builder, found == nil
}
func (buildSettings *BuilderSettingsConfigWrapperYaml) Set(key string, values api_builder.BuildSetting) bool {
	buildSettings.safe()	
	buildSettings.buildSettings.Set(key, values)
	if err := buildSettings.Save(); err != nil {
		log.WithError(err).Error("Could not save build configuration")
		return false
	}
	return true
}
func (buildSettings *BuilderSettingsConfigWrapperYaml) List() []string {
	buildSettings.safe()
	return buildSettings.buildSettings.Order()
}

// Retrieve values by parsing bytes from the wrapper
func (buildSettings *BuilderSettingsConfigWrapperYaml) Load() error {
	buildSettings.buildSettings = api_builder.BuildSettings{} // reset stored settings so that we can repopulate it.

	if sources, err := buildSettings.configWrapper.Get(CONFIG_KEY_BUILDER); err == nil {
		for _, scope := range sources.Order() {
			scopedSource, _ := sources.Get(scope)

			scopedValues := []Yml_BuildSetting{} // temporarily hold all settings for a specific scope in this
			if err := yaml.Unmarshal(scopedSource, &scopedValues); err == nil {
				for index, values := range scopedValues {
					key := scope+"_"+strconv.Itoa(index) // make a unqique key for this setting
					log.WithFields(log.Fields{"ymlSettings": values, "key": key}).Debug("Each yml")
					buildSettings.buildSettings.Set(key, *values.MakeBuildSetting())
				}
			} else {
				log.WithError(err).WithFields(log.Fields{"scope": scope}).Error("Couldn't marshall yml scope")
			}
			log.WithFields(log.Fields{"bytes": string(scopedSource), "values": scopedValues, "settings": buildSettings}).Debug("Builder:Config->Load()")
		}
		return nil
	} else {
		log.WithError(err).Error("Error loading config for " + CONFIG_KEY_SETTINGS)
		return err
	}
}

// Save the current values to the wrapper
func (buildSettings *BuilderSettingsConfigWrapperYaml) Save() error {
	/**
	 * @TODO THIS
     */
	return errors.New("BuilderSettingsConfigWrapperYaml Set operation not yet written.")
}

// A temporary holder of BuildSettings, just for yml parsing (probably not needed even)
type Yml_BuildSetting struct {
	Type string 				                       `yaml:"Type"`
	Implementations []string 	                       `yaml:"Implementations"`
	SettingsProvider Yml_BuildSettingSettingsProvider  `yaml:"Settings"`
}
// Convert this YML struct into a proper BuildSetting struct
func (setting *Yml_BuildSetting) MakeBuildSetting() *api_builder.BuildSetting {
	return api_builder.New_BuildSetting(setting.Type, *api_builder.New_Implementations(setting.Implementations), api_builder.SettingsProvider(setting.SettingsProvider))
}

// Yml builder SettingProvider
type Yml_BuildSettingSettingsProvider struct {
	UnMarshaler func(interface{}) error
}
// Yaml custom UnMarshall handler
func (ymlSettingsProvider *Yml_BuildSettingSettingsProvider) UnmarshalYAML(unmarshal func(interface{}) error) error {
	ymlSettingsProvider.UnMarshaler = unmarshal
	return nil
}

// UnMarshaller function
func (ymlSettingsProvider Yml_BuildSettingSettingsProvider) AssignSettings(target interface{}) error {
	return ymlSettingsProvider.UnMarshaler(target)
}
