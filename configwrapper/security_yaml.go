package configwrapper

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	// "gopkg.in/yaml.v2"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_config "github.com/wunderkraut/radi-api/operation/config"
	api_security "github.com/wunderkraut/radi-api/operation/security"
)

// Constructor for SecurityConfigWrapperYml
func New_SecurityConfigWrapperYml(wrapper api_config.ConfigWrapper) *SecurityConfigWrapperYml {
	return &SecurityConfigWrapperYml{
		wrapper: wrapper,
	}
}

// A SecurityConfirWrapper that reads config as yml
type SecurityConfigWrapperYml struct {
	authHandler SecurityConfigWrapperAuthorizeYmlHandler
	userHandler SecurityConfigWrapperUserYmlHandler
	wrapper     api_config.ConfigWrapper
}

// Convert this into a SecurityConfigWrapper
func (security *SecurityConfigWrapperYml) SecurityConfigWrapper() SecurityConfigWrapper {
	security.safe()
	return SecurityConfigWrapper(security)
}

// Safe lazy initializer
func (security *SecurityConfigWrapperYml) safe() {
	if security.authHandler.Empty() {
		security.LoadAuthorize() // @see security_yaml_authorization.go
		security.LoadUser()      // @see security_yaml_user.go
	}
}

// Save the current values to the wrapper
func (security *SecurityConfigWrapperYml) Save() error {
	err := errors.New("SecurityConfigWrapper.Save() not yet writtent")
	log.WithError(err).Error("Could not save security config")
	return err
}

func (security *SecurityConfigWrapperYml) AuthorizeRules() api_security.AuthorizeOperationRules {
	security.safe()
	return security.authHandler.Rules()
}

// Get an ordered list of rules (SecurityConfigWrapper interface)
func (security *SecurityConfigWrapperYml) AuthorizeOperation(op api_operation.Operation) api_security.RuleResult {
	//log.WithFields(log.Fields{"op": op.Id()}).Info("Authorizing operation")
	return security.authHandler.Rules().AuthorizeOperation(op)
}

// Get an ordered list of rules (SecurityConfigWrapper interface)
func (security *SecurityConfigWrapperYml) CurrentUser() api_security.SecurityUser {
	return security.userHandler.CurrentUser()
}

// Return the default scope string for the wrapper
func (security *SecurityConfigWrapperYml) DefaultScope() string {
	/**
	 * @TODO come up with better scopes, but it has to match local conf path keys
	 */
	return "project"
}
