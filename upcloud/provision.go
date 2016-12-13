package upcloud

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"

	upcloud "github.com/Jalle19/upcloud-go-sdk/upcloud"
	upcloud_request "github.com/Jalle19/upcloud-go-sdk/upcloud/request"

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

	service := up.ServiceWrapper()
	// settings := up.BuilderSettings()
	serverDefinitions := up.ServerDefinitions()

	// track which servers we actually create here
	createdServers := map[string]processedServer{}

	for _, id := range serverDefinitions.Order() {
		serverResult := api_operation.BaseResult{}
		serverResult.Set(true, []error{})

		serverDefinition, _ := serverDefinitions.Get(id)
		createRequest := serverDefinition.CreateServerRequest()

		if requestProp, found := createProperties.Get(UPCLOUD_SERVER_CREATEREQUEST_PROPERTY); found {
			requestProp.Set(createRequest)
		}

		log.WithFields(log.Fields{"id": serverDefinition.Id()}).Info("Creating new server")

		createResult := createOp.Exec()
		if success, errs := createResult.Success(); !success {
			errs = append(errs, errors.New("Could not provision new UpCloud server"))
			serverResult.Set(false, errs)
			result.Merge(createResult)
			continue
		} else {

			var createDetails upcloud.ServerDetails
			if detailsProp, found := createProperties.Get(UPCLOUD_SERVER_DETAILS_PROPERTY); found {
				createDetails = detailsProp.Get().(upcloud.ServerDetails)
			}

			uuid := createDetails.UUID

			createdServers[id] = processedServer{
				uuid:       uuid,
				definition: serverDefinition,
				details:    createDetails,
			}

			log.WithFields(log.Fields{"id": serverDefinition.Id(), "UUID": uuid, "state": createDetails.State}).Info("Created new server")
		}

		result.Merge(api_operation.Result(&serverResult))
	}

	firewallOp := UpcloudServerApplyFirewallRulesOperation{BaseUpcloudServiceOperation: up.BaseUpcloudServiceOperation}
	firewallProperties := firewallOp.Properties()

	// process tags and firewall rules
	for _, createdServer := range createdServers {
		serverResult := api_operation.BaseResult{}
		serverResult.Set(true, []error{})

		uuid := createdServer.uuid
		serverDefinition := createdServer.definition

		// Before running anything, give the server a chance to get into the proper state
		log.WithFields(log.Fields{"id": serverDefinition.Id(), "UUID": uuid}).Info("Waiting for new server to start")
		if serverDetails, err := service.WaitForServerState(&upcloud_request.WaitForServerStateRequest{UUID: uuid, UndesiredState: "maintenance", Timeout: time.Minute * 2}); err != nil {
			if serverDetails == nil {
				serverResult.Set(false, []error{err, errors.New("Server failed to start properly : " + uuid)})
			} else {
				serverResult.Set(false, []error{err, errors.New("Server failed to start properly : " + serverDetails.UUID)})
			}
		} else {
			log.WithFields(log.Fields{"state": serverDetails.State, "UUID": serverDetails.UUID}).Info("Server successfully created, now finalizing provisioning")

			serverDefinition := createdServer.definition
			firewallRules := serverDefinition.GetFirewallRules()

			if firewallProp, found := firewallProperties.Get(UPCLOUD_FIREWALL_RULES_PROPERTY); found {
				firewallProp.Set(firewallRules)
			}
			if uuidProp, found := firewallProperties.Get(UPCLOUD_SERVER_UUID_PROPERTY); found {
				uuidProp.Set(uuid)
			}

			firewallResult := firewallOp.Exec()
			if success, _ := firewallResult.Success(); !success {
				result.Merge(firewallResult)
				continue
			}

			// var serverDetails upcloud.ServerDetails
			// if detailsProp, found := createProperties.Get(UPCLOUD_SERVER_DETAILS_PROPERTY); found {
			// 	serverDetails = detailsProp.Get().(upcloud.ServerDetails)
			// }
		}

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

		props.Add(api_operation.Property(&UpcloudForceProperty{}))

		down.properties = &props
	}
	return down.properties
}

// Execute the Operation
func (down *UpcloudProvisionDownOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	downProperties := down.Properties()
	deleteOp := UpcloudServerDeleteOperation{BaseUpcloudServiceOperation: down.BaseUpcloudServiceOperation}
	deleteProperties := deleteOp.Properties()

	// service := down.ServiceWrapper()
	// settings := down.BuilderSettings()
	serverDefinitions := down.ServerDefinitions()

	// collect UUIDs of project servers
	uuids := []string{}
	for _, id := range serverDefinitions.Order() {
		serverResult := api_operation.BaseResult{}
		serverResult.Set(true, []error{})

		serverDefinition, _ := serverDefinitions.Get(id)

		if serverDefinition.IsCreated() {
			uuid, _ := serverDefinition.UUID()
			log.WithFields(log.Fields{"id": id, "uuid": uuid}).Debug("Down: Server added to list")
			uuids = append(uuids, uuid)
		} else {
			log.WithFields(log.Fields{"id": id}).Info("Down: Server has not been created, so it will be skipped")
		}
	}

	if len(uuids) > 0 {

		if uuidsProp, found := deleteProperties.Get(UPCLOUD_SERVER_UUIDS_PROPERTY); found {
			log.WithFields(log.Fields{"uuids": uuids}).Info("DOWN: Using UUIDs")
			uuidsProp.Set(uuids)
		}
		if downForceProp, found := downProperties.Get(UPCLOUD_FORCE_PROPERTY); found {
			if deleteForceProp, found := deleteProperties.Get(UPCLOUD_FORCE_PROPERTY); found {
				if downForceProp.Get().(bool) {
					log.Info("DOWN: Forcing operation")
					deleteForceProp.Set(true)
				}
			}
		}

		log.WithFields(log.Fields{"uuids": uuids}).Info("Downing project servers")
		result.Merge(deleteOp.Exec())

	} else {
		log.Info("No active servers found to take down.")
	}

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

// hold info about a server that we have processed
type processedServer struct {
	uuid       string
	definition ServerDefinition
	details    upcloud.ServerDetails
}
