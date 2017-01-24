package rancher

/**
 * Some base structs to share upcloud functionality
 */

// Shared base handler
type RancherBaseHandler struct {
	settings RancherSettings

	operations *api_operation.Operations
}

// Constructor for RancherBaseHandler
func New_RancherBaseHandler(settings RancherSettings) *RancherBaseHandler {
	return &RancherBaseHandler{
		settings: settings,
		operations: &api_operation.Opeartions{},
	}
}

// Get the operations from the handler
func (base *RancherBaseHandler) Operations() *api_operation.Operations {
	if base.operations == nil {
		return &api_operation.Operations{}
	} else {
		return base.operations
	}
}

// Retrieve the base settings
func (base *RnacherBaseHandler) Settings() RancherSettings {
	return base.settings
}
