package null

import (
	api_api "github.com/james-nesbitt/kraut-api/api"
	api_builder "github.com/james-nesbitt/kraut-api/builder"
	api_operation "github.com/james-nesbitt/kraut-api/operation"
	api_monitor "github.com/james-nesbitt/kraut-api/operation/monitor"
)

/**
 * Pretty simple builder, which will just add operations for 
 * groups that are activated
 */

// NullBuilder Constructor
func New_NullBuilder() *NullBuilder {
	return &NullBuilder{}
}

// API BUilder that provides many null operations
type NullBuilder struct {
	activated []string
}

// Set a API for this Handler
func (builder *NullBuilder) SetAPI(parent api_api.API) {
	// do nothing, who cares
}

// Initialize and activate the Handler
func (builder *NullBuilder) Activate(implementations api_builder.Implementations, settingsProvider api_builder.SettingsProvider) error {
	for _, implementation := range implementations.Order() {
		found := false
		for _, existing := range builder.activated {
			if existing == implementation {
				found = true
				break;
			}
		}
		if !found {
			builder.activated = append(builder.activated, implementation)
		}
	}

	return nil
}

// Rturn a string identifier for the Handler (not functionally needed yet)
func (builder *NullBuilder) Id() string { 
	return "null"
}

// Return a list of Operations from the Handler
func (builder *NullBuilder) Operations() *api_operation.Operations {
	operations := api_operation.Operations{}

	for _, activated := range builder.activated {
		switch (activated) {
		case "config":
			// Add Null config operations
			operations.Add(api_operation.Operation(&NullConfigReadersOperation{}))
			operations.Add(api_operation.Operation(&NullConfigWritersOperation{}))
		case "setting":
			// Add Null setting operations
			operations.Add(api_operation.Operation(&NullSettingGetOperation{}))
			operations.Add(api_operation.Operation(&NullSettingSetOperation{}))
		case "command":
			// Add Null command operations
			operations.Add(api_operation.Operation(&NullCommandListOperation{}))
			operations.Add(api_operation.Operation(&NullCommandExecOperation{}))
		case "document":
			// Add Null documentation operations
			operations.Add(api_operation.Operation(&NullDocumentTopicListOperation{}))
			operations.Add(api_operation.Operation(&NullDocumentTopicGetOperation{}))
		case "monitor": 
			// Add null monitor operations
			operations.Add(api_operation.Operation(&NullMonitorStatusOperation{}))
			operations.Add(api_operation.Operation(&NullMonitorInfoOperation{}))
			operations.Add(api_operation.Operation(&api_monitor.MonitorStandardLogOperation{}))
		case "orchestrate":
			// Add Null orchestration operations
			operations.Add(api_operation.Operation(&NullOrchestrateUpOperation{}))
			operations.Add(api_operation.Operation(&NullOrchestrateDownOperation{}))
		case "security":
			// Add Null security handlers
			operations.Add(api_operation.Operation(&NullSecurityAuthenticateOperation{}))
			operations.Add(api_operation.Operation(&NullSecurityAuthorizeOperation{}))
			operations.Add(api_operation.Operation(&NullSecurityUserOperation{}))
		}
	}

	return &operations
}
