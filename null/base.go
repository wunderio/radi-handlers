package null

import (
	api_operation "github.com/james-nesbitt/kraut-api/operation"
)

/**
 * Some Base operations for the null operations to extend
 */

// A Null base operation that presents empty properties
type NullNoPropertiesOperation struct{}

// Return operation properties
func (null *NullNoPropertiesOperation) Properties() *api_operation.Properties {
	return &api_operation.Properties{}
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
func (alwaystrue *NullAllwaysTrueOperation) Exec() api_operation.Result {
	baseResult := api_operation.BaseResult{}
	baseResult.Set(true, []error{})
	return api_operation.Result(&baseResult)
}
