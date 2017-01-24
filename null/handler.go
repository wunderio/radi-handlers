package null

/**
 * The NullHandler provides a handlers with a set of operations that are
 * entirly Null provided.
 */

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_monitor "github.com/wunderkraut/radi-api/operation/monitor"
)

// NullHandler Constructor, doesn't do much preprocessing really
func NewNullHandler() *NullHandler {
	nullHandler := NullHandler{}
	return &nullHandler
}

// NullHandler is a handler implementation that provides many core operations, but does very little (but is safe to use)
type NullHandler struct{}

// [Handler.]Id returns a string ID for the handler
func (handler *NullHandler) Id() string {
	return "null"
}

// [Handler.]Init tells the NullHandler to process itself. Return true as Null Handler always validates true
func (handler *NullHandler) Init() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, nil)
	return api_operation.Result(&result)
}

// [Handler.]Operations returns an Operations list of a number of different Null operations
func (handler *NullHandler) Operations() *api_operation.Operations {
	operations := api_operation.Operations{}

	// Add Null config operations
	operations.Add(api_operation.Operation(&NullConfigReadersOperation{}))
	operations.Add(api_operation.Operation(&NullConfigWritersOperation{}))
	// Add Null setting operations
	operations.Add(api_operation.Operation(&NullSettingGetOperation{}))
	operations.Add(api_operation.Operation(&NullSettingSetOperation{}))
	// Add Null command operations
	operations.Add(api_operation.Operation(&NullCommandListOperation{}))
	operations.Add(api_operation.Operation(&NullCommandExecOperation{}))
	// Add Null documentation operations
	operations.Add(api_operation.Operation(&NullDocumentTopicListOperation{}))
	operations.Add(api_operation.Operation(&NullDocumentTopicGetOperation{}))
	// Add null monitor operations
	operations.Add(api_operation.Operation(&NullMonitorStatusOperation{}))
	operations.Add(api_operation.Operation(&NullMonitorInfoOperation{}))
	operations.Add(api_operation.Operation(&api_monitor.MonitorStandardLogOperation{}))
	// Add Null orchestration operations
	operations.Add(api_operation.Operation(&NullOrchestrateUpOperation{}))
	operations.Add(api_operation.Operation(&NullOrchestrateDownOperation{}))
	// Add Null security handlers
	operations.Add(api_operation.Operation(&NullSecurityAuthenticateOperation{}))
	operations.Add(api_operation.Operation(&NullSecurityAuthorizeOperation{}))
	operations.Add(api_operation.Operation(&NullSecurityUserOperation{}))

	return &operations
}
