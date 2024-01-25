package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/nodeset-org/hyperdrive/shared/types"
)

// ===============
// === General ===
// ===============

// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) IsLocalMode() bool {
	return cfg.ClientMode.Value == types.ClientMode_Local
}

// ========================
// === Execution Client ===
// ========================

// Gets the port mapping of the ec container
// Used by text/template to format ec.yml
func (cfg *HyperdriveConfig) GetEcOpenApiPorts() string {
	return cfg.LocalExecutionConfig.getOpenApiPortMapping()
}

// Gets the max peers of the ec container
// Used by text/template to format ec.yml
func (cfg *HyperdriveConfig) GetEcMaxPeers() (uint16, error) {
	if cfg.ClientMode.Value != types.ClientMode_Local {
		return 0, fmt.Errorf("Execution client is external, there is no max peers")
	}
	return cfg.LocalExecutionConfig.getMaxPeers(), nil
}

// Gets the tag of the ec container
// Used by text/template to format ec.yml
func (cfg *HyperdriveConfig) GetEcContainerTag() (string, error) {
	if cfg.ClientMode.Value != types.ClientMode_Local {
		return "", fmt.Errorf("Execution client is external, there is no container tag")
	}
	return cfg.LocalExecutionConfig.getContainerTag(), nil
}

// Used by text/template to format ec.yml
func (cfg *HyperdriveConfig) GetEcAdditionalFlags() (string, error) {
	if cfg.ClientMode.Value != types.ClientMode_Local {
		return "", fmt.Errorf("Execution client is external, there are no additional flags")
	}
	return cfg.LocalExecutionConfig.getAdditionalFlags(), nil
}

// Used by text/template to format ec.yml
func (cfg *HyperdriveConfig) GetExternalIP() string {
	// Get the external IP address
	ip, err := getExternalIP()
	if err != nil {
		fmt.Println("Warning: couldn't get external IP address; if you're using Nimbus or Besu, it may have trouble finding peers:")
		fmt.Println(err.Error())
		return ""
	}

	if ip.To4() == nil {
		fmt.Println("Warning: external IP address is v6; if you're using Nimbus or Besu, it may have trouble finding peers:")
	}

	return ip.String()
}

// ===================
// === Beacon Node ===
// ===================

// Gets the tag of the bn container
// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetBnContainerTag() (string, error) {
	if cfg.ClientMode.Value != types.ClientMode_Local {
		return "", fmt.Errorf("Beacon Node is external, there is no container tag")
	}
	return cfg.LocalBeaconConfig.getContainerTag(), nil
}

// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetBnOpenPorts() []string {
	return cfg.LocalBeaconConfig.getOpenApiPortMapping()
}

// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetEcHttpEndpoint() string {
	if cfg.ClientMode.Value != types.ClientMode_Local {
		return fmt.Sprintf("http://%s:%d", types.ContainerID_ExecutionClient, cfg.LocalExecutionConfig.HttpPort.Value)
	}

	return cfg.ExternalExecutionConfig.HttpUrl.Value
}

// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetEcWsEndpoint() string {
	if cfg.ClientMode.Value != types.ClientMode_Local {
		return fmt.Sprintf("ws://%s:%d", types.ContainerID_ExecutionClient, cfg.LocalExecutionConfig.WebsocketPort.Value)
	}

	return cfg.ExternalExecutionConfig.WebsocketUrl.Value
}

// Gets the max peers of the bn container
// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetBnMaxPeers() (uint16, error) {
	if cfg.ClientMode.Value != types.ClientMode_Local {
		return 0, fmt.Errorf("Beacon Node is external, there is no max peers")
	}
	return cfg.LocalBeaconConfig.getMaxPeers(), nil
}

// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetBnAdditionalFlags() (string, error) {
	if cfg.ClientMode.Value != types.ClientMode_Local {
		return "", fmt.Errorf("Beacon Node is external, there is no additional flags")
	}
	return cfg.LocalBeaconConfig.getAdditionalFlags(), nil
}

// ===============
// === Metrics ===
// ===============

// Used by text/template to format exporter.yml
func (cfg *HyperdriveConfig) GetExporterAdditionalFlags() []string {
	flags := strings.Trim(cfg.Metrics.Exporter.AdditionalFlags.Value, " ")
	if flags == "" {
		return nil
	}
	return strings.Split(flags, " ")
}

// Used by text/template to format prometheus.yml
func (cfg *HyperdriveConfig) GetPrometheusAdditionalFlags() []string {
	flags := strings.Trim(cfg.Metrics.Prometheus.AdditionalFlags.Value, " ")
	if flags == "" {
		return nil
	}
	return strings.Split(flags, " ")
}

// Used by text/template to format prometheus.yml
func (cfg *HyperdriveConfig) GetPrometheusOpenPorts() string {
	portMode := cfg.Metrics.Prometheus.OpenPort.Value
	if !portMode.IsOpen() {
		return ""
	}
	return fmt.Sprintf("\"%s\"", portMode.DockerPortMapping(cfg.Metrics.Prometheus.Port.Value))
}

// Gets the hostname portion of the Execution Client's URI.
// Used by text/template to format prometheus.yml.
func (cfg *HyperdriveConfig) GetExecutionHostname() (string, error) {
	if cfg.ClientMode.Value == types.ClientMode_Local {
		return string(types.ContainerID_ExecutionClient), nil
	}
	ecUrl, err := url.Parse(cfg.ExternalExecutionConfig.HttpUrl.Value)
	if err != nil {
		return "", fmt.Errorf("Invalid External Execution URL %s: %w", cfg.ExternalExecutionConfig.HttpUrl.Value, err)
	}

	return ecUrl.Hostname(), nil
}

// Gets the hostname portion of the Beacon Node's URI.
// Used by text/template to format prometheus.yml.
func (cfg *HyperdriveConfig) GetBeaconHostname() (string, error) {
	if cfg.ClientMode.Value == types.ClientMode_Local {
		return string(types.ContainerID_BeaconNode), nil
	}
	ccUrl, err := url.Parse(cfg.ExternalBeaconConfig.HttpUrl.Value)
	if err != nil {
		return "", fmt.Errorf("Invalid External Consensus URL %s: %w", cfg.ExternalBeaconConfig.HttpUrl.Value, err)
	}

	return ccUrl.Hostname(), nil
}
