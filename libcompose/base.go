package libcompose

import (
	"errors"
	"io"

	"golang.org/x/net/context"

	api_operation "github.com/wunderkraut/radi-api/operation"
	handlers_bytesource "github.com/wunderkraut/radi-handlers/bytesource"
)

/**
 * Some usefull Base classes used by other libcompose operations
 * and Properties
 */

/**
 * Handlers
 */

// A libcompose handler base that can produce a BaseLibcomposeNameFilesOperation for ops base
type BaseLibcomposeHandler struct {
	LibComposeBaseOp *BaseLibcomposeNameFilesOperation
}

// Constructor for BaseLibcomposeHandler
func New_BaseLibcomposeHandler(projectName string, dockerComposeFiles []string, runContext context.Context, outputWriter io.Writer, errorWriter io.Writer, filesettings handlers_bytesource.BytesourceFileSettings) *BaseLibcomposeHandler {
	baseLibcomposeOp, _ := New_BaseLibcomposeNameFilesOperation(projectName, dockerComposeFiles, runContext, outputWriter, errorWriter, filesettings)
	base := &BaseLibcomposeHandler{LibComposeBaseOp: &baseLibcomposeOp}
	return base
}

/**
 * Operations
 */

// A base libcompose operation with Properties for project-name, and yml files
type BaseLibcomposeNameFilesOperation struct {
	projectName        string
	dockerComposeFiles []string
	runContext         context.Context
	outputWriter       io.Writer
	errorWriter        io.Writer
	filesettings       handlers_bytesource.BytesourceFileSettings
}

// Constructor for BaseLibcomposeNameFilesOperation
func New_BaseLibcomposeNameFilesOperation(projectName string, dockerComposeFiles []string, runContext context.Context, outputWriter io.Writer, errorWriter io.Writer, filesettings handlers_bytesource.BytesourceFileSettings) *BaseLibcomposeNameFilesOperation {
	return &BaseLibcomposeNameFilesOperation{
		projectName:        projectName,
		dockerComposeFiles: dockerComposeFiles,
		runContext:         runContext,
		outputWriter:       outputWriter,
		fileSettings:       fileSettings,
	}
}

// Provide static Properties for the operation - set values from the default
func (base *BaseLibcomposeNameFilesOperation) Properties() api_operation.Properties {
	props := &api_operation.Properties{}

	projectName := LibcomposeProjectnameProperty{}
	projectName.Set(base.projectName)
	props.Add(api_operation.Property(&projectName))

	fileSettings := handlers_bytesource.BytesourceFilesettingsProperty{}
	fileSettings.Set(base.fileSettings)
	props.Add(api_operation.Property(&fileSettings))

	composeFiles := LibcomposeComposefilesProperty{}
	composeFiles.Set(base.dockerComposeFiles)
	props.Add(api_operation.Property(&composeFiles))

	runContext := LibcomposeContextProperty{}
	runContext.Set(base.runContext)
	props.Add(api_operation.Property(&runContext))

	output := LibcomposeOutputProperty{}
	output.Set(base.outputWriter)
	props.Add(api_operation.Property(&output))
	err := LibcomposeErrorProperty{}
	err.Set(base.errorWriter)
	props.Add(api_operation.Property(&err))

	return props
}
