package client

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/context"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/urfave/cli/v2"
)

type HyperdriveClient struct {
	Context *context.HyperdriveContext

	cfgDir          string
	systemDir       string
	modMgr          *shared.ModuleManager
	cfgMgr          *config.ConfigurationManager
	primarySettings *config.HyperdriveSettings
}

// Create a new Hyperdrive client from the CLI context
func NewHyperdriveClientFromCtx(c *cli.Context) (*HyperdriveClient, error) {
	hdCtx := context.GetHyperdriveContext(c)
	if hdCtx == nil {
		return nil, fmt.Errorf("Hyperdrive CLI context has not been created")
	}
	cfgDir := hdCtx.UserDirPath
	systemDir := hdCtx.SystemDirPath

	// Config manager
	cfgMgr := config.NewConfigurationManager(cfgDir, systemDir)

	// Module manager
	//adapterKeyPath := shared.GetAdapterKeyPath(cfgDir)
	modulesDir := shared.GetModulesDirectoryPath(systemDir)
	modMgr, err := shared.NewModuleManager(modulesDir)
	if err != nil {
		return nil, fmt.Errorf("error creating module manager: %w", err)
	}

	return &HyperdriveClient{
		cfgDir:    cfgDir,
		systemDir: systemDir,
		modMgr:    modMgr,
		cfgMgr:    cfgMgr,
		Context:   hdCtx,
	}, nil
}

// Get the current Hyperdrive configuration
func (c *HyperdriveClient) GetHyperdriveConfiguration() *config.HyperdriveConfig {
	return c.cfgMgr.HyperdriveConfiguration
}

// Load the config settings
func (c *HyperdriveClient) LoadMainSettingsFile() (*config.HyperdriveSettings, bool, error) {
	primaryConfigPath := filepath.Join(c.cfgDir, config.ConfigFilename)
	cfg, err := c.cfgMgr.HyperdriveConfiguration.LoadSettingsFromFile(primaryConfigPath)
	if err != nil {
		return nil, false, fmt.Errorf("error loading main config settings file: %w", err)
	}

	if cfg == nil {
		return nil, false, nil
	}

	c.primarySettings = cfg
	return cfg, false, nil
}

// Load all of the module info and settings
func (c *HyperdriveClient) LoadModules() ([]*shared.ModuleInfoLoadResult, error) {
	results, err := c.modMgr.LoadModuleInfo(true)
	if err != nil {
		return nil, fmt.Errorf("error loading module info: %w", err)
	}
	for _, result := range results {
		if result.LoadError == nil {
			name := result.Info.Descriptor.GetFullyQualifiedModuleName()
			c.cfgMgr.HyperdriveConfiguration.Modules[name] = result.Info.ModuleInfo
		}
	}
	return results, nil
}

// Initializes a new installation of hyperdrive
func (c *HyperdriveClient) InitializeNewInstallation() error {
	// Make sure the config directory exists
	err := c.createDirectory(c.cfgDir, 0755)
	if err != nil {
		return err
	}

	// Make the subdirectories
	subdirs := []string{
		shared.LogsDir,
		shared.SecretsDir,
		shared.OverrideDir,
		shared.RuntimeDir,
		shared.MetricsDir,
	}
	for _, subdir := range subdirs {
		err = c.createDirectory(filepath.Join(c.cfgDir, subdir), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create a directory if it doesn't already exist
func (c HyperdriveClient) createDirectory(path string, mode os.FileMode) error {
	stat, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		err = os.MkdirAll(path, mode)
		if err != nil {
			return fmt.Errorf("error creating directory [%s]: %w", c.cfgDir, err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking directory [%s]: %w", c.cfgDir, err)
	} else if !stat.IsDir() {
		return fmt.Errorf("directory [%s] is already in use but is not a directory", c.cfgDir)
	}

	return nil
}
