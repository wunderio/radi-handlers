package null

import (
	api_api "github.com/wunderkraut/radi-api/api"
	api_builder "github.com/wunderkraut/radi-api/builder"
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_monitor "github.com/wunderkraut/radi-api/operation/monitor"
	api_result "github.com/wunderkraut/radi-api/result"
)

/**
 * Pretty simple builder, which will just add operations for
 * groups that are activated
 */

// API BUilder that provides many null operations
type NullBuilder struct {
	activated []string
}

// NullBuilder Constructor
func New_NullBuilder() *NullBuilder {
	return &NullBuilder{}
}

// Rturn a string identifier for the Handler (not functionally needed yet)
func (builder *NullBuilder) Id() string {
	return "null"
}

// Set a API for this Handler
func (builder *NullBuilder) SetAPI(parent api_api.API) {
	// do nothing, who cares
}

// Initialize and activate the Handler
func (builder *NullBuilder) Activate(implementations api_builder.Implementations, settingsProvider api_builder.SettingsProvider) api_result.Result {
	for _, implementation := range implementations.Order() {
		found := false
		for _, existing := range builder.activated {
			if existing == implementation {
				found = true
				break
			}
		}
		if !found {
			builder.activated = append(builder.activated, implementation)
		}
	}

	return api_result.MakeSuccessfulResult()
}

// Validate the builder after Activation is complete
func (builder *NullBuilder) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Return a list of Operations from the Handler
func (builder *NullBuilder) Operations() api_operation.Operations {
	ops := api_operation.New_SimpleOperations()

	for _, activated := range builder.activated {
		switch activated {
		case "config":
			// Add Null config operations
			ops.Add(api_operation.Operation(&NullConfigReadersOperation{}))
			ops.Add(api_operation.Operation(&NullConfigWritersOperation{}))
		case "setting":
			// Add Null setting operations
			ops.Add(api_operation.Operation(&NullSettingGetOperation{}))
			ops.Add(api_operation.Operation(&NullSettingSetOperation{}))
		case "command":
			// Add Null command operations
			ops.Add(api_operation.Operation(&NullCommandListOperation{}))
			ops.Add(api_operation.Operation(&NullCommandExecOperation{}))
		case "document":
			// Add Null documentation operations
			ops.Add(api_operation.Operation(&NullDocumentTopicListOperation{}))
			ops.Add(api_operation.Operation(&NullDocumentTopicGetOperation{}))
		case "monitor":
			// Add null monitor operations
			ops.Add(api_operation.Operation(&NullMonitorStatusOperation{}))
			ops.Add(api_operation.Operation(&NullMonitorInfoOperation{}))
			ops.Add(api_operation.Operation(&api_monitor.MonitorStandardLogOperation{}))
		case "orchestrate":
			// Add Null orchestration operations
			ops.Add(api_operation.Operation(&NullOrchestrateUpOperation{}))
			ops.Add(api_operation.Operation(&NullOrchestrateDownOperation{}))
		case "security":
			// Add Null security handlers
			ops.Add(api_operation.Operation(&NullSecurityAuthenticateOperation{}))
			ops.Add(api_operation.Operation(&NullSecurityAuthorizeOperation{}))
			ops.Add(api_operation.Operation(&NullSecurityUserOperation{}))
		}
	}

	return ops.Operations()
}
