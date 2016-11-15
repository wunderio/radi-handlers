package upcloud

import (
	log "github.com/Sirupsen/logrus"
	
	api_operation "github.com/james-nesbitt/kraut-api/operation"
)
/**
 * Monitor handler for Upcloud operations
 */ 
type UpcloudMonitorHandler struct {
	BaseUpcloudServiceHandler
}

// Initialize and activate the Handler
func (monitor *UpcloudMonitorHandler) Init() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})
	return api_operation.Result(&result)
}
// Rturn a string identifier for the Handler (not functionally needed yet)
func (monitor *UpcloudMonitorHandler) Id() string {
	return "upcloud.monitor"
}
// Return a list of Operations from the Handler
func (monitor *UpcloudMonitorHandler) Operations() *api_operation.Operations {
	ops := api_operation.Operations{}

	baseOperation := New_BaseUpcloudServiceOperation(monitor.ServiceWrapper())

	ops.Add(api_operation.Operation(UpcloudMonitorListServersOperation{BaseUpcloudServiceOperation: *baseOperation}))

	return &ops
}

/**
 * Monitor operations for UpCloud
 */
type UpcloudMonitorTestOperation struct {
	BaseUpcloudServiceOperation
}
// Return the string machinename/id of the Operation
func (monTest UpcloudMonitorTestOperation) Id() string {
	return "upcloud.monitor.test"
}
// Return a user readable string label for the Operation
func (monTest UpcloudMonitorTestOperation) Label() string {
	return "Test upcloud monitor operation"
}
// return a multiline string description for the Operation
func (monTest UpcloudMonitorTestOperation) Description() string {
	return "Test upcloud monitor operation"
}

// Is this operation meant to be used only inside the API
func (monTest UpcloudMonitorTestOperation) Internal() bool {
	return false
}

// FUNCTIONAL

// Run a validation check on the Operation
func (monTest UpcloudMonitorTestOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (monTest UpcloudMonitorTestOperation) Properties() *api_operation.Properties {
	props := api_operation.Properties{}

	return &props
}

// Execute the Operation
func (monTest UpcloudMonitorTestOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})
	return api_operation.Result(&result)
}

/**
 * Monitor operations for UpCloud
 */
type UpcloudMonitorListServersOperation struct {
	BaseUpcloudServiceOperation
}
// Return the string machinename/id of the Operation
func (list UpcloudMonitorListServersOperation) Id() string {
	return "upcloud.monitor.list.servers"
}
// Return a user readable string label for the Operation
func (list UpcloudMonitorListServersOperation) Label() string {
	return "UpCloud server list"
}
// return a multiline string description for the Operation
func (list UpcloudMonitorListServersOperation) Description() string {
	return "List UpCloud servers used in the project"
}

// Is this operation meant to be used only inside the API
func (list UpcloudMonitorListServersOperation) Internal() bool {
	return false
}

// FUNCTIONAL

// Run a validation check on the Operation
func (list UpcloudMonitorListServersOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (list UpcloudMonitorListServersOperation) Properties() *api_operation.Properties {
	props := api_operation.Properties{}

	return &props
}

// Execute the Operation
func (list UpcloudMonitorListServersOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	service := list.ServiceWrapper()

	servers, err := service.GetServers()
	if err == nil {
		for index, server := range servers.Servers {
			log.WithFields(log.Fields{"index": index, "server": server}).Info("Server")
		}
	} else {
		log.WithError(err).Error("Could not list UpCloud servers")		
	}

	return api_operation.Result(&result)
}
