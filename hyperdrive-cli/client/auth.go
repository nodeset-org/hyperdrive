package client

import (
	"fmt"
	"os"
	"path/filepath"

	csconfig "github.com/nodeset-org/hyperdrive-constellation/shared/config"
	"github.com/nodeset-org/hyperdrive-daemon/shared/auth"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	swconfig "github.com/nodeset-org/hyperdrive-stakewise/shared/config"
)

const (
	authDirMode os.FileMode = 0700
)

var (
	hdApiKeyRelPath string = filepath.Join(hdconfig.SecretsDir, hdconfig.DaemonKeyFilename)
	swApiKeyRelPath string = filepath.Join(hdconfig.SecretsDir, hdconfig.ModulesName, swconfig.ModuleName, hdconfig.DaemonKeyFilename)
	csApiKeyRelPath string = filepath.Join(hdconfig.SecretsDir, hdconfig.ModulesName, csconfig.ModuleName, hdconfig.DaemonKeyFilename)
)

// Create the metrics and modules folders, and deploy the config templates for Prometheus and Grafana
func (c *HyperdriveClient) GenerateDaemonAuthKeys(config *GlobalConfig) error {
	// Make sure the secrets path exists
	secretsDirPath := filepath.Join(c.Context.UserDirPath, hdconfig.SecretsDir)
	modulesDirPath := filepath.Join(secretsDirPath, hdconfig.ModulesName)
	err := os.MkdirAll(modulesDirPath, authDirMode)
	if err != nil {
		return fmt.Errorf("error creating secrets and modules directories [%s]: %w", modulesDirPath, err)
	}

	// Create the API key for the Hyperdrive daemon
	hdApiKeyPath := filepath.Join(c.Context.UserDirPath, hdApiKeyRelPath)
	err = auth.GenerateAuthKeyIfNotPresent(hdApiKeyPath, auth.DefaultKeyLength)
	if err != nil {
		return fmt.Errorf("error generating Hyperdrive daemon API key: %w", err)
	}

	// Create the API key for the StakeWise module if enabled
	if config.StakeWise.Enabled.Value {
		swDirPath := filepath.Join(modulesDirPath, swconfig.ModuleName)
		err := os.MkdirAll(swDirPath, authDirMode)
		if err != nil {
			return fmt.Errorf("error creating StakeWise API key module directory [%s]: %w", swDirPath, err)
		}
		swApiKeyPath := filepath.Join(c.Context.UserDirPath, swApiKeyRelPath)
		err = auth.GenerateAuthKeyIfNotPresent(swApiKeyPath, auth.DefaultKeyLength)
		if err != nil {
			return fmt.Errorf("error generating StakeWise module API key: %w", err)
		}
	}

	// Create the API key for the Constellation module if enabled
	if config.Constellation.Enabled.Value {
		csDirPath := filepath.Join(modulesDirPath, csconfig.ModuleName)
		err := os.MkdirAll(csDirPath, authDirMode)
		if err != nil {
			return fmt.Errorf("error creating Constellation API key module directory [%s]: %w", csDirPath, err)
		}
		csApiKeyPath := filepath.Join(c.Context.UserDirPath, csApiKeyRelPath)
		err = auth.GenerateAuthKeyIfNotPresent(csApiKeyPath, auth.DefaultKeyLength)
		if err != nil {
			return fmt.Errorf("error generating Constellation module API key: %w", err)
		}
	}
	return nil
}
