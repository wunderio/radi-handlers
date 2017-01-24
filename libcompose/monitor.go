package libcompose

import (
	"errors"
	"io"
	"os"

	"golang.org/x/net/context"
	log "github.com/Sirupsen/logrus"	

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_monitor "github.com/wunderkraut/radi-api/operation/monitor"
)

/**
 * Monitoring operations provided by the libcompose
 * handler.
 */

const (
	OPERATION_ID_COMPOSE_MONITOR_LOGS = api_monitor.OPERATION_ID_MONITOR_LOGS + ".compose"

	OPERATION_ID_COMPOSE_MONITOR_PS = "libcompose.monitor.ps"
)

// An operations which streams the container logs from libcompose
type LibcomposeMonitorLogsOperation struct {
	api_monitor.BaseMonitorLogsOperation
	BaseLibcomposeNameFilesOperation

	properties *api_operation.Properties
}

// Use a different Id() than the parent
func (logs *LibcomposeMonitorLogsOperation) Id() string {
	return OPERATION_ID_COMPOSE_MONITOR_LOGS
}

// Validate
func (logs *LibcomposeMonitorLogsOperation) Validate() bool {
	return true
}

// Provide static properties for the operation
func (logs *LibcomposeMonitorLogsOperation) Properties() *api_operation.Properties {
	if logs.properties == nil {
		newProperties := &api_operation.Properties{}
		newProperties.Add(&LibcomposeDetachProperty{})
		newProperties.Merge(*logs.BaseLibcomposeNameFilesOperation.Properties())
		logs.properties = newProperties
	}
	return logs.properties
}

// Execute the libCompose monitor logs operation
func (logs *LibcomposeMonitorLogsOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, nil)

	properties := logs.Properties()
	// pass all confs to make a project
	project, _ := MakeComposeProject(properties)

	// some confs we will use locally

	var netContext context.Context
	// net context
	if netContextProp, found := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT); found {
		netContext = netContextProp.Get().(context.Context)
	} else {
		result.Set(false, []error{errors.New("Libcompose up operation is missing the context property")})
	}

	var follow bool
	// follow conf
	if followProp, found := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_DETACH); found {
		follow = !followProp.Get().(bool)
	} else {
		result.Set(true, []error{errors.New("Libcompose logs operation is missing the detach property")})
	}

	// output handling test
	if outputProp, found := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_OUTPUT); found {
		outputProp.Set(io.Writer(os.Stdout))
	}

	if success, _ := result.Success(); success {
		if err := project.APIProject.Log(netContext, follow); err != nil {
			result.Set(false, []error{err, errors.New("Could not attach to the project for logs")})
		}
	}

	return api_operation.Result(&result)
}


// LibCompose based ps orchestrate operation
type LibcomposeOrchestratePsOperation struct {
	BaseLibcomposeNameFilesOperation

	properties *api_operation.Properties
}

// Label the operation
func (ps *LibcomposeOrchestratePsOperation) Id() string {
	return "libcompose.monitor.ps"
}

// Label the operation
func (ps *LibcomposeOrchestratePsOperation) Label() string {
	return "Ps"
}

// Description for the operation
func (ps *LibcomposeOrchestratePsOperation) Description() string {
	return "This operation will list all containers."
}

// Is this an internal API operation
func (ps *LibcomposeOrchestratePsOperation) Internal() bool {
	return false
}

// Validate the libCompose Orchestrate Ps operation
func (ps *LibcomposeOrchestratePsOperation) Validate() bool {
	return true
}

// Provide static properties for the operation
func (ps *LibcomposeOrchestratePsOperation) Properties() *api_operation.Properties {
	if ps.properties == nil {
		newProperties := &api_operation.Properties{}
		newProperties.Merge(*ps.BaseLibcomposeNameFilesOperation.Properties())
		ps.properties = newProperties
	}
	return ps.properties
}

// Execute the libCompose Orchestrate Ps operation
func (ps *LibcomposeOrchestratePsOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, nil)

	properties := ps.Properties()
	// pass all props to make a project
	project, _ := MakeComposeProject(properties)

	// some props we will use locally

	var netContext context.Context

	// net context
	if netContextProp, found := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT); found {
		netContext = netContextProp.Get().(context.Context)
	} else {
		result.Set(false, []error{errors.New("Libcompose ps operation is missing the context property")})
	}

	if success, _ := result.Success(); success {
		if infoset, err := project.APIProject.Ps(netContext); err == nil {
			if len(infoset) == 0 {
				log.Info("No running containers found.")
			} else {
				for index, info := range infoset {
					id, _ := info["Id"]
					name, _ := info["Name"]
					state, _ := info["State"]
					log.WithFields(log.Fields{"index": index, "id": id, "name": name, "state": state, "info": info}).Info("Compose info")
				}
			}
		} else {
			result.Set(false, []error{err})
		}
	}

	return api_operation.Result(&result)
}

