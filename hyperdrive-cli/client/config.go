package client

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client/template"
)

const (
	metricsDirMode           os.FileMode = 0755
	prometheusConfigTemplate string      = "prometheus-cfg.tmpl"
	prometheusConfigTarget   string      = "prometheus.yml"
	grafanaConfigTemplate    string      = "grafana-prometheus-datasource.tmpl"
	grafanaConfigTarget      string      = "grafana-prometheus-datasource.yml"
)

// Load the config
// Returns the global config and whether or not it was newly generated
func (c *HyperdriveClient) LoadConfig() (*GlobalConfig, bool, error) {
	if c.cfg != nil {
		return c.cfg, c.isNewCfg, nil
	}

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
		c.cfg = cfg
		return cfg, false, nil
	}

	// Config wasn't loaded, but there was no error- we should create one.
	hdCfg := hdconfig.NewHyperdriveConfig(c.Context.ConfigPath)
	c.cfg = NewGlobalConfig(hdCfg)
	c.isNewCfg = true
	return c.cfg, true, nil
}

// Load the backup config
func (c *HyperdriveClient) LoadBackupConfig() (*GlobalConfig, error) {
	settingsFilePath := filepath.Join(c.Context.ConfigPath, BackupSettingsFile)
	expandedPath, err := homedir.Expand(settingsFilePath)
	if err != nil {
		return nil, fmt.Errorf("error expanding backup settings file path: %w", err)
	}

	return LoadConfigFromFile(expandedPath)
}

// Save the config
func (c *HyperdriveClient) SaveConfig(cfg *GlobalConfig) error {
	settingsFileDirectoryPath, err := homedir.Expand(c.Context.ConfigPath)
	if err != nil {
		return err
	}
	err = SaveConfig(cfg, settingsFileDirectoryPath, SettingsFile)
	if err != nil {
		return fmt.Errorf("error saving config: %w", err)
	}

	// Update the client's config cache
	c.cfg = cfg
	c.isNewCfg = false
	return nil
}

// Create the metrics and modules folders, and deploy the config templates for Prometheus and Grafana
func (c *HyperdriveClient) DeployMetricsConfigurations(config *GlobalConfig) error {
	// Make sure the metrics path exists
	metricsDirPath := filepath.Join(c.Context.ConfigPath, metricsDir)
	modulesDirPath := filepath.Join(metricsDirPath, hdconfig.ModulesName)
	err := os.MkdirAll(modulesDirPath, metricsDirMode)
	if err != nil {
		return fmt.Errorf("error creating metrics and modules directories [%s]: %w", modulesDirPath, err)
	}

	err = updatePrometheusConfiguration(config, metricsDirPath)
	if err != nil {
		return fmt.Errorf("error updating Prometheus configuration: %w", err)
	}
	err = updateGrafanaDatabaseConfiguration(config, metricsDirPath)
	if err != nil {
		return fmt.Errorf("error updating Grafana configuration: %w", err)
	}
	return nil
}

// Load the Prometheus config template, do a template variable substitution, and save it
func updatePrometheusConfiguration(config *GlobalConfig, metricsDirPath string) error {
	prometheusConfigTemplatePath, err := homedir.Expand(filepath.Join(templatesDir, prometheusConfigTemplate))
	if err != nil {
		return fmt.Errorf("error expanding Prometheus config template path: %w", err)
	}

	prometheusConfigTargetPath, err := homedir.Expand(filepath.Join(metricsDirPath, prometheusConfigTarget))
	if err != nil {
		return fmt.Errorf("error expanding Prometheus config target path: %w", err)
	}

	t := template.Template{
		Src: prometheusConfigTemplatePath,
		Dst: prometheusConfigTargetPath,
	}

	return t.Write(config)
}

// Load the Grafana config template, do a template variable substitution, and save it
func updateGrafanaDatabaseConfiguration(config *GlobalConfig, metricsDirPath string) error {
	grafanaConfigTemplatePath, err := homedir.Expand(filepath.Join(templatesDir, grafanaConfigTemplate))
	if err != nil {
		return fmt.Errorf("error expanding Grafana config template path: %w", err)
	}

	grafanaConfigTargetPath, err := homedir.Expand(filepath.Join(metricsDirPath, grafanaConfigTarget))
	if err != nil {
		return fmt.Errorf("error expanding Grafana config target path: %w", err)
	}

	t := template.Template{
		Src: grafanaConfigTemplatePath,
		Dst: grafanaConfigTargetPath,
	}

	return t.Write(config)
}
