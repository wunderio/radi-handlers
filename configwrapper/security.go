package configwrapper

import (
	api_operation "github.com/james-nesbitt/kraut-api/operation"
	api_security "github.com/james-nesbitt/kraut-api/operation/security"
)

/**
 * Security operations that are derived from
 * ConfigWrapper source
 */

const (
	// The Config key for settings
	CONFIG_KEY_SECURITY_AUTHORIZE = "authorize"
	CONFIG_KEY_SECURITY_USER      = "user"
)

// SecurityWrapper definition
type SecurityConfigWrapper interface {
	AuthorizeOperation(api_operation.Operation) api_security.RuleResult
	CurrentUser() api_security.SecurityUser
}

/**
 * A Base operation that holds the SecurityConfigWrapper
 * which can be shared across all security operations
 * that are based on using the Wrapper
 */

// Constructor for SecurityWrapperBaseOperation
func New_SecurityWrapperBaseOperation(wrapper SecurityConfigWrapper) *SecurityWrapperBaseOperation {
	return &SecurityWrapperBaseOperation{
		securityConfigWrapper: wrapper,
	}
}

// SecurityWrapper base operation
type SecurityWrapperBaseOperation struct {
	securityConfigWrapper SecurityConfigWrapper
}

// Get the security wrapper
func (base *SecurityWrapperBaseOperation) SecurityConfigWrapper() SecurityConfigWrapper {
	return base.securityConfigWrapper
}

/**
 * Operations
 */

// ConfigWrapper based security Authorize operation
type SecurityConfigWrapperUserOperation struct {
	SecurityWrapperBaseOperation
	api_security.BaseSecurityUserOperation

	properties *api_operation.Properties
}

// Run a validation check on the Operation
func (userOp *SecurityConfigWrapperUserOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (userOp *SecurityConfigWrapperUserOperation) Properties() *api_operation.Properties {
	if userOp.properties == nil {
		ops := api_operation.Properties{}

		ops.Add(api_operation.Property(&api_security.SecurityUserProperty{}))

		userOp.properties = &ops
	}

	return userOp.properties
}

// Execute the Operation
// @TODO Better error checking is needed in this exec
func (userOp *SecurityConfigWrapperUserOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	securityWrapper := userOp.SecurityConfigWrapper()

	props := userOp.Properties()

	userProp, _ := props.Get(api_security.SECURITY_USER_PROPERTY_KEY)

	currentUser := securityWrapper.CurrentUser()
	userProp.Set(currentUser)

	return api_operation.Result(&result)
}

// ConfigWrapper based security Authorize operation
type SecurityConfigWrapperAuthorizeOperation struct {
	SecurityWrapperBaseOperation
	api_security.BaseSecurityAuthorizeOperation

	targetOperation api_operation.Operation

	properties *api_operation.Properties
}

// Run a validation check on the Operation
func (authorize *SecurityConfigWrapperAuthorizeOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (authorize *SecurityConfigWrapperAuthorizeOperation) Properties() *api_operation.Properties {
	if authorize.properties == nil {
		ops := api_operation.Properties{}

		ops.Add(api_operation.Property(&api_security.SecurityUserProperty{}))
		ops.Add(api_operation.Property(&api_security.SecurityAuthorizationOperationProperty{}))
		ops.Add(api_operation.Property(&api_security.SecurityAuthorizationRuleResultProperty{}))
		ops.Add(api_operation.Property(&api_security.SecurityAuthorizationSucceededProperty{}))

		authorize.properties = &ops
	}

	return authorize.properties
}

// Execute the Operation
func (authorize *SecurityConfigWrapperAuthorizeOperation) Exec() api_operation.Result {
	result := api_operation.BaseResult{}
	result.Set(true, []error{})

	props := authorize.Properties()

	securityWrapper := authorize.SecurityConfigWrapper()

	userProp, _ := props.Get(api_security.SECURITY_USER_PROPERTY_KEY)
	userProp.Set(securityWrapper.CurrentUser())

	/**
	 * Authorize this operation, as opposed to the child, because this operation
	 * also has props for current user etc, which can also be used for the
	 * comparison.  This operation decorates anyway, and pretends to be the child
	 * operation in most ways.
	 */

	ruleResult := securityWrapper.AuthorizeOperation(authorize)

	propRuleResult, _ := props.Get(api_security.SECURITY_AUTHORIZATION_RULERESULT_PROPERTY_KEY)
	propRuleResult.Set(ruleResult)
	propSuccess, _ := props.Get(api_security.SECURITY_AUTHORIZATION_SUCCEEDED_PROPERTY_KEY)
	propSuccess.Set(ruleResult.Allow())

	return api_operation.Result(&result)
}
