package templates

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/adapter"

	"github.com/nodeset-org/hyperdrive/modules"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
)

const (
	CallConfigFunctionCommandString string = adapter.HyperdriveModuleCommand + " call-config-function"
)

type CallConfigFunctionResponse struct {
	Result string `json:"result"`
}

// The data source for module service templates
type ServiceDataSource struct {
	// Public parameters
	ModuleComposeProject string
	ModuleNetwork        string
	ModuleConfigDir      string
	ModuleLogDir         string
	ModuleDataDir        string
	HyperdriveDaemonUrl  string
	HyperdriveJwtKeyFile string

	// Internal fields
	hyperdriveSettings *modconfig.ModuleSettings
	moduleSettingsMap  map[string]*modconfig.ModuleSettings
	moduleInfo         *modconfig.ModuleInfo
}

// Create a new service data source
func NewServiceDataSource(
	hdSettings *modconfig.ModuleSettings,
	moduleSettingsMap map[string]*modconfig.ModuleSettings,
	moduleInfo *modconfig.ModuleInfo,
	adapterSource *AdapterDataSource,
) *ServiceDataSource {
	return &ServiceDataSource{
		ModuleComposeProject: adapterSource.ModuleComposeProject,
		ModuleNetwork:        adapterSource.ModuleNetwork,
		ModuleConfigDir:      adapterSource.ModuleConfigDir,
		ModuleLogDir:         adapterSource.ModuleLogDir,
		//ModuleDataDir:        adapterSource.ModuleDataDir, TODO!
		//HyperdriveDaemonUrl:  hdSettings.DaemonUrl,
		//HyperdriveJwtKeyFile: hdSettings.JwtKeyFile,

		hyperdriveSettings: hdSettings,
		moduleSettingsMap:  moduleSettingsMap,
		moduleInfo:         moduleInfo,
	}
}

// Get the value of a property from its fully qualified path name
func (t *ServiceDataSource) GetValue(fqpn string) (string, error) {
	return t.getPropertyValue(fqpn)
}

// Get the value of a property from its fully qualified path name, splitting it into an array using the delimiter
func (t *ServiceDataSource) GetValueArray(fqpn string, delimiter string) ([]string, error) {
	val, err := t.getPropertyValue(fqpn)
	if err != nil {
		return nil, err
	}
	return strings.Split(val, delimiter), nil
}

func (t *ServiceDataSource) CallConfigFunction(funcName string) (string, error) {
	moduleDir := t.ModuleConfigDir
	adapterKeyPath := filepath.Join(moduleDir, shared.SecretsDir, shared.AdapterKeyFile)
	bytes, err := os.ReadFile(adapterKeyPath)
	containerName := t.moduleInfo.Descriptor.Name
	c, err := adapter.NewAdapterClient(string(containerName), string(bytes))
	if err != nil {
		return "", fmt.Errorf("error creating adapter client: %w", err)
	}

	req := map[string]any{
		"funcName": funcName,
	}
	var response CallConfigFunctionResponse
	err = adapter.RunCommand[map[string]any, CallConfigFunctionResponse](
		c, context.Background(), CallConfigFunctionCommandString, &req, &response)
	if err != nil {
		return "", fmt.Errorf("error calling config function \"%s\": %w", funcName, err)
	}
	return response.Result, nil
}

// Get the value of a property from its fully qualified path name
func (t *ServiceDataSource) getPropertyValue(fqpn string) (string, error) {
	fqmn := ""
	propertyPath := ""
	elements := strings.Split(fqpn, ":")
	if len(elements) == 1 {
		// This is a local property so use the module's fully qualified module name
		fqmn = t.moduleInfo.Descriptor.GetFullyQualifiedModuleName()
		propertyPath = fqpn
	} else {
		// TODO: Error out if there are more than 2 elements?
		fqmn = elements[0]
		propertyPath = elements[1]
	}

	// Get the module settings
	var settings *modconfig.ModuleSettings
	if fqmn == modules.HyperdriveFqmn {
		settings = t.hyperdriveSettings
	} else {
		var exists bool
		settings, exists = t.moduleSettingsMap[fqmn]
		if !exists {
			return "", fmt.Errorf("module settings not found for module \"%s\" in path \"%s\"", fqmn, propertyPath)
		}
	}
	return getModulePropertyValue(settings, propertyPath)
}

// Get the value of a module settings property from its path
func getModulePropertyValue(settings *modconfig.ModuleSettings, paramPath string) (string, error) {
	// Split the param path into its components
	elements := strings.Split(paramPath, "/")
	var container modconfig.IInstanceContainer = settings

	// Iterate through the sections
	level := 0
	for level < len(elements)-1 {
		elementString := elements[level]
		var id modconfig.Identifier
		err := id.UnmarshalText([]byte(elementString))
		if err != nil {
			return "", fmt.Errorf("error converting section \"%s\" in path \"%s\" to identifier: %w", elementString, paramPath, err)
		}
		container, err = container.GetSection(id)
		if err != nil {
			return "", fmt.Errorf("error getting section \"%s\" in path \"%s\": %w", elementString, paramPath, err)
		}
		level++
	}

	// Get the parameter value
	elementString := elements[level]
	var id modconfig.Identifier
	err := id.UnmarshalText([]byte(elementString))
	if err != nil {
		return "", fmt.Errorf("error converting parameter \"%s\" in path \"%s\" to identifier: %w", elementString, paramPath, err)
	}
	param, err := container.GetParameter(id)
	if err != nil {
		return "", fmt.Errorf("error getting parameter \"%s\" in path \"%s\": %w", elementString, paramPath, err)
	}
	return param.String(), nil
}
