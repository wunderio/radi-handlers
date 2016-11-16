package bytesource

import (
	api_operation "github.com/james-nesbitt/kraut-api/operation"
)

/**
 * Base Operation that has a BytesourceFilesettingsProperty
 * property
 */

func New_BaseBytesourceFilesettingsOperation(settings BytesourceFileSettings) *BaseBytesourceFilesettingsOperation {
	return &BaseBytesourceFilesettingsOperation{
		settings: settings,
	}
}

type BaseBytesourceFilesettingsOperation struct {
	settings   BytesourceFileSettings
	properties *api_operation.Properties
}

func (base *BaseBytesourceFilesettingsOperation) Properties() *api_operation.Properties {
	if base.properties == nil {
		settingsProp := BytesourceFilesettingsProperty{}
		settingsProp.Set(base.settings)

		base.properties = &api_operation.Properties{}
		base.properties.Add(api_operation.Property(&settingsProp))
	}
	return base.properties
}
