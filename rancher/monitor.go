package rancher

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
)

/**
 * Monitor Handler using Rancher
 */

// Rancher Monitor Handler
type RancherMonitorHandler struct {
	RancherBaseClientHandler
}

// Initialize and activate the Handler
func (monitor *RancherMonitorHandler) Init() api_operation.Result {
	result := api_operation.New_StandardResult()

	base := New_RancherBaseClientOperation(monitor.ConfigSource())

	ops := api_operation.Operations{}

	ops.Add(api_operation.Operation(&RancherMonitorListEnvironmentsOperation{RancherBaseClientOperation: *base}))

	monitor.operations = &ops

	return api_operation.Result(&result)
}

// Rturn a string identifier for the Handler (not functionally needed yet)
func (monitor *RancherMonitorHandler) Id() string {
	return "rancher.monitor"
}
