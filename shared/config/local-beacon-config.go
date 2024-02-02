package config

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Param IDs
	BnCheckpointSyncUrlID string = "checkpointSyncUrl"
)

// Common parameters shared by all of the Beacon Clients
type LocalBeaconConfig struct {
	// The selected BN
	BeaconNode types.Parameter[types.BeaconNode]

	// The checkpoint sync URL if used
	CheckpointSyncProvider types.Parameter[string]

	// The port to use for gossip traffic
	P2pPort types.Parameter[uint16]

	// The port to expose the HTTP API on
	HttpPort types.Parameter[uint16]

	// Toggle for forwarding the HTTP API port outside of Docker
	OpenHttpPort types.Parameter[types.RpcPortMode]

	// Subconfigs
	Lighthouse *LighthouseBnConfig
	Lodestar   *LodestarBnConfig
	Nimbus     *NimbusBnConfig
	Prysm      *PrysmBnConfig
	Teku       *TekuBnConfig

	// Internal Fields
	parent *HyperdriveConfig
}

// Create a new LocalBeaconConfig struct
func NewLocalBeaconConfig(parent *HyperdriveConfig) *LocalBeaconConfig {
	cfg := &LocalBeaconConfig{
		parent: parent,

		BeaconNode: types.Parameter[types.BeaconNode]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.BnID,
				Name:               "Beacon Node",
				Description:        "Select which Beacon Node client you would like to use.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon, types.ContainerID_BeaconNode, types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: []*types.ParameterOption[types.BeaconNode]{
				{
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Lighthouse",
						Description: "Lighthouse is a Beacon Node with a heavy focus on speed and security. The team behind it, Sigma Prime, is an information security and software engineering firm who have funded Lighthouse along with the Ethereum Foundation, Consensys, and private individuals. Lighthouse is built in Rust and offered under an Apache 2.0 License.",
					},
					Value: types.BeaconNode_Lighthouse,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Lodestar",
						Description: "Lodestar is the fifth open-source Ethereum Beacon Node. It is written in Typescript maintained by ChainSafe Systems. Lodestar, their flagship product, is a production-capable Beacon Chain and Validator Client uniquely situated as the go-to for researchers and developers for rapid prototyping and browser usage.",
					},
					Value: types.BeaconNode_Lodestar,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Nimbus",
						Description: "Nimbus is a Beacon Node implementation that strives to be as lightweight as possible in terms of resources used. This allows it to perform well on embedded systems, resource-restricted devices -- including Raspberry Pis and mobile devices -- and multi-purpose servers.",
					},
					Value: types.BeaconNode_Nimbus,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Prysm",
						Description: "Prysm is a Go implementation of Ethereum Consensus protocol with a focus on usability, security, and reliability. Prysm is developed by Prysmatic Labs, a company with the sole focus on the development of their client. Prysm is written in Go and released under a GPL-3.0 license.",
					},
					Value: types.BeaconNode_Prysm,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Teku",
						Description: "PegaSys Teku (formerly known as Artemis) is a Java-based Ethereum 2.0 client designed & built to meet institutional needs and security requirements. PegaSys is an arm of ConsenSys dedicated to building enterprise-ready clients and tools for interacting with the core Ethereum platform. Teku is Apache 2 licensed and written in Java, a language notable for its maturity & ubiquity.",
					},
					Value: types.BeaconNode_Teku,
				}},
			Default: map[types.Network]types.BeaconNode{
				types.Network_All: types.BeaconNode_Nimbus,
			},
		},

		CheckpointSyncProvider: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:   BnCheckpointSyncUrlID,
				Name: "Checkpoint Sync URL",
				Description: "If you would like to instantly sync using an existing Beacon node, enter its URL.\n" +
					"Example: https://<project ID>:<secret>@eth2-beacon-prater.infura.io\n" +
					"Leave this blank if you want to sync normally from the start of the chain.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},

		P2pPort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.P2pPortID,
				Name:               "P2P Port",
				Description:        "The port to use for P2P (blockchain) traffic.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 9001,
			},
		},

		HttpPort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.HttpPortID,
				Name:               "HTTP API Port",
				Description:        "The port your Beacon Node should run its HTTP API on.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon, types.ContainerID_BeaconNode, types.ContainerID_ValidatorClients, types.ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 5052,
			},
		},

		OpenHttpPort: types.Parameter[types.RpcPortMode]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.OpenHttpPortsID,
				Name:               "Expose API Port",
				Description:        "Select an option to expose your Beacon Node's API port to your localhost or external hosts on the network, so other machines can access it too.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: getPortModes("Allow connections from external hosts. This is safe if you're running your node on your local network. If you're a VPS user, this would expose your node to the internet and could make it vulnerable to MEV/tips theft"),
			Default: map[types.Network]types.RpcPortMode{
				types.Network_All: types.RpcPortMode_Closed,
			},
		},
	}

	cfg.Lighthouse = NewLighthouseBnConfig(cfg)
	cfg.Lodestar = NewLodestarBnConfig(cfg)
	cfg.Nimbus = NewNimbusBnConfig(cfg)
	cfg.Prysm = NewPrysmBnConfig(cfg)
	cfg.Teku = NewTekuBnConfig(cfg)

	return cfg
}

// The title for the config
func (cfg *LocalBeaconConfig) GetTitle() string {
	return "Local Beacon Node Settings"
}

// Get the parameters for this config
func (cfg *LocalBeaconConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.BeaconNode,
		&cfg.CheckpointSyncProvider,
		&cfg.P2pPort,
		&cfg.HttpPort,
		&cfg.OpenHttpPort,
	}
}

// Get the sections underneath this one
func (cfg *LocalBeaconConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{
		"lighthouse": cfg.Lighthouse,
		"lodestar":   cfg.Lodestar,
		"nimbus":     cfg.Nimbus,
		"prysm":      cfg.Prysm,
		"teku":       cfg.Teku,
	}
}

// ==================
// === Templating ===
// ==================

// Get the Docker mapping for the selected API port mode
func (cfg *LocalBeaconConfig) getOpenApiPortMapping() []string {
	bnOpenPorts := make([]string, 0)

	// Handle the standard HTTP API port
	apiPortMode := cfg.OpenHttpPort.Value
	if apiPortMode.IsOpen() {
		apiPort := cfg.HttpPort.Value
		bnOpenPorts = append(bnOpenPorts, apiPortMode.DockerPortMapping(apiPort))
	}

	// Handle Prysm's RPC port
	if cfg.BeaconNode.Value == types.BeaconNode_Prysm {
		prysmRpcPortMode := cfg.Prysm.OpenRpcPort.Value
		if prysmRpcPortMode.IsOpen() {
			prysmRpcPort := cfg.Prysm.RpcPort.Value
			bnOpenPorts = append(bnOpenPorts, prysmRpcPortMode.DockerPortMapping(prysmRpcPort))
		}
	}
	return bnOpenPorts
}

// Gets the max peers of the selected EC
func (cfg *LocalBeaconConfig) getMaxPeers() uint16 {
	switch cfg.BeaconNode.Value {
	case types.BeaconNode_Lighthouse:
		return cfg.Lighthouse.MaxPeers.Value
	case types.BeaconNode_Lodestar:
		return cfg.Lodestar.MaxPeers.Value
	case types.BeaconNode_Nimbus:
		return cfg.Nimbus.MaxPeers.Value
	case types.BeaconNode_Prysm:
		return cfg.Prysm.MaxPeers.Value
	case types.BeaconNode_Teku:
		return cfg.Teku.MaxPeers.Value
	default:
		panic(fmt.Sprintf("Unknown Beacon Node %s", string(cfg.BeaconNode.Value)))
	}
}

// Get the container tag of the selected BN
func (cfg *LocalBeaconConfig) getContainerTag() string {
	switch cfg.BeaconNode.Value {
	case types.BeaconNode_Lighthouse:
		return cfg.Lighthouse.ContainerTag.Value
	case types.BeaconNode_Lodestar:
		return cfg.Lodestar.ContainerTag.Value
	case types.BeaconNode_Nimbus:
		return cfg.Nimbus.ContainerTag.Value
	case types.BeaconNode_Prysm:
		return cfg.Prysm.ContainerTag.Value
	case types.BeaconNode_Teku:
		return cfg.Teku.ContainerTag.Value
	default:
		panic(fmt.Sprintf("Unknown Beacon Node %s", string(cfg.BeaconNode.Value)))
	}
}

// Gets the additional flags of the selected BN
func (cfg *LocalBeaconConfig) getAdditionalFlags() string {
	switch cfg.BeaconNode.Value {
	case types.BeaconNode_Lighthouse:
		return cfg.Lighthouse.AdditionalFlags.Value
	case types.BeaconNode_Lodestar:
		return cfg.Lodestar.AdditionalFlags.Value
	case types.BeaconNode_Nimbus:
		return cfg.Nimbus.AdditionalFlags.Value
	case types.BeaconNode_Prysm:
		return cfg.Prysm.AdditionalFlags.Value
	case types.BeaconNode_Teku:
		return cfg.Teku.AdditionalFlags.Value
	default:
		panic(fmt.Sprintf("Unknown Beacon Node %s", string(cfg.BeaconNode.Value)))
	}
}
