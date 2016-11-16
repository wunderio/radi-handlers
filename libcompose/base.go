package libcompose

import (
	"errors"
	"io"

	"golang.org/x/net/context"

	api_operation "github.com/james-nesbitt/kraut-api/operation"
	handlers_bytesource "github.com/james-nesbitt/kraut-handlers/bytesource"
)

/**
 * Some usefull Base classes used by other libcompose operations
 * and Properties
 */

/**
 * Handlers
 */

func New_BaseLibcomposeHandler(projectName string, dockerComposeFiles []string, runContext context.Context, outputWriter io.Writer, errorWriter io.Writer, filesettings handlers_bytesource.BytesourceFileSettings) *BaseLibcomposeHandler {
	baseLibcomposeOp, _ := New_BaseLibcomposeNameFilesOperation(projectName, dockerComposeFiles, runContext, outputWriter, errorWriter, filesettings)
	base := &BaseLibcomposeHandler{LibComposeBaseOp: &baseLibcomposeOp}
	return base
}

// A libcompose handler base that can produce a BaseLibcomposeNameFilesOperation for ops base
type BaseLibcomposeHandler struct {
	LibComposeBaseOp *BaseLibcomposeNameFilesOperation
}

/**
 * Operations
 */

// A handoff function to make a base orchestration operation, which is really just a lot of linear code.
func New_BaseLibcomposeNameFilesOperation(projectName string, dockerComposeFiles []string, runContext context.Context, outputWriter io.Writer, errorWriter io.Writer, filesettings handlers_bytesource.BytesourceFileSettings) (BaseLibcomposeNameFilesOperation, api_operation.Result) {
	result := api_operation.BaseResult{}
	result.Set(true, nil)

	// This Base operations will be at the root of all of the libCompose operations
	baseLibcomposeOrchestrate := BaseLibcomposeNameFilesOperation{}
	orchestrateProperties := baseLibcomposeOrchestrate.Properties()

	// Set a project name
	if projectNameConf, found := orchestrateProperties.Get(OPERATION_PROPERTY_LIBCOMPOSE_PROJECTNAME); found {
		if !projectNameConf.Set(projectName) {
			result.Set(false, []error{errors.New("Could not set base libCompose project name.  Config set error on base Orchestration operation")})
		}
	} else {
		result.Set(false, []error{errors.New("Could not set base libCompose project name.  Config value not found on base Orchestration operation")})
	}

	// Add project context
	if projectFilesettingsConf, found := orchestrateProperties.Get(handlers_bytesource.OPERATION_PROPERTY_BYTESOURCE_FILESETTINGS); found {
		if !projectFilesettingsConf.Set(filesettings) {
			result.Set(false, []error{errors.New("Could not set base libcompose file settings.  Config set error on base Orchestration operation")})
		}
	} else {
		result.Set(false, []error{errors.New("Could not set base libcompose file settings.  Config not found on base Orchestration operation")})
	}

	// Add project docker-compose yml files
	if projectComposeFilesConf, found := orchestrateProperties.Get(OPERATION_PROPERTY_LIBCOMPOSE_COMPOSEFILES); found {
		if !projectComposeFilesConf.Set(dockerComposeFiles) {
			result.Set(false, []error{errors.New("Could not set base libcompose docker-compose file conf.  Config set error on base Orchestration operation")})
		}
	} else {
		result.Set(false, []error{errors.New("Could not set base libcompose docker-compose file conf.  Config not found on base Orchestration operation")})
	}
	// Add project context
	if projectContextConf, found := orchestrateProperties.Get(OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT); found {
		if !projectContextConf.Set(runContext) {
			result.Set(false, []error{errors.New("Could not set base libcompose net context.  Config set error on base Orchestration operation")})
		}
	} else {
		result.Set(false, []error{errors.New("Could not set base libcompose net context.  Config not found on base Orchestration operation")})
	}
	// Add Stdout as an output writer
	if projectOutputConf, found := orchestrateProperties.Get(OPERATION_PROPERTY_LIBCOMPOSE_OUTPUT); found {
		if !projectOutputConf.Set(outputWriter) {
			result.Set(false, []error{errors.New("Could not set base libcompose output handler.  Config set error on base Orchestration operation")})
		}
	} else {
		result.Set(false, []error{errors.New("Could not set base libcompose output handler.  Config not found on base Orchestration operation")})
	}
	if projectErrorConf, found := orchestrateProperties.Get(OPERATION_PROPERTY_LIBCOMPOSE_ERROR); found {
		if !projectErrorConf.Set(errorWriter) {
			result.Set(false, []error{errors.New("Could not set base libcompose error handler.  Config set error on base Orchestration operation")})
		}
	} else {
		result.Set(false, []error{errors.New("Could not set base libcompose error handler.  Config not found on base Orchestration operation")})
	}

	return baseLibcomposeOrchestrate, api_operation.Result(&result)
}

// A base libcompose operation with Properties for project-name, and yml files
type BaseLibcomposeNameFilesOperation struct {
	properties *api_operation.Properties
}

// Provide static Properties for the operation
func (base *BaseLibcomposeNameFilesOperation) Properties() *api_operation.Properties {
	if base.properties == nil {
		newProperties := &api_operation.Properties{}

		newProperties.Add(api_operation.Property(&LibcomposeProjectnameProperty{}))

		newProperties.Add(api_operation.Property(&handlers_bytesource.BytesourceFilesettingsProperty{}))

		newProperties.Add(api_operation.Property(&LibcomposeComposefilesProperty{}))
		newProperties.Add(api_operation.Property(&LibcomposeContextProperty{}))

		newProperties.Add(api_operation.Property(&LibcomposeOutputProperty{}))
		newProperties.Add(api_operation.Property(&LibcomposeErrorProperty{}))

		base.properties = newProperties
	}
	return base.properties
}
