package rancher

import (
	log "github.com/Sirupsen/logrus"

	api_api "github.com/wunderkraut/radi-api/api"
	api_builder "github.com/wunderkraut/radi-api/builder"
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_handler "github.com/wunderkraut/radi-api/handler"
	api_config "github.com/wunderkraut/radi-api/operation/config"
)

/**
 * Handler builder and settings to create rancher handlers
 */

type RancherBuilder struct {
	settings RancherSettings

	parent   api_api.API
	handlers api_handler.Handlers	
}

// Set a API for this Handler
func (builder *RancherBuilder) SetAPI(parent api_api.API) {
	builder.parent = parent
}
// Initialize and activate the Handler
func (builder *RancherBuilder) Activate(implementations api_builder.Implementations, settingsProvider api_builder.SettingsProvider) error {
	if &builder.handlers == nil {
		builder.handlers = api_handler.Handlers{}
	}

	// create a base handler
	base := builder.base_RancherBaseClientHandler(settingsProvider)

	for _, implementation := range implementations.Order() {
		switch implementation {
		case "command":
			commandHandler := RancherCommandHandler{RancherBaseClientHandler: *base}
			commandHandler.Init()
			builder.handlers.Add(api_handler.Handler(&commandHandler))
		case "monitor":
			monitorHandler := RancherMonitorHandler{RancherBaseClientHandler: *base}
			monitorHandler.Init()
			builder.handlers.Add(api_handler.Handler(&monitorHandler))
		case "orchestrate":
			orchestrateHandler := RancherOrchestrateHandler{RancherBaseClientHandler: *base}
			orchestrateHandler.Init()
			builder.handlers.Add(api_handler.Handler(&orchestrateHandler))
		default:
			log.WithFields(log.Fields{"implementation": implementation}).Error("Unknown implementation in Rancher builder")
		}
	}

	return nil
}

// Rturn a string identifier for the Handler (not functionally needed yet)
func (builder *RancherBuilder) Id() string {
	return "rancher"
}

// Return a list of Operations from the Handler
func (builder *RancherBuilder) Operations() *api_operation.Operations {
	ops := builder.handlers.Operations()
	return &ops
}


// Return a shared BaseUpcloudServiceOperation for any operation that needs it
func (builder *RancherBuilder) base_RancherBaseClientHandler(settingsProvider api_builder.SettingsProvider) *RancherBaseClientHandler {
	// Builder a configwrapper, which will be used to build upcloud service structs
	operations := builder.parent.Operations()
	configWrapper := api_config.New_SimpleConfigWrapper(&operations)

	// try to convert the settings to RancherSettings
	builder.settings = RancherSettings{}
	settingsProvider.AssignSettings(&builder.settings)

	// make a new config handler from the configWrapper (for now just assume yaml)
	ymlConfigSource := New_RancherConfigSourceYaml(configWrapper)

	// Make a YAML based config wrapper
	return New_RancherBaseClientHandler(ymlConfigSource)
}