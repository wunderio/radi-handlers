package configwrapper

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

// Constructor for BuilderComponentsConfigWrapperYaml
func New_BuilderComponentsConfigWrapperYaml(configWrapper api_config.ConfigWrapper) api_builder.BuilderConfigWrapper {
	return api_builder.BuilderConfigWrapper(&BuilderComponentsConfigWrapperYaml{
		configWrapper:   configWrapper,
		buildComponents: api_builder.BuildComponents{},
	})
}

// A BuilderComponentsConfigWRapper, that interprets build config as yml
type BuilderComponentsConfigWrapperYaml struct {
	configWrapper   api_config.ConfigWrapper
	buildComponents api_builder.BuildComponents
}

func (buildComponents *BuilderComponentsConfigWrapperYaml) DefaultScope() string {
	/**
	 * @TODO come up with better scopes, but it has to match local conf path keys
	 * @SEE configconnect/settings_yml.go which has the same issue
	 */
	return "project"
}

func (buildComponents *BuilderComponentsConfigWrapperYaml) safe() {
	if &buildComponents.buildComponents == nil {
		buildComponents.buildComponents = api_builder.BuildComponents{}
	}
	if buildComponents.buildComponents.Empty() {
		if err := buildComponents.Load(); err != nil {
			log.WithError(err).Error("Could not load build configuration")
		}
	}
}
func (buildComponents *BuilderComponentsConfigWrapperYaml) Get(key string) (api_builder.BuildComponent, bool) {
	buildComponents.safe()
	builder, found := buildComponents.buildComponents.Get(key)
	return builder, found == nil
}
func (buildComponents *BuilderComponentsConfigWrapperYaml) Set(key string, values api_builder.BuildComponent) bool {
	buildComponents.safe()
	buildComponents.buildComponents.Set(key, values)
	if err := buildComponents.Save(); err != nil {
		log.WithError(err).Error("Could not save build configuration")
		return false
	}
	return true
}
func (buildComponents *BuilderComponentsConfigWrapperYaml) List() []string {
	buildComponents.safe()
	return buildComponents.buildComponents.Order()
}

// Retrieve values by parsing bytes from the wrapper
func (buildComponents *BuilderComponentsConfigWrapperYaml) Load() error {
	buildComponents.buildComponents = api_builder.BuildComponents{} // reset stored settings so that we can repopulate it.

	if sources, err := buildComponents.configWrapper.Get(CONFIG_KEY_BUILDER); err == nil {
		for _, scope := range sources.Order() {
			scopedSource, _ := sources.Get(scope)

			scopedValues := Yml_BuildDefintion{} // temporarily hold all settings for a specific scope in this
			if err := yaml.Unmarshal(scopedSource, &scopedValues); err == nil {
				for index, values := range scopedValues.Components {
					key := scope + "_" + strconv.Itoa(index) // make a unqique key for this setting
					log.WithFields(log.Fields{"ymlSettings": values, "key": key}).Debug("Each yml")
					buildComponents.buildComponents.Set(key, *values.MakeBuildComponent())
				}
			} else {
				log.WithError(err).WithFields(log.Fields{"scope": scope}).Error("Couldn't marshall yml scope")
			}
			log.WithFields(log.Fields{"bytes": string(scopedSource), "values": scopedValues, "settings": buildComponents}).Debug("Builder:Config->Load()")
		}
		return nil
	} else {
		log.WithError(err).Error("Error loading config for " + CONFIG_KEY_SETTINGS)
		return err
	}
}

// Save the current values to the wrapper
func (buildComponents *BuilderComponentsConfigWrapperYaml) Save() error {
	/**
	 * @TODO THIS
	 */
	return errors.New("BuilderComponentsConfigWrapperYaml Set operation not yet written.")
}

// A temporary holder for the list of components from the yml components file
type Yml_BuildDefintion struct {
	Components []Yml_BuildComponent `yaml:"Components"`
}

// A temporary holder of BuildComponents, just for yml parsing (probably not needed even)
type Yml_BuildComponent struct {
	Type             string                           `yaml:"Type"`
	Implementations  []string                         `yaml:"Implementations"`
	SettingsProvider Yml_BuildSettingSettingsProvider `yaml:"Settings"`
}

// Convert this YML struct into a proper BuildSetting struct
func (component *Yml_BuildComponent) MakeBuildComponent() *api_builder.BuildComponent {
	return api_builder.New_BuildComponent(component.Type, *api_builder.New_Implementations(component.Implementations), api_builder.SettingsProvider(component.SettingsProvider))
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
