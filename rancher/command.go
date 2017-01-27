package rancher

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
)

/**
 * Command Handler using Rancher
 */

// Rancher Command Handler
type RancherCommandHandler struct {
	RancherBaseClientHandler
}

// Initialize and activate the Handler
func (command *RancherCommandHandler) Init() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})


	ops := api_operation.Operations{}

	// ops.Add(api_operation.Operation(&RancherCommandUpOperation{BaseRancherServiceOperation: *baseOperation}))
	// ops.Add(api_operation.Operation(&RancherCommandStopOperation{BaseRancherServiceOperation: *baseOperation}))
	// ops.Add(api_operation.Operation(&RancherCommandDownOperation{BaseRancherServiceOperation: *baseOperation}))

	command.operations = &ops

	return api_operation.Result(&result)
}

// Rturn a string identifier for the Handler (not functionally needed yet)
func (command *RancherCommandHandler) Id() string {
	return "rancher.command"
}
