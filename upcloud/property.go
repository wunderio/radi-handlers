package upcloud

import (
	api_operation "github.com/james-nesbitt/kraut-api/operation"
)

/**
 * Custom properties for the upcloud operations
 */

const (
	UPCLOUD_GLOBAL_PROPERTY       = "upcloud.global"
	UPCLOUD_SERVER_UUID_PROPERTY  = "upcloud.server.uuid"
	UPCLOUD_STORAGE_UUID_PROPERTY = "upcloud.storage.uuid"
	UPCLOUD_ZONE_ID_PROPERTY      = "upcloud.zone.id"
)

// A boolean flag that tells upcloud to consider services/zones outside the scope of the project
// @NOTE this is kind of risky to use, so it should be limited to safe operations
type UpcloudGlobalProperty struct {
	api_operation.BooleanProperty
}

// ID returns string unique property Identifier
func (global *UpcloudGlobalProperty) Id() string {
	return UPCLOUD_GLOBAL_PROPERTY
}

// Label returns a short user readable label for the property
func (global *UpcloudGlobalProperty) Label() string {
	return "Global UpCloud services"
}

// Description provides a longer multi-line string description of what the property does
func (global *UpcloudGlobalProperty) Description() string {
	return "Consider UpCloud Service/Zones outside the scope of the project"
}

// Mark a property as being for internal use only (no shown to users)
func (global *UpcloudGlobalProperty) Internal() bool {
	return false
}

// A string slice property to match to server UUID
type UpcloudServerUUIDProperty struct {
	api_operation.StringSliceProperty
}

// ID returns string unique property Identifier
func (uuid *UpcloudServerUUIDProperty) Id() string {
	return UPCLOUD_SERVER_UUID_PROPERTY
}

// Label returns a short user readable label for the property
func (uuid *UpcloudServerUUIDProperty) Label() string {
	return "UpCloud server UUID"
}

// Description provides a longer multi-line string description of what the property does
func (uuid *UpcloudServerUUIDProperty) Description() string {
	return "Specific UpCloud server UUID"
}

// Mark a property as being for internal use only (no shown to users)
func (uuid *UpcloudServerUUIDProperty) Internal() bool {
	return false
}

// A string slice property to match to storage UUID
type UpcloudStorageUUIDProperty struct {
	api_operation.StringSliceProperty
}

// ID returns string unique property Identifier
func (uuid *UpcloudStorageUUIDProperty) Id() string {
	return UPCLOUD_STORAGE_UUID_PROPERTY
}

// Label returns a short user readable label for the property
func (uuid *UpcloudStorageUUIDProperty) Label() string {
	return "UpCloud server UUID"
}

// Description provides a longer multi-line string description of what the property does
func (uuid *UpcloudStorageUUIDProperty) Description() string {
	return "Specific UpCloud server UUID"
}

// Mark a property as being for internal use only (no shown to users)
func (uuid *UpcloudStorageUUIDProperty) Internal() bool {
	return false
}

// A string slice property to match to zone id
type UpcloudZoneIdProperty struct {
	api_operation.StringSliceProperty
}

// ID returns string unique property Identifier
func (id *UpcloudZoneIdProperty) Id() string {
	return UPCLOUD_ZONE_ID_PROPERTY
}

// Label returns a short user readable label for the property
func (id *UpcloudZoneIdProperty) Label() string {
	return "UpCloud zone ID"
}

// Description provides a longer multi-line string description of what the property does
func (id *UpcloudZoneIdProperty) Description() string {
	return "Specific UpCloud zone ID"
}

// Mark a property as being for internal use only (no shown to users)
func (id *UpcloudZoneIdProperty) Internal() bool {
	return false
}
