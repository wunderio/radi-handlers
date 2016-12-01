package upcloud

import (
	log "github.com/Sirupsen/logrus"

	upcloud "github.com/Jalle19/upcloud-go-sdk/upcloud"

	api_operation "github.com/james-nesbitt/kraut-api/operation"
	api_provision "github.com/james-nesbitt/kraut-api/operation/provision"
)

/**
 * Functionality for provisioning
 */

/**
 * HANDLER
 */

// UpCloud Provisioning Handler
type UpcloudProvisionHandler struct {
	BaseUpcloudServiceHandler
}

// Initialize and activate the Handler
func (provision *UpcloudProvisionHandler) Init() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	baseOperation := provision.BaseUpcloudServiceOperation()

	ops := api_operation.Operations{}

	ops.Add(api_operation.Operation(&UpcloudProvisionUpOperation{BaseUpcloudServiceOperation: *baseOperation}))
	ops.Add(api_operation.Operation(&UpcloudProvisionStopOperation{BaseUpcloudServiceOperation: *baseOperation}))
	ops.Add(api_operation.Operation(&UpcloudProvisionDownOperation{BaseUpcloudServiceOperation: *baseOperation}))

	provision.operations = &ops

	return api_operation.Result(&result)
}

// Rturn a string identifier for the Handler (not functionally needed yet)
func (provision *UpcloudProvisionHandler) Id() string {
	return "upcloud.provision"
}

/**
 * OPERATIONS
 */

// Provision up operation
type UpcloudProvisionUpOperation struct {
	BaseUpcloudServiceOperation
	api_provision.BaseProvisionUpOperation
	properties *api_operation.Properties
}

// Return the string machinename/id of the Operation
func (up *UpcloudProvisionUpOperation) Id() string {
	return "upcloud.provision.up"
}

// Return a user readable string label for the Operation
func (up *UpcloudProvisionUpOperation) Label() string {
	return "Provision UpCloud servers"
}

// return a multiline string description for the Operation
func (up *UpcloudProvisionUpOperation) Description() string {
	return "Provision the UpCloud servers for this project."
}

// Run a validation check on the Operation
func (up *UpcloudProvisionUpOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (up *UpcloudProvisionUpOperation) Properties() *api_operation.Properties {
	if up.properties == nil {
		props := api_operation.Properties{}
		up.properties = &props
	}
	return up.properties
}

/**
 * Execute the Operation
 *
 * The following steps are followed for each server:
 *   1. create the server - then wait for it to be considered running
 *   2. create the firewall rules
 *   3. tag the server
 */
func (up *UpcloudProvisionUpOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	createOp := UpcloudServerCreateOperation{BaseUpcloudServiceOperation: up.BaseUpcloudServiceOperation}
	createProperties := createOp.Properties()

	//service := up.ServiceWrapper()
	// settings := up.BuilderSettings()
	serverDefinitions := up.ServerDefinitions()

	for _, id := range serverDefinitions.Order() {
		serverResult := api_operation.BaseResult{}
		serverResult.Set(true, []error{})

		serverDefinition, _ := serverDefinitions.Get(id)
		createRequest := serverDefinition.CreateServerRequest()

		if requestProp, found := createProperties.Get(UPCLOUD_SERVER_CREATEREQUEST_PROPERTY); found {
			requestProp.Set(createRequest)
		}

		createResult := createOp.Exec()
		if success, _ := createResult.Success(); !success {
			result.Merge(createResult)
			continue
		}

		var serverDetails upcloud.ServerDetails
		if detailsProp, found := createProperties.Get(UPCLOUD_SERVER_DETAILS_PROPERTY); found {
			serverDetails = detailsProp.Get().(upcloud.ServerDetails)
		}

		log.WithFields(log.Fields{"details": serverDetails}).Info("Craeted server")

		result.Merge(api_operation.Result(&serverResult))
	}

	return api_operation.Result(&result)
}

// Provision up operation
type UpcloudProvisionDownOperation struct {
	BaseUpcloudServiceOperation
	api_provision.BaseProvisionDownOperation
	properties *api_operation.Properties
}

// Return the string machinename/id of the Operation
func (down *UpcloudProvisionDownOperation) Id() string {
	return "upcloud.provision.down"
}

// Return a user readable string label for the Operation
func (down *UpcloudProvisionDownOperation) Label() string {
	return "Remove UpCloud servers"
}

// return a multiline string description for the Operation
func (down *UpcloudProvisionDownOperation) Description() string {
	return "Remove the UpCloud servers for this project."
}

// Run a validation check on the Operation
func (down *UpcloudProvisionDownOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (down *UpcloudProvisionDownOperation) Properties() *api_operation.Properties {
	if down.properties == nil {
		props := api_operation.Properties{}

		props.Add(api_operation.Property(&UpcloudGlobalProperty{}))
		//props.Add(api_operation.Property(&UpcloudWaitProperty{})) // This actually doesn't work out well
		props.Add(api_operation.Property(&UpcloudServerUUIDProperty{}))

		down.properties = &props
	}
	return down.properties
}

// Execute the Operation
/**
 * @NOTE this is a first version.
 *
 * We will want to :
 *  1. retrieve servers by tag
 *  2. have a "remove-specific-uuid" option?
 */
func (down *UpcloudProvisionDownOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	return api_operation.Result(&result)
}

// Provision up operation
type UpcloudProvisionStopOperation struct {
	BaseUpcloudServiceOperation
	api_provision.BaseProvisionStopOperation
	properties *api_operation.Properties
}

// Return the string machinename/id of the Operation
func (stop *UpcloudProvisionStopOperation) Id() string {
	return "upcloud.provision.stop"
}

// Return a user readable string label for the Operation
func (stop *UpcloudProvisionStopOperation) Label() string {
	return "Stop UpCloud servers"
}

// return a multiline string description for the Operation
func (stop *UpcloudProvisionStopOperation) Description() string {
	return "Stop the UpCloud servers for this project."
}

// Run a validation check on the Operation
func (stop *UpcloudProvisionStopOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (stop *UpcloudProvisionStopOperation) Properties() *api_operation.Properties {
	if stop.properties == nil {
		props := api_operation.Properties{}

		props.Add(api_operation.Property(&UpcloudGlobalProperty{}))
		props.Add(api_operation.Property(&UpcloudWaitProperty{}))
		props.Add(api_operation.Property(&UpcloudServerUUIDProperty{}))

		stop.properties = &props
	}
	return stop.properties
}

// Execute the Operation
/**
 * @NOTE this is a first version.
 *
 * We will want to :
 *  1. retrieve servers by tag
 *  2. have a "remove-specific-uuid" option?
 */
func (stop *UpcloudProvisionStopOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	return api_operation.Result(&result)
}
