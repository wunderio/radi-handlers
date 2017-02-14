package null

/**
 * The NullHandler provides a handlers with a set of operations that are
 * entirly Null provided.
 */

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_result "github.com/wunderkraut/radi-api/result"

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
func (handler *NullHandler) Init() api_result.Result {
	return api_result.New_StandardResult().Result()
}

// [Handler.]Operations returns an Operations list of a number of different Null operations
func (handler *NullHandler) Operations() api_operation.Operations {
	ops := api_operation.New_SimpleOperations()

	// Add Null config operations
	ops.Add(api_operation.Operation(&NullConfigReadersOperation{}))
	ops.Add(api_operation.Operation(&NullConfigWritersOperation{}))
	// Add Null setting operations
	ops.Add(api_operation.Operation(&NullSettingGetOperation{}))
	ops.Add(api_operation.Operation(&NullSettingSetOperation{}))
	// Add Null command operations
	ops.Add(api_operation.Operation(&NullCommandListOperation{}))
	ops.Add(api_operation.Operation(&NullCommandExecOperation{}))
	// Add Null documentation operations
	ops.Add(api_operation.Operation(&NullDocumentTopicListOperation{}))
	ops.Add(api_operation.Operation(&NullDocumentTopicGetOperation{}))
	// Add null monitor operations
	ops.Add(api_operation.Operation(&NullMonitorStatusOperation{}))
	ops.Add(api_operation.Operation(&NullMonitorInfoOperation{}))
	ops.Add(api_operation.Operation(&api_monitor.MonitorStandardLogOperation{}))
	// Add Null orchestration operations
	ops.Add(api_operation.Operation(&NullOrchestrateUpOperation{}))
	ops.Add(api_operation.Operation(&NullOrchestrateDownOperation{}))
	// Add Null security handlers
	ops.Add(api_operation.Operation(&NullSecurityAuthenticateOperation{}))
	ops.Add(api_operation.Operation(&NullSecurityAuthorizeOperation{}))
	ops.Add(api_operation.Operation(&NullSecurityUserOperation{}))

	return ops.Operations()
}
