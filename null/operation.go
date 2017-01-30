package null

/**
 * Operations that the Null Handler implements
 */

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_command "github.com/wunderkraut/radi-api/operation/command"
	api_config "github.com/wunderkraut/radi-api/operation/config"
	api_document "github.com/wunderkraut/radi-api/operation/document"
	api_monitor "github.com/wunderkraut/radi-api/operation/monitor"
	api_orchestrate "github.com/wunderkraut/radi-api/operation/orchestrate"
	api_security "github.com/wunderkraut/radi-api/operation/security"
	api_setting "github.com/wunderkraut/radi-api/operation/setting"
)

/**
 * Command
 */

// Null operation for listing commands
type NullCommandListOperation struct {
	NullNoPropertiesOperation
	NullAllwaysTrueOperation
	api_command.BaseCommandListOperation
}

// Null operation for executing a command
type NullCommandExecOperation struct {
	NullNoPropertiesOperation
	NullAllwaysTrueOperation
	api_command.BaseCommandExecOperation
}

/**
 * Config
 */

// Null Configuration retreive readers operation
type NullConfigReadersOperation struct {
	NullNoPropertiesOperation
	NullAllwaysTrueOperation
	api_config.BaseConfigReadersOperation
}

// Null Configuration retrieve writers operation
type NullConfigWritersOperation struct {
	NullNoPropertiesOperation
	NullAllwaysTrueOperation
	api_config.BaseConfigWritersOperation
}

/**
 * Setting
 */

// Null Setting retreive accessor operation
type NullSettingGetOperation struct {
	NullNoPropertiesOperation
	NullAllwaysTrueOperation
	api_setting.BaseSettingGetOperation
}

// Null Setting assign accessor operation
type NullSettingSetOperation struct {
	NullNoPropertiesOperation
	NullAllwaysTrueOperation
	api_setting.BaseSettingSetOperation
}

/**
 * Documentationm
 */

// Null operation for listing documentation topics
type NullDocumentTopicListOperation struct {
	NullNoPropertiesOperation
	NullAllwaysTrueOperation
	api_document.BaseDocumentTopicListOperation
}

// Null Operation for retrieving a single documentation topic
type NullDocumentTopicGetOperation struct {
	NullNoPropertiesOperation
	NullAllwaysTrueOperation
	api_document.BaseDocumentTopicGetOperation
}

/**
 * Monitor
 */

// Null operation for monitoring information
type NullMonitorInfoOperation struct {
	NullNoPropertiesOperation
	NullAllwaysTrueOperation
	api_monitor.BaseMonitorInfoOperation
}

// Null status operation exec method
func (info *NullMonitorInfoOperation) Exec(props *api_operation.Properties) api_operation.Result {
	message := "App is using NULL Info handler\n"
	info.WriteMessage(message)

	return info.NullAllwaysTrueOperation.Exec()
}

// Null operation for monitoring status
type NullMonitorStatusOperation struct {
	NullNoPropertiesOperation
	NullAllwaysTrueOperation
	api_monitor.BaseMonitorStatusOperation
}

// Null status operation exec method
func (status *NullMonitorStatusOperation) Exec(props *api_operation.Properties) api_operation.Result {
	message := "App is using NULL status handler\n"
	status.WriteMessage(message)

	return status.NullAllwaysTrueOperation.Exec()
}

/**
 * Orchestration
 */

// Null operation for orchestration UP
type NullOrchestrateUpOperation struct {
	NullNoPropertiesOperation
	NullAllwaysTrueOperation
	api_orchestrate.BaseOrchestrationUpOperation
}

// Null operation for orchestration DOWN
type NullOrchestrateDownOperation struct {
	NullNoPropertiesOperation
	NullAllwaysTrueOperation
	api_orchestrate.BaseOrchestrationDownOperation
}

/**
 * Security
 */

// Null Authenticate always authenticates
type NullSecurityAuthenticateOperation struct {
	NullNoPropertiesOperation
	NullAllwaysTrueOperation
	api_security.BaseSecurityAuthenticateOperation
}

// Null Authorize always authorizes
type NullSecurityAuthorizeOperation struct {
	NullNoPropertiesOperation
	NullAllwaysTrueOperation
	api_security.BaseSecurityAuthorizeOperation
}

// Null User, provides a consistent user value
type NullSecurityUserOperation struct {
	NullNoPropertiesOperation
	NullAllwaysTrueOperation
	api_security.BaseSecurityUserOperation
}
