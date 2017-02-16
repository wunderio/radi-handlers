package local

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_config "github.com/wunderkraut/radi-api/operation/config"

	handler_bytesource "github.com/wunderkraut/radi-handlers/bytesource"
)

// A handler for local config
type LocalHandler_Config struct {
	LocalHandler_Base
}

// Identify the handler
func (handler *LocalHandler_Config) Id() string {
	return "local.config"
}

// Prepare and return operations for the handler
func (handler *LocalHandler_Config) Operations() api_operation.Operations {
	ops := api_operation.New_SimpleOperations()

	// build a ConfigConnector for use with the Config operations.
	connector := handler_bytesource.New_ConfigConnectYmlFiles(handler.settings.ConfigPaths)

	// Build this base operation to be shared across all of our config operations
	baseConnectorOperation := api_config.New_BaseConfigConnectorOperation(connector)

	// Now we can add config operations that use that Base class
	ops.Add(api_operation.Operation(&api_config.ConfigSimpleConnectorReadersOperation{BaseConfigConnectorOperation: *baseConnectorOperation}))
	ops.Add(api_operation.Operation(&api_config.ConfigSimpleConnectorWritersOperation{BaseConfigConnectorOperation: *baseConnectorOperation}))
	ops.Add(api_operation.Operation(&api_config.ConfigSimpleConnectorListOperation{BaseConfigConnectorOperation: *baseConnectorOperation}))

	return ops.Operations()
}

// Make ConfigWrapper
func (handler *LocalHandler_Config) ConfigWrapper() api_config.ConfigWrapper {
	return api_config.ConfigWrapper(api_config.New_SimpleConfigWrapper(handler.Operations()))
}
