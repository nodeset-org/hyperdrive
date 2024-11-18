package client

import (
	"fmt"
	"path/filepath"

	"github.com/nodeset-org/hyperdrive/shared/auth"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
)

var (
	hdApiKeyRelPath     string = filepath.Join(hdconfig.SecretsDir, hdconfig.DaemonKeyFilename)
	moduleApiKeyRelPath string = filepath.Join(hdconfig.SecretsDir, hdconfig.ModulesName)
)

// Create the metrics and modules folders, and deploy the config templates for Prometheus and Grafana
func (c *HyperdriveClient) GenerateDaemonAuthKeys(config *GlobalConfig) error {
	// Create the API key for the Hyperdrive daemon
	hdApiKeyPath := filepath.Join(c.Context.UserDirPath, hdApiKeyRelPath)
	err := auth.GenerateAuthKeyIfNotPresent(hdApiKeyPath, auth.DefaultKeyLength)
	if err != nil {
		return fmt.Errorf("error generating Hyperdrive daemon API key: %w", err)
	}
	return nil
}
