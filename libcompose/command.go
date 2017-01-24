package libcompose

import (
	"errors"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_command "github.com/wunderkraut/radi-api/operation/command"
)

/**
 * Implement command containers that mix into
 * libCompose orchestrated containers
 */

const (
	CONFIG_KEY_COMMAND = "commands" // The Config key for settings
)

// A wrapper interface which pulls command information from a config wrapper backend
type CommandConfigWrapper interface {
	List(parent string) ([]string, error)
	Get(key string) (*CommandYmlCommand, error)
}

/**
 * Operations
 */

// LibCompose Command List operation
type LibcomposeCommandListOperation struct {
	api_command.BaseCommandListOperation
	api_command.BaseCommandKeyKeysOperation
	BaseLibcomposeNameFilesOperation

	Wrapper    CommandConfigWrapper
	properties *api_operation.Properties
}

// Validate the operation
func (list *LibcomposeCommandListOperation) Validate() bool {
	return true
}

// Get properties
func (list *LibcomposeCommandListOperation) Properties() *api_operation.Properties {
	baseProps := api_operation.Properties{}

	keyKeysProps := list.BaseCommandKeyKeysOperation.Properties()
	baseProps.Merge(*keyKeysProps)
	libComposeBaseProps := list.BaseLibcomposeNameFilesOperation.Properties()
	baseProps.Merge(*libComposeBaseProps)

	return &baseProps
}

// Execute the libCompose Command List operation
func (list *LibcomposeCommandListOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, nil)

	props := list.BaseCommandKeyKeysOperation.Properties()
	keyProp, _ := props.Get(api_command.OPERATION_PROPERTY_COMMAND_KEY)
	keysProp, _ := props.Get(api_command.OPERATION_PROPERTY_COMMAND_KEYS)

	parent := ""
	if key, ok := keyProp.Get().(string); ok && key != "" {
		parent = key
	}

	if keyList, err := list.Wrapper.List(parent); err == nil {
		keysProp.Set(keyList)
	} else {
		result.Set(false, []error{err})
	}

	return api_operation.Result(&result)
}

// LibCompose Command Get operation
type LibcomposeCommandGetOperation struct {
	api_command.BaseCommandGetOperation
	api_command.BaseCommandKeyCommandOperation
	BaseLibcomposeNameFilesOperation

	Wrapper    CommandConfigWrapper
	properties *api_operation.Properties
}

// Validate the operation
func (get *LibcomposeCommandGetOperation) Validate() bool {
	return true
}

// Get properties
func (get *LibcomposeCommandGetOperation) Properties() *api_operation.Properties {
	baseProps := api_operation.Properties{}

	keyCommandProps := get.BaseCommandKeyCommandOperation.Properties()
	baseProps.Merge(*keyCommandProps)
	libComposeBaseProps := get.BaseLibcomposeNameFilesOperation.Properties()
	baseProps.Merge(*libComposeBaseProps)

	return &baseProps
}

// Execute the libCompose Command Get operation
func (get *LibcomposeCommandGetOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, nil)

	props := get.BaseCommandKeyCommandOperation.Properties()
	keyProp, _ := props.Get(api_command.OPERATION_PROPERTY_COMMAND_KEY)
	commandProp, _ := props.Get(api_command.OPERATION_PROPERTY_COMMAND_COMMAND)

	if key, ok := keyProp.Get().(string); ok && key != "" {

		if comYml, err := get.Wrapper.Get(key); err == nil {
			// pass all props to make a project
			com := comYml.Command(get.BaseLibcomposeNameFilesOperation.Properties())
			commandProp.Set(com)
		} else {
			result.Set(false, []error{err})
		}

	} else {
		result.Set(false, []error{errors.New("No command name provided.")})
	}

	return api_operation.Result(&result)
}
