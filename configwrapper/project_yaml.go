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

// A ProjectComponentsConfigWRapper, that interprets build config as yml
type ProjectComponentsConfigWrapperYaml struct {
	configWrapper api_config.ConfigWrapper
	components    api_builder.ProjectComponents
}

// Constructor for ProjectComponentsConfigWrapperYaml
func New_ProjectComponentsConfigWrapperYaml(configWrapper api_config.ConfigWrapper) *ProjectComponentsConfigWrapperYaml {
	return &ProjectComponentsConfigWrapperYaml{
		configWrapper: configWrapper,
		components:    api_builder.ProjectComponents{},
	}
}

// Convert this to a api_builder.ProjectConfigWrapper
func (projectComponents *ProjectComponentsConfigWrapperYaml) ProjectConfigWrapper() api_builder.ProjectConfigWrapper {
	return api_builder.ProjectConfigWrapper(projectComponents)
}

func (projectComponents *ProjectComponentsConfigWrapperYaml) DefaultScope() string {
	/**
	 * @TODO come up with better scopes, but it has to match local conf path keys
	 * @SEE configconnect/settings_yml.go which has the same issue
	 */
	return "project"
}

func (projectComponents *ProjectComponentsConfigWrapperYaml) safe() {
	if &projectComponents.components == nil {
		projectComponents.components = api_builder.ProjectComponents{}
	}
	if projectComponents.components.Empty() {
		if err := projectComponents.Load(); err != nil {
			log.WithError(err).Error("Could not load build configuration")
		}
	}
}
func (projectComponents *ProjectComponentsConfigWrapperYaml) Get(key string) (api_builder.ProjectComponent, bool) {
	projectComponents.safe()
	builder, found := projectComponents.components.Get(key)
	return builder, found == nil
}
func (projectComponents *ProjectComponentsConfigWrapperYaml) Set(key string, values api_builder.ProjectComponent) bool {
	projectComponents.safe()
	projectComponents.components.Set(key, values)
	if err := projectComponents.Save(); err != nil {
		log.WithError(err).Error("Could not save build configuration")
		return false
	}
	return true
}
func (projectComponents *ProjectComponentsConfigWrapperYaml) List() []string {
	projectComponents.safe()
	return projectComponents.components.Order()
}

// Retrieve values by parsing bytes from the wrapper
func (projectComponents *ProjectComponentsConfigWrapperYaml) Load() error {
	projectComponents.components = api_builder.ProjectComponents{} // reset stored settings so that we can repopulate it.

	if sources, err := projectComponents.configWrapper.Get(CONFIG_KEY_BUILDER); err == nil {
		for _, scope := range sources.Order() {
			scopedSource, _ := sources.Get(scope)

			scopedValues := Yml_ProjectDefintion{} // temporarily hold all settings for a specific scope in this
			if err := yaml.Unmarshal(scopedSource, &scopedValues); err == nil {
				for index, values := range scopedValues.Components {
					key := scope + "_" + strconv.Itoa(index) // make a unqique key for this setting
					log.WithFields(log.Fields{"ymlSettings": values, "key": key}).Debug("Each yml")
					projectComponents.components.Set(key, *values.MakeProjectComponent())
				}
				break
			} else {
				log.WithError(err).WithFields(log.Fields{"scope": scope}).Error("Couldn't marshall yml scope")
			}
			log.WithFields(log.Fields{"bytes": string(scopedSource), "values": scopedValues, "settings": projectComponents}).Debug("Project:Config->Load()")
		}
		return nil
	} else {
		log.WithError(err).Error("Error loading config for " + CONFIG_KEY_SETTINGS)
		return err
	}
}

// Save the current values to the wrapper
func (projectComponents *ProjectComponentsConfigWrapperYaml) Save() error {
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
	Type             string                             `yaml:"Type"`
	Implementations  []string                           `yaml:"Implementations"`
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
