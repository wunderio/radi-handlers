package libcompose

import (
	"errors"

	// log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_orchestrate "github.com/wunderkraut/radi-api/operation/orchestrate"
)

// Base Start operation
type BaseLibcomposeOrchestrateStartParametrizedOperation struct {
	properties *api_operation.Properties
}

// Provide static Properties for the operation
func (base *BaseLibcomposeOrchestrateStartParametrizedOperation) Properties() *api_operation.Properties {
	if base.properties == nil {
		newProperties := &api_operation.Properties{}

		base.properties = newProperties
	}
	return base.properties
}

// LibCompose based start orchestrate operation
type LibcomposeOrchestrateStartOperation struct {
	api_orchestrate.BaseOrchestrationStartOperation
	BaseLibcomposeNameFilesOperation
	BaseLibcomposeOrchestrateStartParametrizedOperation

	properties *api_operation.Properties
}

// Validate the libCompose Orchestrate Start operation
func (start *LibcomposeOrchestrateStartOperation) Validate() bool {
	return true
}

// Provide static properties for the operation
func (start *LibcomposeOrchestrateStartOperation) Properties() *api_operation.Properties {
	if start.properties == nil {
		newProperties := &api_operation.Properties{}
		newProperties.Merge(*start.BaseLibcomposeOrchestrateStartParametrizedOperation.Properties())
		newProperties.Merge(*start.BaseLibcomposeNameFilesOperation.Properties())
		start.properties = newProperties
	}
	return start.properties
}

// Execute the libCompose Orchestrate Start operation
func (start *LibcomposeOrchestrateStartOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, nil)

	properties := start.Properties()
	// pass all props to make a project
	project, _ := MakeComposeProject(properties)

	// some props we will use locally

	var netContext context.Context

	// net context
	if netContextProp, found := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT); found {
		netContext = netContextProp.Get().(context.Context)
	} else {
		result.Set(false, []error{errors.New("Libcompose start operation is missing the context property")})
	}

	if success, _ := result.Success(); success {
		if err := project.APIProject.Start(netContext); err != nil {
			result.Set(false, []error{err})
		}
	}

	return api_operation.Result(&result)
}
