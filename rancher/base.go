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
		operations: &api_operation.Operations{},
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
func (base *RancherBaseClientHandler) Client() *rancher_client.RancherClient {
	return base.configSource.Client()
}
// Retrieve the base settings
func (base *RancherBaseClientHandler) Settings() RancherSettings {
	return base.configSource.Settings()
}
