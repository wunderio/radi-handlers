package local

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"os"
	"path"

	"github.com/james-nesbitt/kraut-api"	
	"github.com/james-nesbitt/kraut-api/handler"	
	"github.com/james-nesbitt/kraut-api/operation"
	"github.com/james-nesbitt/kraut-api/operation/command"
	"github.com/james-nesbitt/kraut-api/operation/config"
	"github.com/james-nesbitt/kraut-api/operation/orchestrate"
	"github.com/james-nesbitt/kraut-api/operation/setting"
	"github.com/james-nesbitt/kraut-handlers/bytesource"
	"github.com/james-nesbitt/kraut-handlers/libcompose"
	"github.com/james-nesbitt/kraut-handlers/null"
)

// Make a Local based API object, based on a project path
func MakeLocalAPI(settings LocalAPISettings) (*LocalAPI, error) {
	localApi := New_LocalAPI(settings)

	if settings.ProjectDoesntExist {

		// allow local project operations, which could be used to build a project
		localApi.BuildLocalProject()

		// Use null wrappers for those handlers that we don't have (to play safe)
		localApi.Config = config.ConfigWrapper(&null.NullConfigWrapper{})
		localApi.Settings = setting.SettingWrapper(&null.NullSettingWrapper{})
		localApi.Command = command.CommandWrapper(&null.NullCommandWrapper{})

	} else {

		// build all local operations
		localApi.BuildLocalConfigSettings()
		localApi.BuildLocalProject()
		localApi.BuildLocalOrchestrate()
		localApi.BuildLocalCommand()

	}

	return localApi, nil
}

// Settings needed to make a local API
type LocalAPISettings struct {
	bytesource.BytesourceFileSettings
	Context context.Context
}

// Constructor for LocalAPI
func New_LocalAPI(settings LocalAPISettings) *LocalAPI {
	return &LocalAPI{
		settings: settings,
	}
}

// An API based entirely on local handler
type LocalAPI struct {
	api.BaseAPI
	settings LocalAPISettings

	common_base *LocalHandler_Base
	common_libcompose *libcompose.BaseLibcomposeHandler

	Command     command.CommandWrapper
	Config      config.ConfigWrapper
	Settings    setting.SettingWrapper
	Orchestrate orchestrate.OrchestrateWrapper
}

// Validate the local API instance
func (localApi *LocalAPI) Validate() bool {
	return true
}

// Create a shareable common base
func (localApi *LocalAPI) base() *LocalHandler_Base {
	if localApi.common_base == nil {

		// Create a base handler, which just wraps the settings
		localApi.common_base = &LocalHandler_Base{
			settings:   &localApi.settings,
			operations: &operation.Operations{},
		}

	}

	return localApi.common_base
}

// Build a Handler base that produces LibCompose projects
func (localApi *LocalAPI) base_libcompose() *libcompose.BaseLibcomposeHandler {
	if localApi.common_libcompose == nil {
		// Set a project name
		projectName := "default"
		if settingsProjectName, err := localApi.Settings.Get("Project"); err == nil {
			projectName = settingsProjectName
		} else {
			log.WithError(errors.New("Could not set base libCompose project name.  Config value not found in handler config")).Error()
		}

		// Where to get docker-composer files
		dockerComposeFiles := []string{}
		// add the root composer file
		dockerComposeFiles = append(dockerComposeFiles, path.Join(localApi.settings.ProjectRootPath, "docker-compose.yml"))

		// What net context to use
		runContext := localApi.settings.Context

		// Output and Error writers
		outputWriter := os.Stdout
		errorWriter := os.Stderr

		// LibComposeHandlerBase
		localApi.common_libcompose = libcompose.New_BaseLibcomposeHandler(projectName, dockerComposeFiles, runContext, outputWriter, errorWriter, localApi.settings.BytesourceFileSettings)
	}

	return localApi.common_libcompose
}

// Add local Handlers for Config and Settings
func (localApi *LocalAPI) BuildLocalConfigSettings() error {
	// Build a config whandler
	local_config := LocalHandler_Config{
		LocalHandler_Base: *localApi.base(),
	}
	local_config.Init()
	localApi.AddHandler(handler.Handler(&local_config))
	// Get a config wrapper for other handlers
	localApi.Config = local_config.ConfigWrapper()

	// Build a settings handler which uses the configwrapper and the base
	local_settings := LocalHandler_Setting{
		LocalHandler_Base: *localApi.base(),
	}
	local_settings.SetConfigWrapper(localApi.Config)
	local_settings.Init()
	localApi.AddHandler(handler.Handler(&local_settings))
	// Get a settings wrapper for other handlers
	localApi.Settings = local_settings.SettingWrapper()

	return nil
}

// Make a Local based API object for "no existing project found" to allow for project operations
func (localApi *LocalAPI) BuildLocalProject() error {
	// Build a config wrapper using the base for settings
	local_project := LocalHandler_Project{
		LocalHandler_Base: *localApi.base(),
	}
	local_project.Init()
	localApi.AddHandler(handler.Handler(&local_project))

	return nil
}

// Add local Handlers for Orchestrate operations
func (localApi *LocalAPI) BuildLocalOrchestrate() error {
	// Build an orchestration handler
	local_orchestration := LocalHandler_Orchestrate{
		LocalHandler_Base:     *localApi.base(),
		BaseLibcomposeHandler: *localApi.base_libcompose(),
	}
	local_orchestration.SetSettingWrapper(localApi.Settings)
	local_orchestration.Init()
	localApi.AddHandler(handler.Handler(&local_orchestration))
	// Get an orchestrate wrapper for other handlers
	localApi.Orchestrate = local_orchestration.OrchestrateWrapper()

	return nil
}

// Add local Handlers for Command operations
func (localApi *LocalAPI) BuildLocalCommand() error {
	// Build a command Handler
	local_command := LocalHandler_Command{
		LocalHandler_Base:     *localApi.base(),
		BaseLibcomposeHandler: *localApi.base_libcompose(),
	}
	local_command.SetConfigWrapper(localApi.Config)
	local_command.Init()
	localApi.AddHandler(handler.Handler(&local_command))
	// Get an orchestrate wrapper for other handlers
	localApi.Command = local_command.CommandWrapper()

	return nil
}
