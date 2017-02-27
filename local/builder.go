package local

import (
	log "github.com/Sirupsen/logrus"

	api_api "github.com/wunderkraut/radi-api/api"
	api_builder "github.com/wunderkraut/radi-api/builder"
	api_handler "github.com/wunderkraut/radi-api/handler"
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_result "github.com/wunderkraut/radi-api/result"

	api_config "github.com/wunderkraut/radi-api/operation/config"
	api_security "github.com/wunderkraut/radi-api/operation/security"
	api_setting "github.com/wunderkraut/radi-api/operation/setting"
)

/**
 * Provide a builder for the local API
 *
 * The builder works primarily by creating and initializaing
 * the Handlers that are defined in the other files.
 */

// Provide a handler for building all local operations
type LocalBuilder struct {
	settings LocalAPISettings

	parent   api_api.API
	handlers api_handler.Handlers

	common_base *LocalHandler_Base

	Config   api_config.ConfigWrapper
	Setting  api_setting.SettingWrapper
	Security api_security.SecurityWrapper
}

// Constructor for LocalBuilder
func New_LocalBuilder(settings LocalAPISettings) *LocalBuilder {
	return &LocalBuilder{
		settings: settings,
		handlers: api_handler.New_SimpleHandlers().Handlers(),
	}
}

// Convert this to a builder
func (builder *LocalBuilder) Builder() api_builder.Builder {
	return api_builder.Builder(builder)
}

// Builder ID
func (builder *LocalBuilder) Id() string {
	return "local"
}

// Set the parent API, which may need to build Config and Setting Wrappers
func (builder *LocalBuilder) SetAPI(parent api_api.API) {
	builder.parent = parent
}

// Get the list of operations for the implementations
func (builder *LocalBuilder) AddHandler(hand api_handler.Handler) {
	builder.handlers.Add(hand)
}

// Initialize the handler for certain implementations
func (builder *LocalBuilder) Activate(implementations api_builder.Implementations, settingsProvider api_builder.SettingsProvider) api_result.Result {
	for _, implementation := range implementations.Order() {
		switch implementation {
		case "config":
			builder.build_Config()
		case "setting":
			builder.build_Setting()
		case "project":
			builder.build_Project()
		case "security":
			builder.build_Security()

		default:
			log.WithFields(log.Fields{"implementation": implementation}).Error("Local builder implementation not available")
		}
	}

	return api_result.MakeSuccessfulResult()
}

// Validate the builder after Activation is complete
func (builder *LocalBuilder) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Get the list of operations for the implementations
func (builder *LocalBuilder) Operations() api_operation.Operations {
	return builder.handlers.Operations()
}

// Builder Settings
func (builder *LocalBuilder) LocalAPISettings() LocalAPISettings {
	return builder.settings
}

// Create a shareable common base
func (builder *LocalBuilder) Base() *LocalHandler_Base {
	if builder.common_base == nil {

		log.Debug("Building new base handler")

		// Create a base handler, which just wraps the settings
		builder.common_base = New_LocalHandler_Base(&builder.settings)
	}
	return builder.common_base
}

// Add local Handlers for Config and Settings
func (builder *LocalBuilder) build_Config() api_result.Result {
	// Build a config whandler
	local_config := LocalHandler_Config{
		LocalHandler_Base: *builder.Base(),
	}

	res := local_config.Validate()
	<-res.Finished()

	if res.Success() {
		builder.AddHandler(api_handler.Handler(&local_config))
		// Get a config wrapper for other handlers
		builder.Config = local_config.ConfigWrapper()

		log.WithFields(log.Fields{"ConfigWrapper": builder.Config}).Debug("localBuilder: Built Config Handler")
	}

	return res
}

// Add local Handlers for Setting
func (builder *LocalBuilder) build_Setting() api_result.Result {

	// Build a settings handler which uses the configwrapper and the base
	local_setting := LocalHandler_Setting{
		LocalHandler_Base: *builder.Base(),
	}
	local_setting.SetConfigWrapper(builder.Config)

	res := local_setting.Validate()
	<-res.Finished()

	if res.Success() {
		builder.AddHandler(api_handler.Handler(&local_setting))
		// Get a settings wrapper for other handlers
		builder.Setting = local_setting.SettingWrapper()

		log.WithFields(log.Fields{"SettingWrapper": builder.Setting}).Debug("localBuilder: Built Setting Handler")
	}

	return res
}

// Make a Local based API object for "no existing project found" to allow for project operations
func (builder *LocalBuilder) build_Project() api_result.Result {
	// Build a config wrapper using the base for settings
	local_project := LocalHandler_Project{
		LocalHandler_Base: *builder.Base(),
	}

	res := local_project.Validate()
	<-res.Finished()

	if res.Success() {
		builder.AddHandler(api_handler.Handler(&local_project))

		log.Debug("localBuilder: Built Project Handler")
	}

	return res
}

// Add local Handlers for Security operations
func (builder *LocalBuilder) build_Security() api_result.Result {
	// Build a command Handler
	local_security := LocalHandler_Security{
		LocalHandler_Base: *builder.Base(),
	}
	local_security.SetConfigWrapper(builder.Config)

	res := local_security.Validate()
	<-res.Finished()

	if res.Success() {
		builder.AddHandler(api_handler.Handler(&local_security))
		// Get an orchestrate wrapper for other handlers
		builder.Security = local_security.SecurityWrapper()

		log.WithFields(log.Fields{"CSecurityWrapper": builder.Security}).Debug("localBuilder: Built Security Handler")
	}

	return res
}
