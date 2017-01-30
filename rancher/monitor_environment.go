package rancher

import (
	log "github.com/Sirupsen/logrus"

	api_operation "github.com/wunderkraut/radi-api/operation"
)

/**
 * Monitor Environment operations
 */

/**
 * List Rancher Environments
 */
type RancherMonitorListEnvironmentsOperation struct {
	RancherBaseClientOperation
	properties *api_operation.Properties
}

// Return the string machinename/id of the Operation
func (listEnvironments *RancherMonitorListEnvironmentsOperation) Id() string {
	return "rancher.monitor.environments.list"
}

// Return a user readable string label for the Operation
func (listEnvironments *RancherMonitorListEnvironmentsOperation) Label() string {
	return "List Rancher Environments"
}

// return a multiline string description for the Operation
func (listEnvironments *RancherMonitorListEnvironmentsOperation) Description() string {
	return "List information about the Rancher Environments."
}

// Is this operation meant to be used only inside the API
func (listEnvironments *RancherMonitorListEnvironmentsOperation) Internal() bool {
	return false
}

// Run a validation check on the Operation
func (listEnvironments *RancherMonitorListEnvironmentsOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (listEnvironments *RancherMonitorListEnvironmentsOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	return props
}

// Execute the Operation
func (listEnvironments *RancherMonitorListEnvironmentsOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	client := listEnvironments.RancherClient()

	// zones, err := service.GetZones()
	// if err == nil {
	// 	for index, zone := range zones.Zones {
	// 		filterOut := false

	// 		// if some server filters were passed, filter out anything not in the passed list
	// 		if len(idMatch) > 0 {
	// 			found := false
	// 			for _, id := range idMatch {
	// 				if id == zone.Id {
	// 					found = true
	// 				}
	// 			}
	// 			filterOut = !found
	// 		}

	// 		if !filterOut {
	// 			log.WithFields(log.Fields{"index": index, "id": zone.Id, "description": zone.Description}).Info("UpCloud zone")
	// 		}
	// 	}
	// } else {
	// 	result.Set(false, []error{err, errors.New("Could not retrieve UpCloud zones information.")})
	// }

	log.WithFields(log.Fields{"client": client}).Info("RancherClient")

	return api_operation.Result(&result)
}
