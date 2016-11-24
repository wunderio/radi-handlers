package configwrapper

import (
	// "errors"
	// "regexp"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"

	// api_operation "github.com/james-nesbitt/kraut-api/operation"
	// api_config "github.com/james-nesbitt/kraut-api/operation/config"
	api_security "github.com/james-nesbitt/kraut-api/operation/security"
)

// Retrieve values by parsing bytes from the wrapper
func (security *SecurityConfigWrapperYml) LoadUser() error {
	security.userHandler = SecurityConfigWrapperUserYmlHandler{}

	if sources, err := security.wrapper.Get(CONFIG_KEY_SECURITY_USER); err == nil {
		for _, scope := range sources.Order() {
			scopedSource, _ := sources.Get(scope)
			scopedValues := SecurityConfigWrapperUserYmlDefinition{}
			if err := yaml.Unmarshal(scopedSource, &scopedValues); err == nil {
				security.userHandler.Add(scope, scopedValues)
			} else {
				log.WithError(err).WithFields(log.Fields{"scope": scope}).Error("SecurityConfigWrapperYml Couldn't unmarshall user yml scope")
			}
			//log.WithFields(log.Fields{"values": scopedValues, "userHandler": security.userHandler, "scope": scope}).Info("Security:Config->Load()")
		}
		return nil
	} else {
		log.WithError(err).Error("Error loading config for " + CONFIG_KEY_SECURITY_USER)
		return err
	}
}

/**
 * A handler for users from config yml
 */

type SecurityConfigWrapperUserYmlHandler struct {
	def SecurityConfigWrapperUserYmlDefinition
}

// Lazy initializer
func (userHandler *SecurityConfigWrapperUserYmlHandler) safe() {
	if userHandler.def.UserId == "" {
		userHandler.def.UserId = "anonymous"
	}
}

// Add a definition and merge it
func (userHandler *SecurityConfigWrapperUserYmlHandler) Add(scope string, def SecurityConfigWrapperUserYmlDefinition) {
	userHandler.safe()
	userHandler.def.Merge(def)
}

// Add a definition and merge it
func (userHandler *SecurityConfigWrapperUserYmlHandler) CurrentUser() api_security.SecurityUser {
	userHandler.safe()
	return userHandler.def.SecurityUser()
}

// User definition from yml
type SecurityConfigWrapperUserYmlDefinition struct {
	UserId    string `yaml:"ID"`
	UserLabel string `yaml:Label"`
}

// Convert this into a SecurityUser
func (userDef *SecurityConfigWrapperUserYmlDefinition) SecurityUser() api_security.SecurityUser {
	return api_security.SecurityUser(userDef)
}

// Interface:SecurityUser : return a string machine name id
func (userDef *SecurityConfigWrapperUserYmlDefinition) Id() string {
	return userDef.UserId
}

// interface:SecurityUser : returns a readable string label
func (userDef *SecurityConfigWrapperUserYmlDefinition) Label() string {
	return userDef.UserLabel
}

// merge definitions
func (userDef *SecurityConfigWrapperUserYmlDefinition) Merge(merge SecurityConfigWrapperUserYmlDefinition) {
	if merge.UserId != "" {
		userDef.UserId = merge.UserId
	}
	if merge.UserLabel != "" {
		userDef.UserLabel = merge.UserLabel
	}
}
