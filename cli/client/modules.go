package client

import (
	"fmt"

	hdconfig "github.com/nodeset-org/hyperdrive/config"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
)

func (c *HyperdriveClient) StartProjectAdapters(settings *hdconfig.HyperdriveSettings, modInfos map[string]*modconfig.ModuleInfo, moduleSettingsMap map[string]*modconfig.ModuleSettings) error {
	// Start all of the base services and project module adapters
	composeFiles, err := deployTemplates(c.Context.SystemDirPath, c.Context.UserDirPath, settings)
	if err != nil {
		return fmt.Errorf("error deploying templates: %w", err)
	}
	err = deployModules(c.modMgr, c.Context.ModulesDir(), settings, moduleSettingsMap, modInfos)
	if err != nil {
		return fmt.Errorf("error deploying modules: %w", err)
	}
	err = startComposeFiles(c.Context.UserDirPath, settings.ProjectName, modInfos, composeFiles)
	if err != nil {
		return fmt.Errorf("error starting project adapters: %w", err)
	}
	return nil
}
