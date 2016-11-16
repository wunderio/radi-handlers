package upcloud

import (
	api_operation "github.com/james-nesbitt/kraut-api/operation"
)

/**
 * Some base structs which other UpCloud Handler and
 * Operation implementations can include
 */

// Constructor for BaseUpcloudServiceHandler
func New_BaseUpcloudServiceHandler(service *UpcloudServiceWrapper, settings *UpcloudBuilderSettings) *BaseUpcloudServiceHandler {
	return &BaseUpcloudServiceHandler{
		service:  service,
		settings: settings,
	}
}

// Base handler with an upcloud service
type BaseUpcloudServiceHandler struct {
	service    *UpcloudServiceWrapper
	settings   *UpcloudBuilderSettings
	operations *api_operation.Operations
}

// Return the stored operatons
func (base *BaseUpcloudServiceHandler) Operations() *api_operation.Operations {
	return base.operations
}

// Get the service
func (base *BaseUpcloudServiceHandler) ServiceWrapper() *UpcloudServiceWrapper {
	return base.service
}

// Get the settings
func (base *BaseUpcloudServiceHandler) Settings() *UpcloudBuilderSettings {
	return base.settings
}

/**
 * Base operations for Upcloud operations, which
 * allow sharing of Upcloud service across instances
 */

// Constructor for BaseUpcloudServiceOperation
func New_BaseUpcloudServiceOperation(service *UpcloudServiceWrapper, settings *UpcloudBuilderSettings) *BaseUpcloudServiceOperation {
	return &BaseUpcloudServiceOperation{
		service:  service,
		settings: settings,
	}
}

// Base operation with an upcloud service
type BaseUpcloudServiceOperation struct {
	service  *UpcloudServiceWrapper
	settings *UpcloudBuilderSettings
}

// Set the service
func (base *BaseUpcloudServiceOperation) ServiceWrapper() *UpcloudServiceWrapper {
	return base.service
}

// Get the settings
func (base *BaseUpcloudServiceOperation) Settings() *UpcloudBuilderSettings {
	return base.settings
}
