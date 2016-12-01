package upcloud

import (
	"errors"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"

	upcloud "github.com/Jalle19/upcloud-go-sdk/upcloud"
	upcloud_client "github.com/Jalle19/upcloud-go-sdk/upcloud/client"
	upcloud_request "github.com/Jalle19/upcloud-go-sdk/upcloud/request"
	upcloud_service "github.com/Jalle19/upcloud-go-sdk/upcloud/service"

	api_config "github.com/james-nesbitt/kraut-api/operation/config"
)

/**
 * This file provides a ConfigWrapper based tool for reading
 * and creating an UpCloud client/service pair for a project
 */

/**
 * Interpreting build config as yml
 */

// A BuilderSettingsConfigWrapper, that interprets build config as yml
type UpcloudFactoryConfigWrapperYaml struct {
	configWrapper api_config.ConfigWrapper

	scope string

	User    Yml_UpcloudFactory_User     `yaml:"Access"`
	Servers []Yml_UpcloudFactory_Server `yaml:"Servers"`
}

// Constructor for UpcloudFactoryConfigWrapperYaml
func New_UpcloudFactoryConfigWrapperYaml(configWrapper api_config.ConfigWrapper) *UpcloudFactoryConfigWrapperYaml {
	return &UpcloudFactoryConfigWrapperYaml{
		configWrapper: configWrapper,
	}
}

// Covnert this into an UpcloudFactory interface
func (configFactory *UpcloudFactoryConfigWrapperYaml) UpcloudFactory() UpcloudFactory {
	return UpcloudFactory(configFactory)
}

func (configFactory *UpcloudFactoryConfigWrapperYaml) DefaultScope() string {
	/**
	 * @TODO come up with better scopes, but it has to match local conf path keys
	 * @SEE configconnect/settings_yml.go which has the same issue
	 */
	return "project"
}

func (configFactory *UpcloudFactoryConfigWrapperYaml) safe() {
	if configFactory.User.Empty() {
		if err := configFactory.Load(); err != nil {
			log.WithError(err).Error("Could not load build configuration")
		}
	}
}

// Convert this YML struct into a Client
func (configFactory *UpcloudFactoryConfigWrapperYaml) Client() *upcloud_client.Client {
	return configFactory.User.Client()
}

// Convert this YML struct into a Service
func (configFactory *UpcloudFactoryConfigWrapperYaml) Service() *upcloud_service.Service {
	client := configFactory.Client()
	return New_UpcloudServiceWrapperFactory(*client).Service()
}

// Convert this YML struct into a Service
func (configFactory *UpcloudFactoryConfigWrapperYaml) ServiceWrapper() *UpcloudServiceWrapper {
	client := configFactory.Client()
	return New_UpcloudServiceWrapperFactory(*client).ServiceWrapper()
}

// Retieve a slice of ServerDefinitions
func (configFactory *UpcloudFactoryConfigWrapperYaml) ServerDefinitions() ServerDefinitions {
	defs := ServerDefinitions{}
	for _, ymlServer := range configFactory.Servers {
		defs.Add(ymlServer.ServerDefinition())
	}
	return defs
}

// Retrieve values by parsing bytes from the wrapper
func (configFactory *UpcloudFactoryConfigWrapperYaml) Load() error {
	log.Debug("Loading UpCloud config")

	if sources, err := configFactory.configWrapper.Get(CONFIG_KEY_UPCLOUD); err == nil {
		for _, scope := range sources.Order() {
			scopedSource, _ := sources.Get(scope)

			// empty out this oobject
			configFactory.scope = scope
			configFactory.User = Yml_UpcloudFactory_User{}
			configFactory.Servers = []Yml_UpcloudFactory_Server{}

			if err := yaml.Unmarshal(scopedSource, &configFactory); err == nil {
				log.WithFields(log.Fields{"servers": configFactory.Servers, "scope": configFactory.scope}).Debug("UpCloud settings parsed from config yml")
				break
			} else {
				log.WithError(err).WithFields(log.Fields{"scope": scope}).Error("Couldn't marshall yml scope for upcloud settings for scope.")
			}
		}
		return nil
	} else {
		log.WithError(err).Error("Error loading Upcloud config")
		return err
	}
}

// Save the current values to the wrapper
func (configFactory *UpcloudFactoryConfigWrapperYaml) Save() error {
	/**
	 * @TODO THIS
	 */
	return errors.New("UpcloudFactoryConfigWrapperYaml Set operation not yet written.")
}

/**
 * Various configuration holders
 */

// A holder of user configuration,
type Yml_UpcloudFactory_User struct {
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
}

// Is this struct populated?
func (ymlFactoryUser *Yml_UpcloudFactory_User) Empty() bool {
	return ymlFactoryUser.User == ""
}

// Convert this YML struct into a Client
func (ymlFactoryUser *Yml_UpcloudFactory_User) Client() *upcloud_client.Client {
	return upcloud_client.New(ymlFactoryUser.User, ymlFactoryUser.Password)
}

// A holder for server configuration from yaml
type Yml_UpcloudFactory_Server struct {
	id   string
	zone string
	plan string

	serverDefinition Yml_UpcloudFactory_ServerDefinition
	firewallRules    Yml_UpcloudFactory_ServerFirewall
}

func (server *Yml_UpcloudFactory_Server) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// MetaData unmarshall
	metaHolder := struct {
		Id   string `yaml:"Id"`
		Zone string `yaml:"Zone"`
		Plan string `yaml:"Plan"`
	}{}
	if err := unmarshal(&metaHolder); err != nil {
		return err
	}
	server.id = metaHolder.Id
	server.zone = metaHolder.Zone
	server.plan = metaHolder.Plan
	// log.WithFields(log.Fields{"id": server.id, "zone": server.zone, "holder": metaHolder}).Info("UPCLOUD:FACTORY:YML:ID")

	// Create Server Request unmarshall
	serverHolder := Yml_UpcloudFactory_ServerDefinition{}
	if err := unmarshal(&serverHolder); err != nil {
		log.Error("YML ERROR SERVER")
		return err
	}
	server.serverDefinition = serverHolder
	// log.WithFields(log.Fields{"server": server.serverDefinition, "holder": serverHolder}).Info("UPCLOUD:FACTORY:YML:SERVER")

	// Firewall rules unmarshall
	firewallHolder := struct {
		Firewall Yml_UpcloudFactory_ServerFirewall `yaml:"Firewall"`
	}{}
	if err := unmarshal(&firewallHolder); err != nil {
		log.Error("YML ERROR FIREWALL")
		return err
	}
	server.firewallRules = firewallHolder.Firewall
	// log.WithFields(log.Fields{"rules": server.firewallRules, "holder": firewallHolder}).Info("UPCLOUD:FACTORY:YML:FIREWALL")

	return nil
}

// Convert this to a ServerDefinition interface
func (server *Yml_UpcloudFactory_Server) ServerDefinition() ServerDefinition {
	return ServerDefinition(server)
}

// Build an upcloud CreateServerReequest
func (server *Yml_UpcloudFactory_Server) Id() string {
	return server.id
}

// Build an upcloud CreateServerReequest
func (server *Yml_UpcloudFactory_Server) CreateServerRequest() upcloud_request.CreateServerRequest {
	return server.serverDefinition.CreateServerRequest()
}

// Build upcloud FirewallRules for the server
func (server *Yml_UpcloudFactory_Server) FirewallRules() upcloud.FirewallRules {
	return server.firewallRules.FirewallRules()
}

/**
 * Various configuration holders
 *
 * These are used as providers of the upcloud structs
 * and are needed because of poor YML handling in the
 * upcloud project.  If the upcloud github project
 * improves it's yaml handling then we can probable
 * just remove these.
 *
 */

// A holder for server firewall rules configuration from yaml
type Yml_UpcloudFactory_ServerFirewall struct {
	Rules []Yml_UpcloudFactory_ServerFirewall_Rule `yaml:"Rules"`
}

// Get upcloud FirewallRules
func (firewall *Yml_UpcloudFactory_ServerFirewall) FirewallRules() upcloud.FirewallRules {
	rules := upcloud.FirewallRules{}
	for _, rule := range firewall.Rules {
		rules.FirewallRules = append(rules.FirewallRules, rule.FirewallRule())
	}
	return rules
}

type Yml_UpcloudFactory_ServerFirewall_Rule struct {
	Action                  string `yaml:"Action"`
	Comment                 string `yaml:"Comment,omitempty"`
	DestinationAddressStart string `yaml:"DestinationAddressStart,omitempty"`
	DestinationAddressEnd   string `yaml:"DestinationAddressEnd,omitempty"`
	DestinationPortStart    int    `yaml:"DestinationPortStart,omitempty"`
	DestinationPortEnd      int    `yaml:"DestinationPortEnd,omitempty"`
	Direction               string `yaml:"Direction"`
	Family                  string `yaml:"Family"`
	ICMPType                string `yaml:"ICMPType,omitempty"`
	Position                int    `yaml:"Position"`
	Protocol                string `yaml:"Protocol,omitempty"`
	SourceAddressStart      string `yaml:"SourceAddressStart,omitempty"`
	SourceAddressEnd        string `yaml:"SourceAddressEnd,omitempty"`
	SourcePortStart         int    `yaml:"SourcePortStart,omitempty"`
	SourcePortEnd           int    `yaml:"SourcePortEnd,omitempty"`
}

// Get upcloud FirewallRules
func (rule *Yml_UpcloudFactory_ServerFirewall_Rule) FirewallRule() upcloud.FirewallRule {
	return upcloud.FirewallRule{
		Action:                  rule.Action,
		Comment:                 rule.Comment,
		DestinationAddressStart: rule.DestinationAddressStart,
		DestinationAddressEnd:   strconv.Itoa(rule.DestinationPortStart),
		DestinationPortEnd:      strconv.Itoa(rule.DestinationPortEnd),
		Direction:               rule.Direction,
		Family:                  rule.Family,
		ICMPType:                rule.ICMPType,
		Position:                rule.Position,
		Protocol:                rule.Protocol,
		SourceAddressStart:      rule.SourceAddressStart,
		SourceAddressEnd:        rule.SourceAddressEnd,
		SourcePortStart:         strconv.Itoa(rule.SourcePortStart),
		SourcePortEnd:           strconv.Itoa(rule.SourcePortEnd),
	}
}

// A horrible copy of the upcloud Server definition, only to add the yml parsing definitions
type Yml_UpcloudFactory_ServerDefinition struct {
	AvoidHost        string                                        `yaml:"AvoidHost,omitempty"`
	BootOrder        string                                        `yaml:"BootOrder,omitempty"`
	CoreNumber       int                                           `yaml:"CoreNumber,omitempty"`
	Hostname         string                                        `yaml:"Hostname"`
	Networks         []Yml_UpcloudFactory_ServerDefinition_Network `yaml:"Networks"`
	LoginUser        Yml_UpcloudFactory_ServerDefinition_User      `yaml:"Userser,omitempty"`
	MemoryAmount     int                                           `yaml:"Memory,omitempty"`
	PasswordDelivery string                                        `yaml:"PasswordDelivery,omitempty"`
	Plan             string                                        `yaml:"Plan,omitempty"`
	StorageDevices   []Yml_UpcloudFactory_ServerDefinition_Storage `yaml:"Storage"`
	TimeZone         string                                        `yaml:"Timezone,omitempty"`
	Title            string                                        `yaml:"Id"`
	UserData         string                                        `yaml:"UserData,omitempty"`
	VideoModel       string                                        `yaml:"VideoModel,omitempty"`
	VNC              bool                                          `yaml:"Vnc,omitempty"`
	VNCPassword      string                                        `yaml:"VncPassword,omitempty"`
	Zone             string                                        `yaml:"Zone"`
}

// Build an upcloud CreateServerReequest
func (server *Yml_UpcloudFactory_ServerDefinition) CreateServerRequest() upcloud_request.CreateServerRequest {
	request := upcloud_request.CreateServerRequest{
		AvoidHost:        server.AvoidHost,
		BootOrder:        server.BootOrder,
		CoreNumber:       server.CoreNumber,
		Hostname:         server.Hostname,
		MemoryAmount:     server.MemoryAmount,
		PasswordDelivery: server.PasswordDelivery,
		Plan:             server.Plan,
		TimeZone:         server.TimeZone,
		Title:            server.Title,
		UserData:         server.UserData,
		VideoModel:       server.VideoModel,
		VNCPassword:      server.VNCPassword,
		Zone:             server.Zone,
		Firewall:         convertBoolToString(false, "onoff"), // enable this below if any firewall rules were provided
		VNC:              convertBoolToString(server.VNC, "onoff"),
		IPAddresses:      []upcloud_request.CreateServerIPAddress{},
		StorageDevices:   []upcloud.CreateServerStorageDevice{},
	}

	request.LoginUser = &upcloud_request.LoginUser{
		CreatePassword: convertBoolToString(server.LoginUser.CreatePassword, "yesno"),
		Username:       server.LoginUser.Username,
		SSHKeys:        server.LoginUser.SSHKeys,
	}

	if len(server.Networks) > 0 {
		request.Firewall = convertBoolToString(true, "onoff")
		for _, networkHolder := range server.Networks {
			request.IPAddresses = append(request.IPAddresses, upcloud_request.CreateServerIPAddress{
				Access: networkHolder.Access,
				Family: networkHolder.Family,
			})
		}
	}

	for _, StorageDeviceHolder := range server.StorageDevices {
		request.StorageDevices = append(request.StorageDevices, upcloud.CreateServerStorageDevice{
			Action:  StorageDeviceHolder.Action,
			Address: StorageDeviceHolder.Address,
			Storage: StorageDeviceHolder.Storage,
			Title:   StorageDeviceHolder.Title,
			Size:    StorageDeviceHolder.Size,
			Tier:    StorageDeviceHolder.Tier,
			Type:    StorageDeviceHolder.Type,
		})
	}

	// Some safe values / sanity checks

	if request.Hostname == "" {
		request.Hostname = request.Title
	}

	return request
}

// Small function to convert a boolean value to the upcloud string boolean value used.
func convertBoolToString(value bool, format string) string {
	if value {
		switch format {
		case "onoff":
			return "on"
		default: 
			return "yes"
		}
	} else {
		switch format {
		case "onoff":
			return "off"
		default: 
			return "no"
		}
	}
}

type Yml_UpcloudFactory_ServerDefinition_User struct {
	CreatePassword bool     `yaml:"CreatePassword,omitempty"`
	Username       string   `yaml:"Username,omitempty"`
	SSHKeys        []string `yaml:"SSHKeys,omitempty"`
}

type Yml_UpcloudFactory_ServerDefinition_Network struct {
	Access string `yaml:"Access"`
	Family string `yaml:"Family"`
}

type Yml_UpcloudFactory_ServerDefinition_Storage struct {
	Action  string `yaml:"Action,omitempty"`
	Address string `yaml:"Address,omitempty"`
	Storage string `yaml:"Storage"`
	Title   string `yaml:"Title,omitempty"`
	Size    int    `yaml:"Size,omitempty"`
	Tier    string `yaml:"Tier,omitempty"`
	Type    string `yaml:"Type,omitempty"`
}
