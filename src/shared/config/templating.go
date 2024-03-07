package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/rocket-pool/node-manager-core/config"
)

// =================
// === Constants ===
// =================

func (c *HyperdriveConfig) BeaconNodeContainerName() string {
	return string(config.ContainerID_BeaconNode)
}

func (c *HyperdriveConfig) DaemonContainerName() string {
	return string(config.ContainerID_Daemon)
}

func (c *HyperdriveConfig) ExecutionClientContainerName() string {
	return string(config.ContainerID_ExecutionClient)
}

func (c *HyperdriveConfig) ExporterContainerName() string {
	return string(config.ContainerID_Exporter)
}

func (c *HyperdriveConfig) GrafanaContainerName() string {
	return string(config.ContainerID_Grafana)
}

func (c *HyperdriveConfig) PrometheusContainerName() string {
	return string(config.ContainerID_Prometheus)
}

func (c *HyperdriveConfig) ExecutionClientDataVolume() string {
	return ExecutionClientDataVolume
}

func (c *HyperdriveConfig) BeaconNodeDataVolume() string {
	return BeaconNodeDataVolume
}

// ===============
// === General ===
// ===============

// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) IsLocalMode() bool {
	return cfg.ClientMode.Value == config.ClientMode_Local
}

// Gets the full name of the Docker container or volume with the provided suffix (name minus the project ID prefix)
func (cfg *HyperdriveConfig) GetDockerArtifactName(entity string) string {
	return fmt.Sprintf("%s_%s", cfg.ProjectName.Value, entity)
}

// Gets the name of the Execution Client start script
func (cfg *HyperdriveConfig) GetEcStartScript() string {
	return EcStartScript
}

// Gets the name of the Beacon Node start script
func (cfg *HyperdriveConfig) GetBnStartScript() string {
	return BnStartScript
}

// Gets the name of the Validator Client start script
func (cfg *HyperdriveConfig) GetVcStartScript() string {
	return VcStartScript
}

func (cfg *HyperdriveConfig) BnHttpUrl() (string, error) {
	/*
		// Check if Rescue Node is in-use
		cc, _ := cfg.GetSelectedConsensusClient()
		overrides, err := cfg.RescueNode.(*rescue_node.RescueNode).GetOverrides(cc)
		if err != nil {
			return "", fmt.Errorf("error using Rescue Node: %w", err)
		}
		if overrides != nil {
			// Use the rescue node
			return overrides.CcApiEndpoint, nil
		}
	*/
	if cfg.IsLocalMode() {
		return fmt.Sprintf("http://%s:%d", config.ContainerID_BeaconNode, cfg.LocalBeaconConfig.HttpPort.Value), nil
	}
	return cfg.ExternalBeaconConfig.HttpUrl.Value, nil
}

func (cfg *HyperdriveConfig) BnRpcUrl() (string, error) {
	/*
		// Check if Rescue Node is in-use
		cc, _ := cfg.GetSelectedConsensusClient()
		if cc != config.ConsensusClient_Prysm {
			return "", nil
		}

		overrides, err := cfg.RescueNode.(*rescue_node.RescueNode).GetOverrides(cc)
		if err != nil {
			return "", fmt.Errorf("error using Rescue Node: %w", err)
		}
		if overrides != nil {
			// Use the rescue node
			return overrides.CcRpcEndpoint, nil
		}
	*/
	if cfg.IsLocalMode() {
		return fmt.Sprintf("%s:%d", config.ContainerID_BeaconNode, cfg.LocalBeaconConfig.Prysm.RpcPort.Value), nil
	}
	return cfg.ExternalBeaconConfig.PrysmRpcUrl.Value, nil
}

func (cfg *HyperdriveConfig) FallbackBnHttpUrl() string {
	if !cfg.Fallback.UseFallbackClients.Value {
		return ""
	}
	return cfg.Fallback.BnHttpUrl.Value
}

func (cfg *HyperdriveConfig) FallbackBnRpcUrl() string {
	if !cfg.Fallback.UseFallbackClients.Value {
		return ""
	}
	return cfg.Fallback.PrysmRpcUrl.Value
}

func (cfg *HyperdriveConfig) AutoTxMaxFeeInt() uint64 {
	return uint64(cfg.AutoTxMaxFee.Value)
}

func (cfg *HyperdriveConfig) AutoTxGasThresholdInt() uint64 {
	return uint64(cfg.AutoTxGasThreshold.Value)
}

// ==============
// === Daemon ===
// ==============

func (cfg *HyperdriveConfig) GetDaemonContainerTag() string {
	return hyperdriveTag
}

// ========================
// === Execution Client ===
// ========================

// Get the selected Beacon Node
func (cfg *HyperdriveConfig) GetSelectedExecutionClient() config.ExecutionClient {
	if cfg.IsLocalMode() {
		return cfg.LocalExecutionConfig.ExecutionClient.Value
	}
	return cfg.ExternalExecutionConfig.ExecutionClient.Value
}

// Gets the port mapping of the ec container
// Used by text/template to format ec.yml
func (cfg *HyperdriveConfig) GetEcOpenApiPorts() string {
	return cfg.LocalExecutionConfig.GetOpenApiPortMapping()
}

// Gets the max peers of the ec container
// Used by text/template to format ec.yml
func (cfg *HyperdriveConfig) GetEcMaxPeers() (uint16, error) {
	if cfg.ClientMode.Value != config.ClientMode_Local {
		return 0, fmt.Errorf("Execution client is external, there is no max peers")
	}
	return cfg.LocalExecutionConfig.GetMaxPeers(), nil
}

// Gets the tag of the ec container
// Used by text/template to format ec.yml
func (cfg *HyperdriveConfig) GetEcContainerTag() (string, error) {
	if cfg.ClientMode.Value != config.ClientMode_Local {
		return "", fmt.Errorf("Execution client is external, there is no container tag")
	}
	return cfg.LocalExecutionConfig.GetContainerTag(), nil
}

// Used by text/template to format ec.yml
func (cfg *HyperdriveConfig) GetEcAdditionalFlags() (string, error) {
	if cfg.ClientMode.Value != config.ClientMode_Local {
		return "", fmt.Errorf("Execution client is external, there are no additional flags")
	}
	return cfg.LocalExecutionConfig.GetAdditionalFlags(), nil
}

// Used by text/template to format ec.yml
func (cfg *HyperdriveConfig) GetExternalIP() string {
	// Get the external IP address
	ip, err := config.GetExternalIP()
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

// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetEcHttpEndpoint() string {
	if cfg.ClientMode.Value == config.ClientMode_Local {
		return fmt.Sprintf("http://%s:%d", config.ContainerID_ExecutionClient, cfg.LocalExecutionConfig.HttpPort.Value)
	}

	return cfg.ExternalExecutionConfig.HttpUrl.Value
}

// Get the endpoints of the EC, including the fallback if applicable
func (cfg *HyperdriveConfig) GetEcHttpEndpointsWithFallback() string {
	endpoints := cfg.GetEcHttpEndpoint()

	if cfg.Fallback.UseFallbackClients.Value {
		endpoints = fmt.Sprintf("%s,%s", endpoints, cfg.Fallback.EcHttpUrl.Value)
	}
	return endpoints
}

// ===================
// === Beacon Node ===
// ===================

// Get the selected Beacon Node
func (cfg *HyperdriveConfig) GetSelectedBeaconNode() config.BeaconNode {
	if cfg.IsLocalMode() {
		return cfg.LocalBeaconConfig.BeaconNode.Value
	}
	return cfg.ExternalBeaconConfig.BeaconNode.Value
}

// Gets the tag of the bn container
// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetBnContainerTag() (string, error) {
	if cfg.ClientMode.Value != config.ClientMode_Local {
		return "", fmt.Errorf("Beacon Node is external, there is no container tag")
	}
	return cfg.LocalBeaconConfig.GetContainerTag(), nil
}

// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetBnOpenPorts() []string {
	return cfg.LocalBeaconConfig.GetOpenApiPortMapping()
}

// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetEcWsEndpoint() string {
	if cfg.ClientMode.Value == config.ClientMode_Local {
		return fmt.Sprintf("ws://%s:%d", config.ContainerID_ExecutionClient, cfg.LocalExecutionConfig.WebsocketPort.Value)
	}

	return cfg.ExternalExecutionConfig.WebsocketUrl.Value
}

// Gets the max peers of the bn container
// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetBnMaxPeers() (uint16, error) {
	if cfg.ClientMode.Value != config.ClientMode_Local {
		return 0, fmt.Errorf("Beacon Node is external, there is no max peers")
	}
	return cfg.LocalBeaconConfig.GetMaxPeers(), nil
}

// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetBnAdditionalFlags() (string, error) {
	if cfg.ClientMode.Value != config.ClientMode_Local {
		return "", fmt.Errorf("Beacon Node is external, there is no additional flags")
	}
	return cfg.LocalBeaconConfig.GetAdditionalFlags(), nil
}

// Get the HTTP API endpoint for the provided BN
func (cfg *HyperdriveConfig) GetBnHttpEndpoint() string {
	if cfg.ClientMode.Value == config.ClientMode_Local {
		return fmt.Sprintf("http://%s:%d", config.ContainerID_BeaconNode, cfg.LocalBeaconConfig.HttpPort.Value)
	}

	return cfg.ExternalBeaconConfig.HttpUrl.Value
}

// Get the endpoints of the BN, including the fallback if applicable
func (cfg *HyperdriveConfig) GetBnHttpEndpointsWithFallback() string {
	endpoints := cfg.GetBnHttpEndpoint()

	if cfg.Fallback.UseFallbackClients.Value {
		endpoints = fmt.Sprintf("%s,%s", endpoints, cfg.Fallback.BnHttpUrl.Value)
	}
	return endpoints
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
	if cfg.ClientMode.Value == config.ClientMode_Local {
		return string(config.ContainerID_ExecutionClient), nil
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
	if cfg.ClientMode.Value == config.ClientMode_Local {
		return string(config.ContainerID_BeaconNode), nil
	}
	ccUrl, err := url.Parse(cfg.ExternalBeaconConfig.HttpUrl.Value)
	if err != nil {
		return "", fmt.Errorf("Invalid External Consensus URL %s: %w", cfg.ExternalBeaconConfig.HttpUrl.Value, err)
	}

	return ccUrl.Hostname(), nil
}

// Used by text/template to format validator.yml
// Only returns the the prefix
func (cfg *HyperdriveConfig) GraffitiPrefix() string {
	identifier := ""
	versionString := fmt.Sprintf("v%s", shared.HyperdriveVersion)
	if len(versionString) < 8 {
		var ec config.ExecutionClient
		var bn config.BeaconNode
		if cfg.IsLocalMode() {
			ec = cfg.LocalExecutionConfig.ExecutionClient.Value
			bn = cfg.LocalBeaconConfig.BeaconNode.Value
		} else {
			ec = cfg.ExternalExecutionConfig.ExecutionClient.Value
			bn = cfg.ExternalBeaconConfig.BeaconNode.Value
		}

		ecInitial := strings.ToUpper(string(ec)[:1])

		var ccInitial string
		switch bn {
		case config.BeaconNode_Lodestar:
			ccInitial = "S" // Lodestar is special because it conflicts with Lighthouse
		default:
			ccInitial = strings.ToUpper(string(bn)[:1])
		}
		identifier = fmt.Sprintf("-%s%s", ecInitial, ccInitial)
	}

	return fmt.Sprintf("HD%s %s", identifier, versionString)
}
