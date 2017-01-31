package rancher

import (
	"errors"

	log "github.com/Sirupsen/logrus"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_orchestrate "github.com/wunderkraut/radi-api/operation/orchestrate"
)

/**
 * Orchestrate Up operation for rancher
 */

type RancherOrchestrateUpOperation struct {
	api_orchestrate.BaseOrchestrationUpOperation
	RancherBaseClientOperation
	properties *api_operation.Properties
}

// Alter the ID of the parent operation
func (up *RancherOrchestrateUpOperation) Id() string {
	return "rancher." + up.BaseOrchestrationUpOperation.Id()
}

// Run a validation check on the Operation
func (up *RancherOrchestrateUpOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (up *RancherOrchestrateUpOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	return props
}

// Execute the operation
func (up *RancherOrchestrateUpOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	log.WithFields(log.Fields{"clientsettings": up.RancherClientSettings(), "envsettings": up.RancherEnvironmentSettings()}).Info("SETTINGS")

	result.AddError(errors.New("RANCHER UP OPERATION NOT YET WRITTEN"))
	result.MarkFailed()

	result.MarkFinished()
	return api_operation.Result(result)
}
