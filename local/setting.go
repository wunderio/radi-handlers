package local

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_setting "github.com/wunderkraut/radi-api/operation/setting"

	handler_configwrapper "github.com/wunderkraut/radi-handlers/configwrapper"
)

// A handler for local settings
type LocalHandler_Setting struct {
	LocalHandler_Base
	LocalHandler_ConfigWrapperBase
}

// Identify the handler
func (handler *LocalHandler_Setting) Id() string {
	return "local.setting"
}

// [Handler.]Init tells the LocalHandler_Orchestrate to prepare it's operations
func (handler *LocalHandler_Setting) Operations() api_operation.Operations {
	ops := api_operation.New_SimpleOperations()

	// Make a wrapper for the Settings Config interpretation, based on itnerpreting YML settings
	wrapper := handler_configwrapper.SettingsConfigWrapper(handler_configwrapper.New_BaseSettingConfigWrapperYmlOperation(handler.ConfigWrapper()))

	// Now we can add config operations that use that Base class
	ops.Add(api_operation.Operation(&handler_configwrapper.SettingConfigWrapperGetOperation{Wrapper: wrapper}))
	ops.Add(api_operation.Operation(&handler_configwrapper.SettingConfigWrapperSetOperation{Wrapper: wrapper}))
	ops.Add(api_operation.Operation(&handler_configwrapper.SettingConfigWrapperListOperation{Wrapper: wrapper}))

	return ops.Operations()
}

// Make ConfigWrapper
func (handler *LocalHandler_Setting) SettingWrapper() api_setting.SettingWrapper {
	return api_setting.New_SimpleSettingWrapper(handler.Operations())
}
