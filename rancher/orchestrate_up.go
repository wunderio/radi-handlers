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
func (up *RancherOrchestrateUpOperation) Properties() *api_operation.Properties {
	if up.properties == nil {
		props := api_operation.Properties{}
		up.properties = &props
	}
	return up.properties
}

// Execute the operation
func (up *RancherOrchestrateUpOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	log.WithFields(log.Fields{"clientsettings": up.RancherClientSettings(), "envsettings": up.RancherEnvironmentSettings()}).Info("SETTINGS")

	result.Set(false, []error{errors.New("RANCHER UP OPERATION NOT YET WRITTEN")})

	return api_operation.Result(&result)
}