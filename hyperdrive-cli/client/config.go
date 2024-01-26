package client

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client/template"
	"github.com/nodeset-org/hyperdrive/shared/config"
)

const (
	prometheusConfigTemplate string = "prometheus-cfg.tmpl"
	prometheusConfigTarget   string = "prometheus.yml"
	grafanaConfigTemplate    string = "grafana-prometheus-datasource.tmpl"
	grafanaConfigTarget      string = "grafana-prometheus-datasource.yml"
)

// Load the config
// Returns the RocketPoolConfig and whether or not it was newly generated
func (c *Client) LoadConfig() (*config.HyperdriveConfig, bool, error) {
	settingsFilePath := filepath.Join(c.Context.ConfigPath, SettingsFile)
	expandedPath, err := homedir.Expand(settingsFilePath)
	if err != nil {
		return nil, false, fmt.Errorf("error expanding settings file path: %w", err)
	}

	cfg, err := LoadConfigFromFile(expandedPath)
	if err != nil {
		return nil, false, err
	}

	if cfg != nil {
		// A config was loaded, return it now
		return cfg, false, nil
	}

	// Config wasn't loaded, but there was no error- we should create one.
	return config.NewHyperdriveConfig(c.Context.ConfigPath), true, nil
}

// Load the backup config
func (c *Client) LoadBackupConfig() (*config.HyperdriveConfig, error) {
	settingsFilePath := filepath.Join(c.Context.ConfigPath, BackupSettingsFile)
	expandedPath, err := homedir.Expand(settingsFilePath)
	if err != nil {
		return nil, fmt.Errorf("error expanding backup settings file path: %w", err)
	}

	return LoadConfigFromFile(expandedPath)
}

// Save the config
func (c *Client) SaveConfig(cfg *config.HyperdriveConfig) error {
	settingsFileDirectoryPath, err := homedir.Expand(c.Context.ConfigPath)
	if err != nil {
		return err
	}
	return SaveConfig(cfg, settingsFileDirectoryPath, SettingsFile)
}

// Remove the upgrade flag file
func (c *Client) RemoveUpgradeFlagFile() error {
	expandedPath, err := homedir.Expand(c.Context.ConfigPath)
	if err != nil {
		return err
	}
	return RemoveUpgradeFlagFile(expandedPath)
}

// Returns whether or not this is the first run of the configurator since a previous installation
func (c *Client) IsFirstRun() (bool, error) {
	expandedPath, err := homedir.Expand(c.Context.ConfigPath)
	if err != nil {
		return false, fmt.Errorf("error expanding settings file path: %w", err)
	}
	return IsFirstRun(expandedPath), nil
}

// Load the Prometheus config template, do a template variable substitution, and save it
func (c *Client) UpdatePrometheusConfiguration(config *config.HyperdriveConfig) error {
	prometheusConfigTemplatePath, err := homedir.Expand(filepath.Join(templatesDir, prometheusConfigTemplate))
	if err != nil {
		return fmt.Errorf("Error expanding Prometheus config template path: %w", err)
	}

	prometheusConfigTargetPath, err := homedir.Expand(filepath.Join(c.Context.ConfigPath, prometheusConfigTarget))
	if err != nil {
		return fmt.Errorf("Error expanding Prometheus config target path: %w", err)
	}

	t := template.Template{
		Src: prometheusConfigTemplatePath,
		Dst: prometheusConfigTargetPath,
	}

	return t.Write(config)
}

// Load the Grafana config template, do a template variable substitution, and save it
func (c *Client) UpdateGrafanaDatabaseConfiguration(config *config.HyperdriveConfig) error {
	grafanaConfigTemplatePath, err := homedir.Expand(filepath.Join(templatesDir, grafanaConfigTemplate))
	if err != nil {
		return fmt.Errorf("Error expanding Grafana config template path: %w", err)
	}

	grafanaConfigTargetPath, err := homedir.Expand(filepath.Join(c.Context.ConfigPath, grafanaConfigTarget))
	if err != nil {
		return fmt.Errorf("Error expanding Grafana config target path: %w", err)
	}

	t := template.Template{
		Src: grafanaConfigTemplatePath,
		Dst: grafanaConfigTargetPath,
	}

	return t.Write(config)
}
