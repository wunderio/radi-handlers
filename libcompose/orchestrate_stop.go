package libcompose

import (
	"errors"

	// log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_orchestrate "github.com/wunderkraut/radi-api/operation/orchestrate"
)

// Base Stop operation
type BaseLibcomposeOrchestrateStopParametrizedOperation struct{}

// Provide static Properties for the operation
func (base *BaseLibcomposeOrchestrateStopParametrizedOperation) Properties() api_operation.Properties {
	if base.properties == nil {
		newProperties := &api_operation.Properties{}

		newProperties.Add(api_operation.Property(&LibcomposeTimeoutProperty{}))

		base.properties = newProperties
	}
	return base.properties
}

// LibCompose based stop orchestrate operation
type LibcomposeOrchestrateStopOperation struct {
	api_orchestrate.BaseOrchestrationStopOperation
	BaseLibcomposeNameFilesOperation
	BaseLibcomposeOrchestrateStopParametrizedOperation

	properties *api_operation.Properties
}

// Validate the libCompose Orchestrate Stop operation
func (stop *LibcomposeOrchestrateStopOperation) Validate() bool {
	return true
}

// Provide static properties for the operation
func (stop *LibcomposeOrchestrateStopOperation) Properties() *api_operation.Properties {
	if stop.properties == nil {
		newProperties := &api_operation.Properties{}
		newProperties.Merge(*stop.BaseLibcomposeOrchestrateStopParametrizedOperation.Properties())
		newProperties.Merge(*stop.BaseLibcomposeNameFilesOperation.Properties())
		stop.properties = newProperties
	}
	return stop.properties
}

// Execute the libCompose Orchestrate Stop operation
func (stop *LibcomposeOrchestrateStopOperation) Exec() api_operation.Result {
	result := api_operation.StandardResult{}
	result.Set(true, nil)

	properties := stop.Properties()
	// pass all props to make a project
	project, _ := MakeComposeProject(properties)

	// some props we will use locally

	var netContext context.Context

	// net context
	if netContextProp, found := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT); found {
		netContext = netContextProp.Get().(context.Context)
	} else {
		result.Set(false, []error{errors.New("Libcompose stop operation is missing the context property")})
	}

	// stop options
	timeout := 10 // timeout in seconds

	if stopOptionsProp, found := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_TIMEOUT); found {
		timeout = stopOptionsProp.Get().(int)
	}

	if success, _ := result.Success(); success {
		if err := project.APIProject.Stop(netContext, timeout); err != nil {
			result.Set(false, []error{err})
		}
	}

	return api_operation.Result(&result)
}
