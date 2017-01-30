package rancher

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
)

/**
 * Orchestration Handler using Rancher
 */

// Rancher Orchestration Handler
type RancherOrchestrateHandler struct {
	RancherBaseClientHandler
}

// Initialize and activate the Handler
func (orchestrate *RancherOrchestrateHandler) Init() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	base := New_RancherBaseClientOperation(orchestrate.ConfigSource())

	ops := api_operation.Operations{}

	ops.Add(api_operation.Operation(&RancherOrchestrateUpOperation{RancherBaseClientOperation: *base}))
	ops.Add(api_operation.Operation(&RancherOrchestrateDownOperation{RancherBaseClientOperation: *base}))

	orchestrate.operations = &ops

	return api_operation.Result(&result)
}

// Rturn a string identifier for the Handler (not functionally needed yet)
func (orchestrate *RancherOrchestrateHandler) Id() string {
	return "rancher.orchestrate"
}
