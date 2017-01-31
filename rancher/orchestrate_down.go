package rancher

import (
	"errors"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_orchestrate "github.com/wunderkraut/radi-api/operation/orchestrate"
)

/**
 * Orchestrate Down operation for rancher
 */

type RancherOrchestrateDownOperation struct {
	api_orchestrate.BaseOrchestrationDownOperation
	RancherBaseClientOperation
	properties *api_operation.Properties
}

// Alter the ID of the parent operation
func (down *RancherOrchestrateDownOperation) Id() string {
	return "rancher." + down.BaseOrchestrationDownOperation.Id()
}

// Run a validation check on the Operation
func (down *RancherOrchestrateDownOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (down *RancherOrchestrateDownOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	return props
}

// Execute the operation
func (down *RancherOrchestrateDownOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	result.AddError(errors.New("RANCHER DOWN OPERATION NOT YET WRITTEN"))
	result.MarkFailed()

	result.MarkFinished()
	return api_operation.Result(result)
}
