package null

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_property "github.com/wunderkraut/radi-api/property"
	api_result "github.com/wunderkraut/radi-api/result"
	api_usage "github.com/wunderkraut/radi-api/usage"
)

/**
 * Some Base operations for the null operations to extend
 */

// A Null base operation that presents empty properties
type NullNoPropertiesOperation struct{}

// Return operation properties
func (null *NullNoPropertiesOperation) Properties() api_property.Properties {
	return api_property.New_SimplePropertiesEmpty().Properties()
}

// Null base operation which always execs TRUE
type NullInternalUsageOperation struct{}

// Validate the operation
func (internal *NullInternalUsageOperation) Usage() api_usage.Usage {
	return api_operation.Usage_Internal()
}

// Null base operation which always execs TRUE
type NullAllwaysTrueOperation struct{}

// Exec the operation
func (alwaystrue *NullAllwaysTrueOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Exec the operation
func (alwaystrue *NullAllwaysTrueOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	res.MarkSuccess()
	res.MarkFinished()

	return res.Result()
}
