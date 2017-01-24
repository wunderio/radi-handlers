package libcompose

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	libCompose_options "github.com/docker/libcompose/project/options"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_orchestrate "github.com/wunderkraut/radi-api/operation/orchestrate"
)

const (
	// config for up orchestration compose settings
	OPERATION_PROPERTY_LIBCOMPOSE_SETTINGS_UP = "compose.up"
)

/**
 * Up specific Properties
 */

// A libcompose Property for net context limiting
type LibcomposeOptionsUpProperty struct {
	value libCompose_options.Up
}

// Id for the Property
func (optionsConf *LibcomposeOptionsUpProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_SETTINGS_UP
}

// Label for the Property
func (optionsConf *LibcomposeOptionsUpProperty) Label() string {
	return "Up operation options"
}

// Description for the Property
func (optionsConf *LibcomposeOptionsUpProperty) Description() string {
	return "Options to configure the Up.  See github.com/docker/libcompose/project/options for more information."
}

// Is the Property internal only
func (optionsConf *LibcomposeOptionsUpProperty) Internal() bool {
	return false
}

// Give an idea of what type of value the property consumes
func (optionsConf *LibcomposeOptionsUpProperty) Type() string {
	return "github.com/docker/libcompose/project/options.Up"
}

func (optionsConf *LibcomposeOptionsUpProperty) Get() interface{} {
	return interface{}(optionsConf.value)
}
func (optionsConf *LibcomposeOptionsUpProperty) Set(value interface{}) bool {
	if converted, ok := value.(libCompose_options.Up); ok {
		optionsConf.value = converted
		return true
	} else {
		log.WithFields(log.Fields{"value": value}).Error("Could not assign Property value, because the passed parameter was the wrong type. Expected github.com/docker/libcompose/project/options.Up")
		return false
	}
}

/**
 * Operation
 */

// Base Up operation
type BaseLibcomposeOrchestrateUpSingleOperation struct {
	properties *api_operation.Properties
}

// Provide static Properties for the operation
func (base *BaseLibcomposeOrchestrateUpSingleOperation) Properties() *api_operation.Properties {
	if base.properties == nil {
		newProperties := &api_operation.Properties{}

		newProperties.Add(api_operation.Property(&LibcomposeOptionsUpProperty{}))

		base.properties = newProperties
	}
	return base.properties
}

// Base Up operation
type BaseLibcomposeOrchestrateUpParametrizedOperation struct {
	properties *api_operation.Properties
}

// Provide static Properties for the operation
func (base *BaseLibcomposeOrchestrateUpParametrizedOperation) Properties() *api_operation.Properties {
	if base.properties == nil {
		newProperties := &api_operation.Properties{}

		newProperties.Add(api_operation.Property(&LibcomposeNoRecreateProperty{}))
		newProperties.Add(api_operation.Property(&LibcomposeForceRecreateProperty{}))
		newProperties.Add(api_operation.Property(&LibcomposeNoBuildProperty{}))
		newProperties.Add(api_operation.Property(&LibcomposeForceRebuildProperty{}))

		base.properties = newProperties
	}
	return base.properties
}

// LibCompose based up orchestrate operation
type LibcomposeOrchestrateUpOperation struct {
	api_orchestrate.BaseOrchestrationUpOperation
	BaseLibcomposeNameFilesOperation
	BaseLibcomposeOrchestrateUpParametrizedOperation

	properties *api_operation.Properties
}

// Validate the libCompose Orchestrate Up operation
func (up *LibcomposeOrchestrateUpOperation) Validate() bool {
	return true
}

// Provide static properties for the operation
func (up *LibcomposeOrchestrateUpOperation) Properties() *api_operation.Properties {
	if up.properties == nil {
		newProperties := &api_operation.Properties{}
		newProperties.Merge(*up.BaseLibcomposeOrchestrateUpParametrizedOperation.Properties())
		newProperties.Merge(*up.BaseLibcomposeNameFilesOperation.Properties())
		up.properties = newProperties
	}
	return up.properties
}

// Execute the libCompose Orchestrate Up operation
func (up *LibcomposeOrchestrateUpOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, nil)

	properties := up.Properties()
	// pass all props to make a project
	project, _ := MakeComposeProject(properties)

	// some props we will use locally

	var netContext context.Context
	var upOptions libCompose_options.Up
	// net context
	if netContextProp, found := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT); found {
		netContext = netContextProp.Get().(context.Context)
	} else {
		result.Set(false, []error{errors.New("Libcompose up operation is missing the context property")})
	}

	// up options
	upOptions = libCompose_options.Up{}
	if upOptionsProp, found := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_NORECREATE); found {
		upOptions.NoRecreate = upOptionsProp.Get().(bool)
	}
	if upOptionsProp, found := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_FORCERECREATE); found {
		upOptions.ForceRecreate = upOptionsProp.Get().(bool)
	}
	if upOptionsProp, found := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_NOBUILD); found {
		upOptions.NoBuild = upOptionsProp.Get().(bool)
	}
	if upOptionsProp, found := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_FORCEREBUILD); found {
		upOptions.ForceBuild = upOptionsProp.Get().(bool)
	}

	if success, _ := result.Success(); success {
		if err := project.APIProject.Up(netContext, upOptions); err != nil {
			result.Set(false, []error{err})
		}
	}

	return api_operation.Result(&result)
}
