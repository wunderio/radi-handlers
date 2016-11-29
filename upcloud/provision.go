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

	baseOperation := New_BaseUpcloudServiceOperation(provision.ServiceWrapper(), provision.Settings())

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

		// props.Add(api_operation.Property(&UpcloudGlobalProperty{}))
		// props.Add(api_operation.Property(&UpcloudZoneIdProperty{}))

		up.properties = &props
	}
	return up.properties
}

// Execute the Operation
/**
 * @note this is a first version of the operation.  It does not implement
 *   the following checks/functionality:
 *     1. are the servies already provisioned?
 *     2. get the servers defintions from settings
 */
func (up *UpcloudProvisionUpOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	service := up.ServiceWrapper()
	// settings := up.Settings()

	log.Info("Provisioning project server on Upcloud")

	prov_project := "krauttest"
	prov_zone := "fi-hel1"

	prov_initscript := "" // Initialize script. Can be a URL

	prov_user := upcloud_request.LoginUser{
		CreatePassword: "no", // Allow SSH only with key
		Username:       "kraut",
		SSHKeys: []string{
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDEI7j4EaK2RRKgp7rA9gDIL279WtNBWsPQwKn6YNjb7i1EUAM+IYzdQbPgpYr0rMx67DhvbK1pBeL0HTXfk1ZnSbbZe2xktk+YJo6l8zQ7wYydWMcCcB5HUvgG1/ugTj6wxImYAx7sEuXY4MVO7aHmfMnjV+7Re0uXHjAPL9k5O2Xvt75RmrgG8YpE6MvZtYTzIRmINbuSAX9CWKi46ZuRNYKDyZTSarqA1TOaGKiO6vf2dM7bWqsvitOxwEC6Z0c5nIAjcAGhg+yBEEloWTqNqkxPzbh0AIIO9HjGlnbSaIffwrv78UzrHatukUQOcsg6PBvMPvhjdoB0JrscLneDy0DhF6ptAQporg3SieypB/3hiZ0RfT94c35DQufFphfsphIBXIsqENJKR383sz57PPDtVgyXKSu5ujhXUPgC1rwldGUqVtMLsvC4tcnOIbOK917QjUQ+8cJoC08XBUG10knUoIWP8ebv55sfnBTHW27g+4B1V6ub3Zyq/ZRzeJXWzSh1QmOUXp1Q47bEz3eT2/VRtKyUYAo3ChvceMSERsVtwfRgIcAreIqGy2GJQPe7NkYOAPrirwhfppoFJ/nx3xGFjg0iZIg4Z1nTpwEWBWcC5eo/yfORnhdAooJWRYO37nOrjryUZJsRbFC/Uj7JOIX2QrZEX1bm4SwgSF8tTQ==", // JRN
		},
	}
	prov_networks := []upcloud_request.CreateServerIPAddress{
		upcloud_request.CreateServerIPAddress{
			Access: "private",
			Family: "IPv4",
		},
		upcloud_request.CreateServerIPAddress{
			Access: "public",
			Family: "IPv4",
		},
		upcloud_request.CreateServerIPAddress{
			Access: "public",
			Family: "IPv6",
		},
	}
	prov_storages := []upcloud.CreateServerStorageDevice{
		upcloud.CreateServerStorageDevice{ // primary disk
			Action:  "clone",
			Storage: "01000000-0000-4000-8000-000080010200", //size=5 title="CoreOS Stable 1068.8.0" type=template
			Title:   "coreos-install",
			Size:    10, // Storage size in gigabytes, if cloning it has to be larger the source size
			Tier:    "maxiops",
		},
		// upcloud.CreateServerStorageDevice{ // primary disk
		// 	Action:  "create",
		// 	Address: "virtio:0",
		// 	Title:   "coreos-root",
		// 	Size:    10,
		// 	Tier:    "maxiops",
		// 	Type:    "disk",
		// },
	}
	prov_firewall_rules := []upcloud.FirewallRule{
		// upcloud.FirewallRule{
		// 	Action:                  "accept",
		// 	Comment:                 "Alow HTTP from anywhere",
		// 	DestinationAddressEnd:   "",
		// 	DestinationAddressStart: "",
		// 	DestinationPortEnd:      "80",
		// 	DestinationPortStart:    "80",
		// 	Direction:               "in",
		// 	Family:                  "IPv4",
		// 	ICMPType:                "",
		// 	Position:                1,
		// 	Protocol:                "",
		// 	SourceAddressEnd:        "",
		// 	SourceAddressStart:      "",
		// 	SourcePortEnd:           "",
		// 	SourcePortStart:         "",
		// },
		upcloud.FirewallRule{
			Action:               "accept",
			Comment:              "Alow HTTP from anywhere",
			DestinationPortEnd:   "80",
			DestinationPortStart: "80",
			Direction:            "in",
			Family:               "IPv4",
			Position:             1,
			Protocol:             "tcp",
		},
		upcloud.FirewallRule{
			Action:               "accept",
			Comment:              "Allow SSH from a specific network only",
			DestinationPortEnd:   "22",
			DestinationPortStart: "22",
			Direction:            "in",
			Family:               "IPv4",
			ICMPType:             "tcp",
			Position:             2,
			Protocol:             "tcp",
			SourceAddressEnd:     "192.168.1.255",
			SourceAddressStart:   "192.168.1.1",
		},
		// upcloud.FirewallRule{
		// 	Action:                  "accept",
		// 	Comment:                 "Allow SSH over IPv6 from this range",
		// 	DestinationAddressEnd:   "",
		// 	DestinationAddressStart: "",
		// 	DestinationPortEnd:      "22",
		// 	DestinationPortStart:    "22",
		// 	Direction:               "in",
		// 	Family:                  "IPv6",
		// 	ICMPType:                "",
		// 	Position:                3,
		// 	Protocol:                "tcp",
		// 	SourceAddressEnd:        "2a04:3540:1000:aaaa:bbbb:cccc:d001",
		// 	SourceAddressStart:      "2a04:3540:1000:aaaa:bbbb:cccc:d001",
		// 	SourcePortEnd:           "",
		// 	SourcePortStart:         "",
		// },
		upcloud.FirewallRule{
			Action:    "accept",
			Comment:   "Allow ICMP echo request (ping)",
			Direction: "in",
			Family:    "IPv4",
			Position:  4,
			Protocol:  "icmp",
			ICMPType:  "8",
		},
		upcloud.FirewallRule{
			Action:    "drop",
			Direction: "in",
			Position:  5,
		},
	}

	//hardcoded_tag := "kraut-provisioned-" + prov_project

	request := upcloud_request.CreateServerRequest{
		//AvoidHost  string `xml:"avoid_host,omitempty"`
		//BootOrder  string `xml:"boot_order,omitempty"`
		//CoreNumber int    `xml:"core_number,omitempty"`
		Firewall:    "on",
		Hostname:    prov_project,
		IPAddresses: prov_networks,
		LoginUser:   &prov_user,
		// MemoryAmount: 2048,
		PasswordDelivery: "none",
		Plan:             "1xCPU-1GB",
		StorageDevices:   prov_storages,
		TimeZone:         "Europe/Helsinki",
		Title:            prov_project + ": provisioned automatically by kraut",
		UserData:         prov_initscript,
		//VNC: "off",
		Zone: prov_zone,
	}

	serverDetails, err := service.CreateServer(&request)

	if err == nil {

		log.Info("Server created, waiting for it to start, before enabling firewall")

		tagRequest := upcloud_request.TagServerRequest{
			UUID: serverDetails.UUID,
			Tags: []string{"kraut-provisioned"},
		}
		if tagDetails, err := service.TagServer(&tagRequest); err == nil {
			log.WithFields(log.Fields{"UUID": tagDetails.UUID, "hostname": tagDetails.Hostname, "ips": tagDetails.IPAddresses, "state": tagDetails.State, "progress": serverDetails.Progress}).Info("Created and tagged custom server")
		} else {
			result.Set(true, []error{err})
			log.WithError(err).WithFields(log.Fields{"UUID": serverDetails.UUID, "hostname": serverDetails.Hostname, "ips": serverDetails.IPAddresses, "state": serverDetails.State, "progress": serverDetails.Progress}).Warn("Created custom server, but could not tag it")
		}

		waitRequest := upcloud_request.WaitForServerStateRequest{
			UUID:           serverDetails.UUID,
			DesiredState:   "stopped",
			UndesiredState: "started",
			Timeout:        time.Duration(60) * time.Second,
		}
		if _, err := service.WaitForServerState(&waitRequest); err != nil {
			log.WithError(err).Error("Machine startup timed out")
		}

		for _, rule := range prov_firewall_rules {

			ruleRequest := upcloud_request.CreateFirewallRuleRequest{
				ServerUUID:   serverDetails.UUID,
				FirewallRule: rule,
			}
			ruleDetails, err := service.CreateFirewallRule(&ruleRequest)

			if err == nil {
				log.WithFields(log.Fields{"rule": ruleDetails}).Debug("Created firewall rule.")
			} else {
				log.WithError(err).Warn("Could not create firewall rule")
			}

		}

	} else {
		result.Set(false, []error{err, errors.New("Unable to provision new server.")})
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

	service := down.ServiceWrapper()
	settings := down.Settings()

	global := false
	properties := down.Properties()
	if globalProp, found := properties.Get(UPCLOUD_GLOBAL_PROPERTY); found {
		global = globalProp.Get().(bool)
		log.WithFields(log.Fields{"key": UPCLOUD_GLOBAL_PROPERTY, "prop": globalProp, "value": global}).Debug("Operate on servers outside of the project")
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

	service := stop.ServiceWrapper()
	settings := stop.Settings()

	global := false
	properties := stop.Properties()
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
