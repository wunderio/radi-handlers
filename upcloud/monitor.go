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

	baseOperation := New_BaseUpcloudServiceOperation(monitor.ServiceWrapper(), monitor.Settings())

	ops := api_operation.Operations{}
	ops.Add(api_operation.Operation(&UpcloudMonitorListZonesOperation{BaseUpcloudServiceOperation: *baseOperation}))
	ops.Add(api_operation.Operation(&UpcloudMonitorListServersOperation{BaseUpcloudServiceOperation: *baseOperation}))
	monitor.operations = &ops

	return api_operation.Result(&result)
}

// Rturn a string identifier for the Handler (not functionally needed yet)
func (monitor *UpcloudMonitorHandler) Id() string {
	return "upcloud.monitor"
}

/**
 * Monitor operations for UpCloud
 */

/**
 * List UpCloud Zones
 */
type UpcloudMonitorListZonesOperation struct {
	BaseUpcloudServiceOperation
	properties *api_operation.Properties
}

// Return the string machinename/id of the Operation
func (listZones *UpcloudMonitorListZonesOperation) Id() string {
	return "upcloud.monitor.list.zones"
}

// Return a user readable string label for the Operation
func (listZones *UpcloudMonitorListZonesOperation) Label() string {
	return "List UpCloud Zones"
}

// return a multiline string description for the Operation
func (listZones *UpcloudMonitorListZonesOperation) Description() string {
	return "List information about the UpCloud zones."
}

// Is this operation meant to be used only inside the API
func (listZones *UpcloudMonitorListZonesOperation) Internal() bool {
	return false
}

// Run a validation check on the Operation
func (listZones *UpcloudMonitorListZonesOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (listZones *UpcloudMonitorListZonesOperation) Properties() *api_operation.Properties {
	if listZones.properties == nil {
		props := api_operation.Properties{}

		props.Add(api_operation.Property(&UpcloudGlobalProperty{}))
		props.Add(api_operation.Property(&UpcloudZoneIdProperty{}))

		listZones.properties = &props
	}
	return listZones.properties
}

// Execute the Operation
func (listZones *UpcloudMonitorListZonesOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	service := listZones.ServiceWrapper()
	settings := listZones.Settings()

	global := false
	properties := listZones.Properties()
	if globalProp, found := properties.Get(UPCLOUD_GLOBAL_PROPERTY); found {
		global = globalProp.Get().(bool)
		log.WithFields(log.Fields{"key": UPCLOUD_GLOBAL_PROPERTY, "prop": globalProp, "value": global}).Debug("Filter: global")
	}
	idMatch := []string{}
	if idProp, found := properties.Get(UPCLOUD_ZONE_ID_PROPERTY); found {
		ids := idProp.Get().([]string)
		idMatch = append(idMatch, ids...)
		log.WithFields(log.Fields{"key": UPCLOUD_ZONE_ID_PROPERTY, "prop": idProp, "value": idMatch}).Debug("Filter: zone id")
	}

	zones, err := service.GetZones()
	if err == nil {
		for index, zone := range zones.Zones {
			filterOut := false

			// filter out zones that are no a part of the current project
			if !global {
				filterOut = !settings.ZoneAllowed(zone)
			}

			// if some server filters were passed, filter out anything not in the passed list
			if len(idMatch) > 0 {
				found := false
				for _, id := range idMatch {
					if id == zone.Id {
						found = true
					}
				}
				filterOut = !found
			}

			if !filterOut {
				log.WithFields(log.Fields{"index": index, "id": zone.Id, "description": zone.Description}).Info("UpCloud zone")
			}
		}
	} else {
		log.WithError(err).Error("Could not retrieve UpCloud zones information.")
	}

	return api_operation.Result(&result)
}

/**
 * Monitor operations for UpCloud
 */
type UpcloudMonitorListServersOperation struct {
	BaseUpcloudServiceOperation
	properties *api_operation.Properties
}

// Return the string machinename/id of the Operation
func (listServers *UpcloudMonitorListServersOperation) Id() string {
	return "upcloud.monitor.list.servers"
}

// Return a user readable string label for the Operation
func (listServers *UpcloudMonitorListServersOperation) Label() string {
	return "UpCloud server list"
}

// return a multiline string description for the Operation
func (listServers *UpcloudMonitorListServersOperation) Description() string {
	return "List UpCloud servers used in the project"
}

// Is this operation meant to be used only inside the API
func (listServers *UpcloudMonitorListServersOperation) Internal() bool {
	return false
}

// Run a validation check on the Operation
func (listServers *UpcloudMonitorListServersOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (listServers *UpcloudMonitorListServersOperation) Properties() *api_operation.Properties {
	if listServers.properties == nil {
		props := api_operation.Properties{}

		props.Add(api_operation.Property(&UpcloudGlobalProperty{}))
		props.Add(api_operation.Property(&UpcloudServerUUIDProperty{}))

		listServers.properties = &props
	}
	return listServers.properties
}

// Execute the Operation
func (listServers *UpcloudMonitorListServersOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	service := listServers.ServiceWrapper()
	settings := listServers.Settings()

	global := false
	properties := listServers.Properties()
	if globalProp, found := properties.Get(UPCLOUD_GLOBAL_PROPERTY); found {
		global = globalProp.Get().(bool)
		log.WithFields(log.Fields{"key": UPCLOUD_GLOBAL_PROPERTY, "prop": globalProp, "value": global}).Debug("GLOBAL")
	}
	uuidMatch := []string{}
	if uuidProp, found := properties.Get(UPCLOUD_SERVER_UUID_PROPERTY); found {
		newUUIDs := uuidProp.Get().([]string)
		uuidMatch = append(uuidMatch, newUUIDs...)
		log.WithFields(log.Fields{"key": UPCLOUD_SERVER_UUID_PROPERTY, "prop": uuidMatch, "value": uuidMatch}).Debug("Filter: Server UUID")
	}

	servers, err := service.GetServers()
	if err == nil {
		for index, server := range servers.Servers {
			filterOut := false

			// filter out servers that are no a part of the current project
			if !global {
				filterOut = !settings.ServerAllowed(server)
			}

			// if some server filters were passed, filter out anything not in the passed list
			if len(uuidMatch) > 0 {
				found := false
				for _, uuid := range uuidMatch {
					if uuid == server.UUID {
						found = true
					}
				}
				filterOut = !found
			}

			if !filterOut {
				log.WithFields(log.Fields{"index": index, "uuid": server.UUID, "title": server.Title, "plan": server.Plan, "zone": server.Zone}).Info("Server")
			}
		}
	} else {
		log.WithError(err).Error("Could not list UpCloud servers")
	}

	return api_operation.Result(&result)
}
