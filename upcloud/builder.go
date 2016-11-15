package upcloud

import (
	log "github.com/Sirupsen/logrus"

	api_api "github.com/james-nesbitt/kraut-api/api"
	api_builder "github.com/james-nesbitt/kraut-api/builder"
	api_handler "github.com/james-nesbitt/kraut-api/handler"
	api_operation "github.com/james-nesbitt/kraut-api/operation"
	api_config "github.com/james-nesbitt/kraut-api/operation/config"
)

/**
 * A kraut builder for upcloud handlers
 */


// Upcloud Builder
type UpcloudBuilder struct {
	parent api_api.API

	handlers *api_handler.Handlers
	base_UpcloudServiceHandler *BaseUpcloudServiceHandler
}

// Set a API for this Handler
func (builder *UpcloudBuilder) SetAPI(parent api_api.API) {
	// Keep that api, so that we can use it to make a ConfigWrapper later on
	builder.parent = parent
}
// Initialize and activate the Handler
func (builder *UpcloudBuilder) Activate(implementations api_builder.Implementations, settings interface{}) error {
	if builder.handlers == nil {
		builder.handlers = &api_handler.Handlers{}
	}

	// This base handler is commonly used in the implementation handlers, so get it once here.
	baseHandler := builder.base_BaseUpcloudServiceOperation()

	for _, implementation := range implementations.Order() {
		switch implementation {
		case "monitor":
			monitorHandler := UpcloudMonitorHandler{BaseUpcloudServiceHandler: *baseHandler}
			builder.handlers.Add(api_handler.Handler(&monitorHandler))
		default:
			log.WithFields(log.Fields{"implementation": implementation}).Error("Unknown implementation in UpCloud builder")
		}
	}

	return nil
}
// Rturn a string identifier for the Handler (not functionally needed yet)
func (builder *UpcloudBuilder) Id() string {
	return "upcloud"
}
// Return a list of Operations from the Handler
func (builder *UpcloudBuilder) Operations() *api_operation.Operations {
	ops := builder.handlers.Operations()
	return &ops
}

// Return a shared BaseUpcloudServiceOperation for any operation that needs it
func (builder *UpcloudBuilder) base_BaseUpcloudServiceOperation() *BaseUpcloudServiceHandler {
	if builder.base_UpcloudServiceHandler == nil {
		// Builder a configwrapper, which will be used to build upcloud service structs
		operations := builder.parent.Operations()
		configWrapper := api_config.New_SimpleConfigWrapper(&operations)
		// get an upcloud factory, using the config wrapper (probably a file like upcloud.yml)
		upcloudFactory := New_UpcloudFactoryConfigWrapperYaml(configWrapper)

		// Ask the factory to build the service wrapper, use that to make the base operations
		serviceWrapper := upcloudFactory.ServiceWrapper()
		builder.base_UpcloudServiceHandler = New_BaseUpcloudServiceHandler(serviceWrapper)
	}
	return builder.base_UpcloudServiceHandler
}
