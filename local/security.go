package local

import (
	log "github.com/Sirupsen/logrus"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_config "github.com/wunderkraut/radi-api/operation/config"
	api_security "github.com/wunderkraut/radi-api/operation/security"
	api_property "github.com/wunderkraut/radi-api/property"
	api_result "github.com/wunderkraut/radi-api/result"
	api_usage "github.com/wunderkraut/radi-api/usage"
	handlers_configwrapper "github.com/wunderkraut/radi-handlers/configwrapper"
)

/**
 * Local implementation of security operations
 *
 * This security implementation bases authentication
 * on local user assessment, and authorization on a
 * configwrapper source, based on a local file.
 *
 * The local file is not considered safe, but rather
 * is considered an initial implementation, which
 * can also be considered as a project based import
 * source to populate external authorization tools.
 */

// A handler for local security
type LocalHandler_Security struct {
	LocalHandler_Base
	LocalHandler_ConfigWrapperBase
}

// Identify the handler
func (handler *LocalHandler_Security) Id() string {
	return "local.security"
}

// [Handler.]Init tells the LocalHandler_Orchestrate to prepare it's operations
func (handler *LocalHandler_Security) Operations() api_operation.Operations {
	ops := api_operation.New_SimpleOperations()

	// Make a SecurityWrapper Base operation
	securityWrapper := handlers_configwrapper.New_SecurityConfigWrapperYml(handler.ConfigWrapper()).SecurityConfigWrapper()
	base := handlers_configwrapper.New_SecurityWrapperBaseOperation(securityWrapper)

	// Add operations from using the base
	ops.Add(api_operation.Operation(New_LocalCurrentUserOperation(handler.LocalHandler_Base.settings, base)))
	ops.Add(api_operation.Operation(&handlers_configwrapper.SecurityConfigWrapperAuthorizeOperation{SecurityWrapperBaseOperation: *base}))

	return ops.Operations()
}

// Make ConfigWrapper
func (handler *LocalHandler_Security) SecurityWrapper() api_security.SecurityWrapper {
	return api_security.New_SimpleSecurityWrapper(handler.Operations()).SecurityWrapper()
}

/**
 * A local config based CurrentUser operation
 *
 * We use this as a wrapper around the SecurityConfigWrapperUserOperation
 * in order to provide a fallback case for when no user configwrapper source
 * is available.
 */

// Local Current user
type LocalCurrentUserOperation struct {
	handlers_configwrapper.SecurityConfigWrapperUserOperation
	settings      *LocalAPISettings
	configWrapper api_config.ConfigWrapper
}

func New_LocalCurrentUserOperation(settings *LocalAPISettings, base *handlers_configwrapper.SecurityWrapperBaseOperation) *LocalCurrentUserOperation {
	configWrapperUserOperation := handlers_configwrapper.SecurityConfigWrapperUserOperation{
		SecurityWrapperBaseOperation: *base,
	}
	return &LocalCurrentUserOperation{
		SecurityConfigWrapperUserOperation: configWrapperUserOperation,
		settings: settings,
	}
}

func (userOp *LocalCurrentUserOperation) Exec(props api_property.Properties) api_result.Result {
	result := api_result.New_StandardResult()
	securityWrapper := userOp.SecurityConfigWrapper()

	userProp, _ := props.Get(api_security.SECURITY_USER_PROPERTY_KEY)

	currentUser := securityWrapper.CurrentUser()

	if currentUser == nil || currentUser.Id() == "anonymous" {
		settings := userOp.settings
		localUser := &settings.User
		if localUser != nil {
			currentUser = api_security.New_CoreUserSecurityUser(localUser).SecurityUser()
			log.WithFields(log.Fields{"id": currentUser.Id(), "label": currentUser.Label()}).Debug("Retrieved current user from OS user default")
		} else {
			log.WithFields(log.Fields{"id": currentUser.Id(), "label": currentUser.Label()}).Debug("Using anonymous user, as no better user was available")
		}
	} else {
		log.WithFields(log.Fields{"id": currentUser.Id(), "label": currentUser.Label()}).Debug("Retrieved current user from config")
	}

	userProp.Set(currentUser)
	result.MarkSuccess()

	result.MarkFinished()

	return api_result.Result(result)
}

func (userOp *LocalCurrentUserOperation) Usage() api_usage.Usage {
	return api_operation.Usage_External()
}
