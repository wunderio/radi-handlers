package bytesource

import (
	api_property "github.com/wunderkraut/radi-api/property"
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
	settings BytesourceFileSettings
}

func (base *BaseBytesourceFilesettingsOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	settingsProp := BytesourceFilesettingsProperty{}
	settingsProp.Set(base.settings)

	props.Add(api_property.Property(&settingsProp))

	return props.Properties()
}
