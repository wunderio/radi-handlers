package bytesource

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
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

func (base *BaseBytesourceFilesettingsOperation) Properties() api_operation.Properties {
	settingsProp := BytesourceFilesettingsProperty{}
	settingsProp.Set(base.settings)

	props := api_operation.Properties{}
	props.Add(api_operation.Property(&settingsProp))

	return props
}
