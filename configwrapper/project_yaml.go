package configwrapper

import (
	"errors"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"

	api_builder "github.com/wunderkraut/radi-api/builder"
	api_config "github.com/wunderkraut/radi-api/operation/config"
)

/**
 * Interpreting build config as yml
 */

// Constructor for ProjectComponentsConfigWrapperYaml
func New_ProjectComponentsConfigWrapperYaml(configWrapper api_config.ConfigWrapper) api_builder.ProjectConfigWrapper {
	return api_builder.ProjectConfigWrapper(&ProjectComponentsConfigWrapperYaml{
		configWrapper:   configWrapper,
		buildComponents: api_builder.ProjectComponents{},
	})
}

// A ProjectComponentsConfigWRapper, that interprets build config as yml
type ProjectComponentsConfigWrapperYaml struct {
	configWrapper   api_config.ConfigWrapper
	buildComponents api_builder.ProjectComponents
}

func (buildComponents *ProjectComponentsConfigWrapperYaml) DefaultScope() string {
	/**
	 * @TODO come up with better scopes, but it has to match local conf path keys
	 * @SEE configconnect/settings_yml.go which has the same issue
	 */
	return "project"
}

func (buildComponents *ProjectComponentsConfigWrapperYaml) safe() {
	if &buildComponents.buildComponents == nil {
		buildComponents.buildComponents = api_builder.ProjectComponents{}
	}
	if buildComponents.buildComponents.Empty() {
		if err := buildComponents.Load(); err != nil {
			log.WithError(err).Error("Could not load build configuration")
		}
	}
}
func (buildComponents *ProjectComponentsConfigWrapperYaml) Get(key string) (api_builder.ProjectComponent, bool) {
	buildComponents.safe()
	builder, found := buildComponents.buildComponents.Get(key)
	return builder, found == nil
}
func (buildComponents *ProjectComponentsConfigWrapperYaml) Set(key string, values api_builder.ProjectComponent) bool {
	buildComponents.safe()
	buildComponents.buildComponents.Set(key, values)
	if err := buildComponents.Save(); err != nil {
		log.WithError(err).Error("Could not save build configuration")
		return false
	}
	return true
}
func (buildComponents *ProjectComponentsConfigWrapperYaml) List() []string {
	buildComponents.safe()
	return buildComponents.buildComponents.Order()
}

// Retrieve values by parsing bytes from the wrapper
func (buildComponents *ProjectComponentsConfigWrapperYaml) Load() error {
	buildComponents.buildComponents = api_builder.ProjectComponents{} // reset stored settings so that we can repopulate it.

	if sources, err := buildComponents.configWrapper.Get(CONFIG_KEY_BUILDER); err == nil {
		for _, scope := range sources.Order() {
			scopedSource, _ := sources.Get(scope)

			scopedValues := Yml_ProjectDefintion{} // temporarily hold all settings for a specific scope in this
			if err := yaml.Unmarshal(scopedSource, &scopedValues); err == nil {
				for index, values := range scopedValues.Components {
					key := scope + "_" + strconv.Itoa(index) // make a unqique key for this setting
					log.WithFields(log.Fields{"ymlSettings": values, "key": key}).Debug("Each yml")
					buildComponents.buildComponents.Set(key, *values.MakeProjectComponent())
				}
				break
			} else {
				log.WithError(err).WithFields(log.Fields{"scope": scope}).Error("Couldn't marshall yml scope")
			}
			log.WithFields(log.Fields{"bytes": string(scopedSource), "values": scopedValues, "settings": buildComponents}).Debug("Project:Config->Load()")
		}
		return nil
	} else {
		log.WithError(err).Error("Error loading config for " + CONFIG_KEY_SETTINGS)
		return err
	}
}

// Save the current values to the wrapper
func (buildComponents *ProjectComponentsConfigWrapperYaml) Save() error {
	/**
	 * @TODO THIS
	 */
	return errors.New("ProjectComponentsConfigWrapperYaml Set operation not yet written.")
}

// A temporary holder for the list of components from the yml components file
type Yml_ProjectDefintion struct {
	Components []Yml_ProjectComponent `yaml:"Components"`
}

// A temporary holder of ProjectComponents, just for yml parsing (probably not needed even)
type Yml_ProjectComponent struct {
	Type             string                           `yaml:"Type"`
	Implementations  []string                         `yaml:"Implementations"`
	SettingsProvider Yml_ProjectSettingSettingsProvider `yaml:"Settings"`
}

// Convert this YML struct into a proper ProjectSetting struct
func (component *Yml_ProjectComponent) MakeProjectComponent() *api_builder.ProjectComponent {
	return api_builder.New_ProjectComponent(component.Type, *api_builder.New_Implementations(component.Implementations), api_builder.SettingsProvider(component.SettingsProvider))
}

// Yml builder SettingProvider
type Yml_ProjectSettingSettingsProvider struct {
	UnMarshaler func(interface{}) error
}

// Yaml custom UnMarshall handler
func (ymlSettingsProvider *Yml_ProjectSettingSettingsProvider) UnmarshalYAML(unmarshal func(interface{}) error) error {
	ymlSettingsProvider.UnMarshaler = unmarshal
	return nil
}

// UnMarshaller function
func (ymlSettingsProvider Yml_ProjectSettingSettingsProvider) AssignSettings(target interface{}) error {
	if ymlSettingsProvider.UnMarshaler != nil {
		return ymlSettingsProvider.UnMarshaler(target)
	}
	return nil
}
