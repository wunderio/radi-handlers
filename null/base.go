package null

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
)

/**
 * Some Base operations for the null operations to extend
 */

// A Null base operation that presents empty properties
type NullNoPropertiesOperation struct{}

// Return operation properties
func (null *NullNoPropertiesOperation) Properties() api_operation.Properties {
	return api_operation.Properties{}
}

// Null base operation which always execs TRUE
type NullAllwaysTrueOperation struct{}

// Validate the operation
func (alwaystrue *NullAllwaysTrueOperation) Validate() bool {
	return true
}

// return empty Configuraitons
// func (alwaystrue *NullAllwaysTrueOperation) Configurations() *operation.Configurations {
// 	return &operation.Configurations{}
// }
// Exec the operation
func (alwaystrue *NullAllwaysTrueOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()
	result.MarkSuccess()
	result.MarkFinished()

	return api_operation.Result(&baseResult)
}
