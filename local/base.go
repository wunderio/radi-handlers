package local

import (
	api_config "github.com/wunderkraut/radi-api/operation/config"
	api_setting "github.com/wunderkraut/radi-api/operation/setting"
	api_result "github.com/wunderkraut/radi-api/result"
)

/**
 * Local handlers provides operations based entirely on the
 * Local environment, primarily based on config files in
 * a project, based on the current path'
 */

// Constructor for a localHandlerBase
func New_LocalHandler_Base(settings *LocalAPISettings) *LocalHandler_Base {
	return &LocalHandler_Base{
		settings: settings,
	}
}

// A handler for base local handlers
type LocalHandler_Base struct {
	settings *LocalAPISettings
}

// Validate the handler
func (base *LocalHandler_Base) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// An accessor for the local settings
func (base *LocalHandler_Base) LocalAPISettings() *LocalAPISettings {
	return base.settings
}

// A handler for base local handlers that use a config source (like a yml file)
type LocalHandler_ConfigWrapperBase struct {
	configWrapper api_config.ConfigWrapper
}

// An accessor for the ConfigBase ConfigWrapper
func (base *LocalHandler_ConfigWrapperBase) ConfigWrapper() api_config.ConfigWrapper {
	return base.configWrapper
}

// An accessor to set the ConfigBase ConfigWrapper
func (base *LocalHandler_ConfigWrapperBase) SetConfigWrapper(configWrapper api_config.ConfigWrapper) {
	base.configWrapper = configWrapper
}

// A handler for local settings
type LocalHandler_SettingWrapperBase struct {
	settingWrapper api_setting.SettingWrapper
}

// An accessor for the SettingBase SettingWrapper
func (base *LocalHandler_SettingWrapperBase) SettingWrapper() api_setting.SettingWrapper {
	return base.settingWrapper
}

// An accessor to set the SettingsBase SettingWrapper
func (base *LocalHandler_SettingWrapperBase) SetSettingWrapper(settingWrapper api_setting.SettingWrapper) {
	base.settingWrapper = settingWrapper
}
