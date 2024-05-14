package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/nodeset-org/hyperdrive-daemon/shared"
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

// Container tag for the daemon
func (cfg *HyperdriveConfig) GetDaemonContainerTag() string {
	return cfg.ContainerTag.Value
}

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
	return cfg.GetBnHttpEndpoint(), nil
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
		return fmt.Sprintf("%s:%d", config.ContainerID_BeaconNode, cfg.LocalBeaconClient.Prysm.RpcPort.Value), nil
	}
	return cfg.ExternalBeaconClient.PrysmRpcUrl.Value, nil
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

// ========================
// === Execution Client ===
// ========================

// Get the selected Beacon Node
func (cfg *HyperdriveConfig) GetSelectedExecutionClient() config.ExecutionClient {
	if cfg.IsLocalMode() {
		return cfg.LocalExecutionClient.ExecutionClient.Value
	}
	return cfg.ExternalExecutionClient.ExecutionClient.Value
}

// Gets the port mapping of the ec container
// Used by text/template to format ec.yml
func (cfg *HyperdriveConfig) GetEcOpenApiPorts() string {
	return cfg.LocalExecutionClient.GetOpenApiPortMapping()
}

// Gets the max peers of the ec container
// Used by text/template to format ec.yml
func (cfg *HyperdriveConfig) GetEcMaxPeers() (uint16, error) {
	if !cfg.IsLocalMode() {
		return 0, fmt.Errorf("Execution client is external, there is no max peers")
	}
	return cfg.LocalExecutionClient.GetMaxPeers(), nil
}

// Gets the tag of the ec container
// Used by text/template to format ec.yml
func (cfg *HyperdriveConfig) GetEcContainerTag() (string, error) {
	if !cfg.IsLocalMode() {
		return "", fmt.Errorf("Execution client is external, there is no container tag")
	}
	return cfg.LocalExecutionClient.GetContainerTag(), nil
}

// Used by text/template to format ec.yml
func (cfg *HyperdriveConfig) GetEcAdditionalFlags() (string, error) {
	if !cfg.IsLocalMode() {
		return "", fmt.Errorf("Execution client is external, there are no additional flags")
	}
	return cfg.LocalExecutionClient.GetAdditionalFlags(), nil
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
		return fmt.Sprintf("http://%s:%d", config.ContainerID_ExecutionClient, cfg.LocalExecutionClient.HttpPort.Value)
	}

	return cfg.ExternalExecutionClient.HttpUrl.Value
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
		return cfg.LocalBeaconClient.BeaconNode.Value
	}
	return cfg.ExternalBeaconClient.BeaconNode.Value
}

// Gets the tag of the bn container
// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetBnContainerTag() (string, error) {
	if cfg.ClientMode.Value != config.ClientMode_Local {
		return "", fmt.Errorf("Beacon Node is external, there is no container tag")
	}
	return cfg.LocalBeaconClient.GetContainerTag(), nil
}

// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetBnOpenPorts() []string {
	return cfg.LocalBeaconClient.GetOpenApiPortMapping()
}

// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetEcWsEndpoint() string {
	if cfg.ClientMode.Value == config.ClientMode_Local {
		return fmt.Sprintf("ws://%s:%d", config.ContainerID_ExecutionClient, cfg.LocalExecutionClient.WebsocketPort.Value)
	}

	return cfg.ExternalExecutionClient.WebsocketUrl.Value
}

// Gets the max peers of the bn container
// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetBnMaxPeers() (uint16, error) {
	if cfg.ClientMode.Value != config.ClientMode_Local {
		return 0, fmt.Errorf("Beacon Node is external, there is no max peers")
	}
	return cfg.LocalBeaconClient.GetMaxPeers(), nil
}

// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetBnAdditionalFlags() (string, error) {
	if cfg.ClientMode.Value != config.ClientMode_Local {
		return "", fmt.Errorf("Beacon Node is external, there is no additional flags")
	}
	return cfg.LocalBeaconClient.GetAdditionalFlags(), nil
}

// Get the HTTP API endpoint for the provided BN
func (cfg *HyperdriveConfig) GetBnHttpEndpoint() string {
	if cfg.IsLocalMode() {
		return fmt.Sprintf("http://%s:%d", config.ContainerID_BeaconNode, cfg.LocalBeaconClient.HttpPort.Value)
	}

	return cfg.ExternalBeaconClient.HttpUrl.Value
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
	ecUrl, err := url.Parse(cfg.ExternalExecutionClient.HttpUrl.Value)
	if err != nil {
		return "", fmt.Errorf("Invalid External Execution URL %s: %w", cfg.ExternalExecutionClient.HttpUrl.Value, err)
	}

	return ecUrl.Hostname(), nil
}

// Gets the hostname portion of the Beacon Node's URI.
// Used by text/template to format prometheus.yml.
func (cfg *HyperdriveConfig) GetBeaconHostname() (string, error) {
	if cfg.ClientMode.Value == config.ClientMode_Local {
		return string(config.ContainerID_BeaconNode), nil
	}
	bnUrl, err := url.Parse(cfg.ExternalBeaconClient.HttpUrl.Value)
	if err != nil {
		return "", fmt.Errorf("Invalid External Consensus URL %s: %w", cfg.ExternalBeaconClient.HttpUrl.Value, err)
	}

	return bnUrl.Hostname(), nil
}

// Used by text/template to format validator.yml
// Only returns the the prefix
func (cfg *HyperdriveConfig) GraffitiPrefix() string {
	identifier := ""
	versionString := fmt.Sprintf("v%s", shared.HyperdriveVersion)
	if len(versionString) < 8 {
		ecInitial := strings.ToUpper(string(cfg.GetSelectedExecutionClient())[:1])

		var bnInitial string
		bn := cfg.GetSelectedBeaconNode()
		switch bn {
		case config.BeaconNode_Lodestar:
			bnInitial = "S" // Lodestar is special because it conflicts with Lighthouse
		default:
			bnInitial = strings.ToUpper(string(bn)[:1])
		}

		var modeFlag string
		if cfg.IsLocalMode() {
			modeFlag = "L"
		} else {
			modeFlag = "X"
		}
		identifier = fmt.Sprintf("%s%s%s", ecInitial, bnInitial, modeFlag)
	}

	return fmt.Sprintf("HD%s %s", identifier, versionString)
}
