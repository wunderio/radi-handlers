package upcloud

import (
	upcloud_client "github.com/Jalle19/upcloud-go-sdk/upcloud/client"
)

// Constructor for UpcloudClientSettings
func New_UpcloudClientSettings(user string, password string) *UpcloudClientSettings {
	return &UpcloudClientSettings{
		user:     user,
		password: password,
	}
}

// Client Settings property
type UpcloudClientSettings struct {
	user     string
	password string // don't make this public
}

// Builder an UpCloud client from the settings
func (settings *UpcloudClientSettings) Client() *upcloud_client.Client {
	return upcloud_client.New(settings.user, settings.password)
}
