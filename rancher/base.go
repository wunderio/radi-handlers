package rancher

import (
	rancher_client "github.com/rancher/go-rancher/client"

	api_operation "github.com/wunderkraut/radi-api/operation"
)

/**
 * Some base structs to share upcloud functionality
 */

// Shared base handler
type RancherBaseClientHandler struct {
	configSource RancherConfigSource

	operations *api_operation.Operations
}

// Constructor for RancherBaseClientHandler
func New_RancherBaseClientHandler(configSource RancherConfigSource) *RancherBaseClientHandler {
	return &RancherBaseClientHandler{
		configSource: configSource,
		operations:   &api_operation.Operations{},
	}
}

// Get the operations from the handler
func (base *RancherBaseClientHandler) Operations() *api_operation.Operations {
	if base.operations == nil {
		return &api_operation.Operations{}
	} else {
		return base.operations
	}
}

// Retrieve the base settings
func (base *RancherBaseClientHandler) ConfigSource() RancherConfigSource {
	return base.configSource
}

// Share base operation
type RancherBaseClientOperation struct {
	configSource RancherConfigSource
}

// constructor for configSource RancherConfigSource
func New_RancherBaseClientOperation(configSource RancherConfigSource) *RancherBaseClientOperation {
	return &RancherBaseClientOperation{
		configSource: configSource,
	}
}

// Retrieve a rancher client
func (base *RancherBaseClientOperation) RancherClient() *rancher_client.RancherClient {
	return base.configSource.RancherClient()
}

// Retrieve the base settings
func (base *RancherBaseClientOperation) RancherClientSettings() RancherClientSettings {
	return base.configSource.RancherClientSettings()
}

// Retrieve the base settings
func (base *RancherBaseClientOperation) RancherEnvironmentSettings() RancherEnvironmentSettings {
	return base.configSource.RancherEnvironmentSettings()
}
