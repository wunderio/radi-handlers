package upcloud

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"

	upcloud_request "github.com/Jalle19/upcloud-go-sdk/upcloud/request"

	api_operation "github.com/james-nesbitt/kraut-api/operation"
)

/**
 * Here are a number of server provision related operations, none
 * of which are public, but all of which are used together in the
 * more public provision operations.
 */

/**
 * HANDLER
 *
 * @Note that this handler is not typically needed as it would
 * only add internal operations.  For upcloud provision ops, the
 * tend to build the related operations directly.
 */

// UpCloud Provisioning Handler
type UpcloudServerHandler struct {
	BaseUpcloudServiceHandler
}

// Initialize and activate the Handler
func (server *UpcloudServerHandler) Init() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	baseOperation := server.BaseUpcloudServiceOperation()

	ops := api_operation.Operations{}

	ops.Add(api_operation.Operation(&UpcloudServerCreateOperation{BaseUpcloudServiceOperation: *baseOperation}))
	ops.Add(api_operation.Operation(&UpcloudServerStopOperation{BaseUpcloudServiceOperation: *baseOperation}))
	ops.Add(api_operation.Operation(&UpcloudServerDeleteOperation{BaseUpcloudServiceOperation: *baseOperation}))

	server.operations = &ops

	return api_operation.Result(&result)
}

// Rturn a string identifier for the Handler (not functionally needed yet)
func (server *UpcloudServerHandler) Id() string {
	return "upcloud.server"
}

/**
 * OPERATIONS
 */

// Create a new server operation
type UpcloudServerCreateOperation struct {
	BaseUpcloudServiceOperation
	properties *api_operation.Properties
}

// Return the string machinename/id of the Operation
func (create *UpcloudServerCreateOperation) Id() string {
	return "upcloud.server.create"
}

// Return a user readable string label for the Operation
func (create *UpcloudServerCreateOperation) Label() string {
	return "Create UpCloud server"
}

// return a multiline string description for the Operation
func (create *UpcloudServerCreateOperation) Description() string {
	return "Create an UpCloud server for this project."
}

// Run a validation check on the Operation
func (create *UpcloudServerCreateOperation) Validate() bool {
	return true
}

// Is this operation an internal Operation
func (create *UpcloudServerCreateOperation) Internal() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (create *UpcloudServerCreateOperation) Properties() *api_operation.Properties {
	if create.properties == nil {
		props := api_operation.Properties{}

		props.Add(api_operation.Property(&UpcloudServerCreateRequestProperty{}))
		props.Add(api_operation.Property(&UpcloudServerDetailsProperty{}))

		create.properties = &props
	}
	return create.properties
}

// Execute the Operation
/**
 * @note this is a first version of the operation.  It does not implement
 *   the following checks/functionality:
 *     1. are the servies already provisioned?
 *     2. get the servers defintions from settings
 */
func (create *UpcloudServerCreateOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	service := create.ServiceWrapper()
	// settings := create.BuilderSettings()

	log.Info("Provisioning project server on Upcloud")

	properties := create.Properties()

	request := upcloud_request.CreateServerRequest{}
	if requestProp, found := properties.Get(UPCLOUD_SERVER_CREATEREQUEST_PROPERTY); found {
		request = requestProp.Get().(upcloud_request.CreateServerRequest)
		log.WithFields(log.Fields{"key": UPCLOUD_SERVER_CREATEREQUEST_PROPERTY, "prop": requestProp, "value": request}).Debug("Retrieved create server request")
	}

	serverDetails, err := service.CreateServer(&request)

	if err == nil {

		if detailsProp, found := properties.Get(UPCLOUD_SERVER_DETAILS_PROPERTY); found {
			detailsProp.Set(*serverDetails)
		}
		log.WithFields(log.Fields{"UUID": serverDetails.UUID}).Info("Server created")

	} else {
		result.Set(false, []error{err, errors.New("Unable to provision new server.")})
	}

	return api_operation.Result(&result)
}

// Delete a server operation
type UpcloudServerDeleteOperation struct {
	BaseUpcloudServiceOperation
	properties *api_operation.Properties
}

// Return the string machinename/id of the Operation
func (delete *UpcloudServerDeleteOperation) Id() string {
	return "upcloud.server.delete"
}

// Return a user readable string label for the Operation
func (delete *UpcloudServerDeleteOperation) Label() string {
	return "Remove UpCloud servers"
}

// return a multiline string description for the Operation
func (delete *UpcloudServerDeleteOperation) Description() string {
	return "Remove UpCloud servers for this project."
}

// Run a validation check on the Operation
func (delete *UpcloudServerDeleteOperation) Validate() bool {
	return true
}

// Is this operation an internal Operation
func (delete *UpcloudServerDeleteOperation) Internal() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (delete *UpcloudServerDeleteOperation) Properties() *api_operation.Properties {
	if delete.properties == nil {
		props := api_operation.Properties{}

		props.Add(api_operation.Property(&UpcloudWaitProperty{}))
		props.Add(api_operation.Property(&UpcloudServerUUIDProperty{}))

		delete.properties = &props
	}
	return delete.properties
}

// Execute the Operation
/**
 * @NOTE this is a first version.
 *
 * We will want to :
 *  1. retrieve servers by tag
 *  2. have a "remove-specific-uuid" option?
 */
func (delete *UpcloudServerDeleteOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	service := delete.ServiceWrapper()
	settings := delete.BuilderSettings()

	properties := delete.Properties()

	global := false
	if globalProp, found := properties.Get(UPCLOUD_GLOBAL_PROPERTY); found {
		global = globalProp.Get().(bool)
		log.WithFields(log.Fields{"key": UPCLOUD_GLOBAL_PROPERTY, "prop": globalProp, "value": global}).Debug("Allowing global access")
	}
	wait := false
	if waitProp, found := properties.Get(UPCLOUD_WAIT_PROPERTY); found {
		wait = waitProp.Get().(bool)
		log.WithFields(log.Fields{"key": UPCLOUD_WAIT_PROPERTY, "prop": waitProp, "value": wait}).Debug("Wait for operation to complete")
	}
	uuidMatch := []string{}
	if uuidProp, found := properties.Get(UPCLOUD_SERVER_UUID_PROPERTY); found {
		newUUIDs := uuidProp.Get().([]string)
		uuidMatch = append(uuidMatch, newUUIDs...)
		log.WithFields(log.Fields{"key": UPCLOUD_SERVER_UUID_PROPERTY, "prop": uuidMatch, "value": uuidMatch}).Debug("Filter: Server UUID")
	}

	if len(uuidMatch) > 0 {

		count := 0
		for _, uuid := range uuidMatch {
			if !(global || settings.ServerUUIDAllowed(uuid)) {
				log.WithFields(log.Fields{"uuid": uuid}).Error("Server UUID not a part of the project. Details will not be shown.")
				continue
			}

			request := upcloud_request.DeleteServerRequest{
				UUID: uuid,
			}

			err := service.DeleteServer(&request)

			if err == nil {
				if wait {
					waitRequest := upcloud_request.WaitForServerStateRequest{
						UUID:           uuid,
						DesiredState:   "stopped",
						UndesiredState: "started",
						Timeout:        time.Duration(60) * time.Second,
					}
					details, err := service.WaitForServerState(&waitRequest)

					if err == nil {
						count++
						log.WithFields(log.Fields{"UUID": uuid, "state": details.State, "progress": details.Progress}).Info("Removed UpCloud server")
					} else {
						result.Set(false, []error{err, errors.New("timeout waiting for server be removed.")})
					}
				} else {
					count++
					log.WithFields(log.Fields{"UUID": uuid}).Info("Removed UpCloud server")
				}
			} else {
				result.Set(false, []error{err, errors.New("Could not remove UpCloud server")})
			}
		}

	} else {
		log.Info("No servers requested.  You should have passed a server UUID") // @TODO remove this when we are tagging servers
	}

	return api_operation.Result(&result)
}

// Provision up operation
type UpcloudServerStopOperation struct {
	BaseUpcloudServiceOperation
	properties *api_operation.Properties
}

// Return the string machinename/id of the Operation
func (stop *UpcloudServerStopOperation) Id() string {
	return "upcloud.server.stop"
}

// Return a user readable string label for the Operation
func (stop *UpcloudServerStopOperation) Label() string {
	return "Stop UpCloud server"
}

// return a multiline string description for the Operation
func (stop *UpcloudServerStopOperation) Description() string {
	return "Stop UpCloud servers."
}

// Run a validation check on the Operation
func (stop *UpcloudServerStopOperation) Validate() bool {
	return true
}

// Is this operation an internal Operation
func (stop *UpcloudServerStopOperation) Internal() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (stop *UpcloudServerStopOperation) Properties() *api_operation.Properties {
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
 */
func (stop *UpcloudServerStopOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	service := stop.ServiceWrapper()
	settings := stop.BuilderSettings()

	properties := stop.Properties()

	global := false
	if globalProp, found := properties.Get(UPCLOUD_GLOBAL_PROPERTY); found {
		global = globalProp.Get().(bool)
		log.WithFields(log.Fields{"key": UPCLOUD_GLOBAL_PROPERTY, "prop": globalProp, "value": global}).Debug("Allowing global access")
	}
	wait := false
	if waitProp, found := properties.Get(UPCLOUD_WAIT_PROPERTY); found {
		wait = waitProp.Get().(bool)
		log.WithFields(log.Fields{"key": UPCLOUD_WAIT_PROPERTY, "prop": waitProp, "value": wait}).Debug("Wait for operation to complete")
	}
	uuidMatch := []string{}
	if uuidProp, found := properties.Get(UPCLOUD_SERVER_UUID_PROPERTY); found {
		newUUIDs := uuidProp.Get().([]string)
		uuidMatch = append(uuidMatch, newUUIDs...)
		log.WithFields(log.Fields{"key": UPCLOUD_SERVER_UUID_PROPERTY, "prop": uuidMatch, "value": uuidMatch}).Debug("Filter: Server UUID")
	}

	if len(uuidMatch) > 0 {

		count := 0
		for _, uuid := range uuidMatch {
			if !(global || settings.ServerUUIDAllowed(uuid)) {
				log.WithFields(log.Fields{"uuid": uuid}).Error("Server UUID not a part of the project. Details will not be shown.")
				continue
			}

			request := upcloud_request.StopServerRequest{
				UUID: uuid,
			}

			log.WithFields(log.Fields{"uuid": uuid}).Info("Stopping server.")
			details, err := service.StopServer(&request)

			if err == nil {
				count++
				if wait {
					waitRequest := upcloud_request.WaitForServerStateRequest{
						UUID:           uuid,
						DesiredState:   "stopped",
						UndesiredState: "started",
						Timeout:        time.Duration(60) * time.Second,
					}
					details, err = service.WaitForServerState(&waitRequest)

					if err == nil {
						log.WithFields(log.Fields{"UUID": uuid, "state": details.State, "progress": details.Progress}).Info("Stopped UpCloud server")
					} else {
						result.Set(false, []error{err, errors.New("timeout waiting for server stop.")})
					}
				} else {
					log.WithFields(log.Fields{"UUID": uuid, "state": details.State, "progress": details.Progress}).Info("Stopped UpCloud server")
				}
			} else {
				result.Set(false, []error{err, errors.New("Could not stop UpCloud server")})
			}
		}

	} else {
		log.Info("No servers requested.  You should have passed a server UUID") // @TODO remove this when we are tagging servers
	}

	return api_operation.Result(&result)
}
