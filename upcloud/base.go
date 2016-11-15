package upcloud

// Constructor for BaseUpcloudServiceHandler
func New_BaseUpcloudServiceHandler(service *UpcloudServiceWrapper) *BaseUpcloudServiceHandler {
	return &BaseUpcloudServiceHandler{
		service: service,
	}
}

// Base handler with an upcloud service
type BaseUpcloudServiceHandler struct {
	service *UpcloudServiceWrapper
}

// Set the service
func (base *BaseUpcloudServiceHandler) ServiceWrapper() *UpcloudServiceWrapper {
	return base.service
}

/**
 * Base operations for Upcloud operations, which
 * allow sharing of Upcloud service across instances
 */

// Constructor for BaseUpcloudServiceOperation
func New_BaseUpcloudServiceOperation(service *UpcloudServiceWrapper) *BaseUpcloudServiceOperation {
	return &BaseUpcloudServiceOperation{
		service: service,
	}
}

// Base operation with an upcloud service
type BaseUpcloudServiceOperation struct {
	service *UpcloudServiceWrapper
}

// Set the service
func (base *BaseUpcloudServiceOperation) ServiceWrapper() *UpcloudServiceWrapper {
	return base.service
}
