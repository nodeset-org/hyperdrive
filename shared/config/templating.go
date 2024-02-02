package config

import (
	"fmt"
	"net/url"
	"strings"

	swshared "github.com/nodeset-org/hyperdrive/modules/stakewise/shared"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

// ===============
// === General ===
// ===============

// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) IsLocalMode() bool {
	return cfg.ClientMode.Value == types.ClientMode_Local
}

// Check if any of the services have doppelganger detection enabled
// NOTE: update this with each new service that runs a VC!
func (cfg *HyperdriveConfig) IsDoppelgangerEnabled() bool {
	return cfg.Modules.Stakewise.VcCommon.DoppelgangerDetection.Value
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
		return fmt.Sprintf("http://%s:%d", types.ContainerID_BeaconNode, cfg.LocalBeaconConfig.HttpPort.Value), nil
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
		return fmt.Sprintf("%s:%d", types.ContainerID_BeaconNode, cfg.LocalBeaconConfig.Prysm.RpcPort.Value), nil
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

// ==============
// === Daemon ===
// ==============

func (cfg *HyperdriveConfig) GetDaemonContainerTag() string {
	return hyperdriveTag
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

// Used by text/template to format bn.yml
func (cfg *HyperdriveConfig) GetEcHttpEndpoint() string {
	if cfg.ClientMode.Value == types.ClientMode_Local {
		return fmt.Sprintf("http://%s:%d", types.ContainerID_ExecutionClient, cfg.LocalExecutionConfig.HttpPort.Value)
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
func (cfg *HyperdriveConfig) GetEcWsEndpoint() string {
	if cfg.ClientMode.Value == types.ClientMode_Local {
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

// Get the HTTP API endpoint for the provided BN
func (cfg *HyperdriveConfig) GetBnHttpEndpoint() string {
	if cfg.ClientMode.Value == types.ClientMode_Local {
		return fmt.Sprintf("http://%s:%d", types.ContainerID_BeaconNode, cfg.LocalBeaconConfig.HttpPort.Value)
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

// =================
// === Stakewise ===
// =================
// TODO: Find a better way to break this out so it isn't in the root config

// Get the container tag of the selected VC
func (cfg *HyperdriveConfig) GetStakewiseVcContainerTag() string {
	var bn types.BeaconNode
	if cfg.IsLocalMode() {
		bn = cfg.LocalBeaconConfig.BeaconNode.Value
	} else {
		bn = cfg.ExternalBeaconConfig.BeaconNode.Value
	}

	switch bn {
	case types.BeaconNode_Lighthouse:
		return cfg.Modules.Stakewise.Lighthouse.ContainerTag.Value
	case types.BeaconNode_Lodestar:
		return cfg.Modules.Stakewise.Lodestar.ContainerTag.Value
	case types.BeaconNode_Nimbus:
		return cfg.Modules.Stakewise.Nimbus.ContainerTag.Value
	case types.BeaconNode_Prysm:
		return cfg.Modules.Stakewise.Prysm.ContainerTag.Value
	case types.BeaconNode_Teku:
		return cfg.Modules.Stakewise.Teku.ContainerTag.Value
	default:
		panic(fmt.Sprintf("Unknown Beacon Node %s", string(cfg.LocalBeaconConfig.BeaconNode.Value)))
	}
}

// Gets the additional flags of the selected VC
func (cfg *HyperdriveConfig) GetStakewiseVcAdditionalFlags() string {
	var bn types.BeaconNode
	if cfg.IsLocalMode() {
		bn = cfg.LocalBeaconConfig.BeaconNode.Value
	} else {
		bn = cfg.ExternalBeaconConfig.BeaconNode.Value
	}

	switch bn {
	case types.BeaconNode_Lighthouse:
		return cfg.Modules.Stakewise.Lighthouse.AdditionalFlags.Value
	case types.BeaconNode_Lodestar:
		return cfg.Modules.Stakewise.Lodestar.AdditionalFlags.Value
	case types.BeaconNode_Nimbus:
		return cfg.Modules.Stakewise.Nimbus.AdditionalFlags.Value
	case types.BeaconNode_Prysm:
		return cfg.Modules.Stakewise.Prysm.AdditionalFlags.Value
	case types.BeaconNode_Teku:
		return cfg.Modules.Stakewise.Teku.AdditionalFlags.Value
	default:
		panic(fmt.Sprintf("Unknown Beacon Node %s", string(cfg.LocalBeaconConfig.BeaconNode.Value)))
	}
}

// Used by text/template to format validator.yml
func (cfg *HyperdriveConfig) StakewiseGraffiti() (string, error) {
	prefix := cfg.graffitiPrefix()
	customGraffiti := cfg.Modules.Stakewise.VcCommon.Graffiti.Value
	if customGraffiti == "" {
		return prefix, nil
	}
	return fmt.Sprintf("%s (%s)", prefix, customGraffiti), nil
}

func (cfg *HyperdriveConfig) StakewiseFeeRecipient() string {
	res := swshared.NewStakewiseResources(cfg.Network.Value)
	return res.FeeRecipient.Hex()
}

func (cfg *HyperdriveConfig) StakewiseVault() string {
	res := swshared.NewStakewiseResources(cfg.Network.Value)
	return res.Vault.Hex()
}

func (cfg *HyperdriveConfig) StakewiseNetwork() string {
	res := swshared.NewStakewiseResources(cfg.Network.Value)
	return res.NodesetNetwork
}

// Used by text/template to format validator.yml
// Only returns the the prefix
func (cfg *HyperdriveConfig) graffitiPrefix() string {
	identifier := ""
	versionString := fmt.Sprintf("v%s", shared.HyperdriveVersion)
	if len(versionString) < 8 {
		var ec types.ExecutionClient
		var bn types.BeaconNode
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
		case types.BeaconNode_Lodestar:
			ccInitial = "S" // Lodestar is special because it conflicts with Lighthouse
		default:
			ccInitial = strings.ToUpper(string(bn)[:1])
		}
		identifier = fmt.Sprintf("-%s%s", ecInitial, ccInitial)
	}

	return fmt.Sprintf("HD%s %s", identifier, versionString)
}
