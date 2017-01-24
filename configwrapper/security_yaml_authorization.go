package configwrapper

import (
	// "errors"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"

	api_operation "github.com/wunderkraut/radi-api/operation"
	// api_config "github.com/wunderkraut/radi-api/operation/config"
	api_security "github.com/wunderkraut/radi-api/operation/security"
)

/**
 * Add the function for loading auth to the security config wrapper
 */

// Retrieve values by parsing bytes from the wrapper
func (security *SecurityConfigWrapperYml) LoadAuthorize() error {
	security.authHandler = SecurityConfigWrapperAuthorizeYmlHandler{}

	if sources, err := security.wrapper.Get(CONFIG_KEY_SECURITY_AUTHORIZE); err == nil {
		for _, scope := range sources.Order() {
			scopedSource, _ := sources.Get(scope)
			scopedValues := SecurityConfigWrapperAuthorizeYmlDefinition{}
			if err := yaml.Unmarshal(scopedSource, &scopedValues); err == nil {
				security.authHandler.Add(scope, scopedValues)
			} else {
				log.WithError(err).WithFields(log.Fields{"scope": scope}).Error("SecurityConfigWrapper Couldn't unmarshall auth yml scope")
			}
			//log.WithFields(log.Fields{"values": scopedValues, "authHandler": security.authHandler, "scope": scope}).Info("Security:Config->Load()")
		}
		return nil
	} else {
		log.WithError(err).Error("Error loading config for " + CONFIG_KEY_SECURITY_AUTHORIZE)
		return err
	}
}

/**
 * YML holder structs
 */

// Yml Rule set container
type SecurityConfigWrapperAuthorizeYmlHandler struct {
	definitions map[string]*SecurityConfigWrapperAuthorizeYmlDefinition
	order       []string
}

// Is this bject uninitialized
func (handler *SecurityConfigWrapperAuthorizeYmlHandler) Empty() bool {
	return len(handler.order) == 0
}

// safe lazy initializer
func (handler *SecurityConfigWrapperAuthorizeYmlHandler) safe() {
	// @TODO we need a better safe definition than this if statement
	if &handler.definitions == nil || len(handler.definitions) == 0 {
		handler.definitions = map[string]*SecurityConfigWrapperAuthorizeYmlDefinition{}
		handler.order = []string{}
	}
}

// Add a definition
func (handler *SecurityConfigWrapperAuthorizeYmlHandler) Add(id string, definition SecurityConfigWrapperAuthorizeYmlDefinition) {
	handler.safe()
	if _, found := handler.definitions[id]; !found {
		handler.order = append(handler.order, id)
	}
	handler.definitions[id] = &definition
}

// Get a definition
func (handler *SecurityConfigWrapperAuthorizeYmlHandler) Get(id string) (*SecurityConfigWrapperAuthorizeYmlDefinition, bool) {
	handler.safe()
	definition, found := handler.definitions[id]
	return definition, found
}

// List the key order
func (handler *SecurityConfigWrapperAuthorizeYmlHandler) Order() []string {
	handler.safe()
	return handler.order
}

// Get an ordered list of rules (SecurityConfigWrapper interface)
func (handler *SecurityConfigWrapperAuthorizeYmlHandler) Rules() api_security.AuthorizeOperationRules {
	handler.safe()

	rules := api_security.SimpleAuthorizeOperationRules{}

	for _, id := range handler.order {
		definition, _ := handler.Get(id)

		for _, ymlRule := range definition.Rules() {
			if rule, err := ymlRule.Rule(); err == nil {
				//log.WithFields(log.Fields{"rule": ymlRule.Id, "rule-auth": ymlRule.Authorize, "*rule": &ymlRule}).Info("Collected ymlRule")
				rules.Set(ymlRule.Id, rule)
			} else {
				log.WithError(err).WithFields(log.Fields{"rule": ymlRule.Id, "rule-auth": ymlRule.Authorize, "*rule": &rule}).Info("ymlRule could not create a valid rule.")
			}

		}
	}

	return api_security.AuthorizeOperationRules(&rules)
}

// Yml Rule set container
type SecurityConfigWrapperAuthorizeYmlDefinition struct {
	Settings    SecurityConfigWrapperAuthorizeYmlSettings `yaml:"Settings"`
	SourceRules []*SecurityConfigWrapperAuthorizeYmlRule  `yaml:"Rules"`
}

// Get an ordered list of rules
func (definition *SecurityConfigWrapperAuthorizeYmlDefinition) Rules() []*SecurityConfigWrapperAuthorizeYmlRule {
	definition.applyDefaults()
	return definition.SourceRules
}

// Apply defaults to all the rules
func (definition *SecurityConfigWrapperAuthorizeYmlDefinition) applyDefaults() {
	for index, rule := range definition.SourceRules {
		if rule.Message == "" {
			rule.Message = definition.Settings.DefaultMessage
		}
		if rule.Authorize == "" {
			rule.Authorize = definition.Settings.DefaultAuthorize
		}
		if rule.Aggregate == "" {
			rule.Aggregate = definition.Settings.DefaultAggregate
		}

		definition.SourceRules[index] = rule
	}
}

// Yml Rule container
type SecurityConfigWrapperAuthorizeYmlSettings struct {
	DefaultAuthorize string `yaml:"Authorize"`
	DefaultAggregate string `yaml:"Aggregate"`
	DefaultMessage   string `yaml:"Message"`
}

// Yml Rule container
type SecurityConfigWrapperAuthorizeYmlRule struct {
	Id         string              `yaml:"Id"`
	Message    string              `yaml:"Message"`
	Operation  string              `yaml:"Operation"`
	Authorize  string              `yaml:"Authorize"`
	Aggregate  string              `yaml:"Aggregate"`
	Properties map[string][]string `yaml:"Property"`
}

// Conver this YmlRule to an api_security Rule
func (ymlRule *SecurityConfigWrapperAuthorizeYmlRule) Rule() (api_security.AuthorizeOperationRule, error) {
	return api_security.AuthorizeOperationRule(ymlRule), nil
}

// Conver this YmlRule to an api_security Rule
func (ymlRule *SecurityConfigWrapperAuthorizeYmlRule) AuthorizeOperation(op api_operation.Operation) api_security.RuleResult {
	//log.WithFields(log.Fields{"rule": ymlRule, "op": op.Id()}).Info("Checking Rule")

	if len(ymlRule.Operation) == 0 {
		log.WithFields(log.Fields{"rule": ymlRule}).Error("Rule had no operation match definition")
		return ymlRule.result(0)
	}

	// match operation id
	if ymlRule.Operation != "*" {
		if match, _ := regexp.MatchString(ymlRule.Operation, op.Id()); !match {
			//log.WithFields(log.Fields{"match": ymlRule.Operation, "op": op.Id(), "rule": ymlRule}).Info("Rule id did not match")
			return ymlRule.result(0)
		}
	}

	// match properties
	opProps := op.Properties()
	for propId, propValues := range ymlRule.Properties {
		if prop, found := opProps.Get(propId); found {
			if authMatchProperty(prop, propValues) {
				return ymlRule.result(authStringToInt(ymlRule.Authorize, false))
			} else {
				//log.WithFields(log.Fields{"rule": ymlRule, "prop": prop}).Info("failed to match rule")
				return ymlRule.result(authStringToInt(ymlRule.Authorize, true))
			}
		} else if ymlRule.Aggregate == "AND" {
			//log.WithFields(log.Fields{"rule": ymlRule, "prop": propId}).Info("failed to match rule as property was not found in the rule")
			return ymlRule.result(authStringToInt(ymlRule.Authorize, true))
		}

	}

	log.WithFields(log.Fields{"auth": ymlRule.Authorize, "rule": ymlRule.Id, "op": op.Id()}).Debug("SecurityConfigWrapperAuthorizeYmlRule.AuthorizeOperation() : Rule Matched. Applying rule auth")
	return ymlRule.result(authStringToInt(ymlRule.Authorize, false))
}

// Convert this rule into a RuleResult depending on value
func (ymlRule *SecurityConfigWrapperAuthorizeYmlRule) result(authValue int) api_security.RuleResult {
	return api_security.New_SimpleRuleResult(ymlRule.Id, ymlRule.Message, authValue).RuleResult()
}
