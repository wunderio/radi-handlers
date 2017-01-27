package rancher

import (
	log "github.com/Sirupsen/logrus"
	
	rancher_client "github.com/rancher/go-rancher/client"
)

/**
 * Properties used to configure the rancher connections
 */

const (
	OPERATION_PROPERTY_RANCHER_CLIENTOPTS = "rancher.clientopts"
)

// A single Rancher client options configuration
type RancherClientOptsProperty struct {
	options rancher_client.ClientOpts
}

// Id for the Property
func (clientOpts *RancherClientOptsProperty) Id() string {
	return OPERATION_PROPERTY_RANCHER_CLIENTOPTS
}

// Label for the Property
func (clientOpts *RancherClientOptsProperty) Label() string {
	return "Rancher client opts"
}

// Description for the Property
func (clientOpts *RancherClientOptsProperty) Description() string {
	return "ClientOpts for the rancher client.  See github.com/rancher/go-rancher/client/common.go"
}

// Is the Property internal only
func (clientOpts *RancherClientOptsProperty) Internal() bool {
	return true
}

// Give an idea of what type of value the property consumes
func (clientOpts *RancherClientOptsProperty) Type() string {
	return "go-rancher.clientopts"
}

func (clientOpts *RancherClientOptsProperty) Get() interface{} {
	return interface{}(clientOpts.options)
}
func (clientOpts *RancherClientOptsProperty) Set(value interface{}) bool {
	if converted, ok := value.(rancher_client.ClientOpts); ok {
		clientOpts.options = converted
		return true
	} else {
		log.WithFields(log.Fields{"value": value}).Error("Could not assign Property value, because the passed parameter was the wrong type. Expected go-rancher client ClientOpts")
		return false
	}
}
