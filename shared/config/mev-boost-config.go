package config

import (
	"fmt"
	"strings"

	ids "github.com/nodeset-org/hyperdrive-daemon/shared/config/ids"
	"github.com/rocket-pool/node-manager-core/config"
	nmc_ids "github.com/rocket-pool/node-manager-core/config/ids"
)

// Constants
const (
	mevBoostTag string = "flashbots/mev-boost:1.7"
)

// A MEV relay
type MevRelay struct {
	ID          MevRelayID
	Name        string
	Description string
	Urls        map[string]string
}

// Configuration for MEV-Boost
type MevBoostConfig struct {
	// Toggle to enable / disable
	Enable config.Parameter[bool]

	// Ownership mode
	Mode config.Parameter[config.ClientMode]

	// The mode for relay selection
	SelectionMode config.Parameter[MevSelectionMode]

	// Flashbots relay
	FlashbotsRelay config.Parameter[bool]

	// bloXroute max profit relay
	BloxRouteMaxProfitRelay config.Parameter[bool]

	// bloXroute regulated relay
	BloxRouteRegulatedRelay config.Parameter[bool]

	// Titan regional relay
	TitanRegionalRelay config.Parameter[bool]

	// Custom relays provided by the user
	CustomRelays config.Parameter[string]

	// The RPC port
	Port config.Parameter[uint16]

	// Toggle for forwarding the HTTP port outside of Docker
	OpenRpcPort config.Parameter[config.RpcPortMode]

	// The Docker Hub tag for MEV-Boost
	ContainerTag config.Parameter[string]

	// Custom command line flags
	AdditionalFlags config.Parameter[string]

	// The URL of an external MEV-Boost client
	ExternalUrl config.Parameter[string]

	///////////////////////////
	// Non-editable settings //
	///////////////////////////

	parent   *HyperdriveConfig
	relays   []MevRelay
	relayMap map[MevRelayID]MevRelay
}

// Generates a new MEV-Boost configuration
func NewMevBoostConfig(parent *HyperdriveConfig) *MevBoostConfig {
	// Generate the relays
	relays := createDefaultRelays()
	relayMap := map[MevRelayID]MevRelay{}
	for _, relay := range relays {
		relayMap[relay.ID] = relay
	}

	rpcPortModes := config.GetPortModes("")

	return &MevBoostConfig{
		parent: parent,

		Enable: config.Parameter[bool]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.MevBoostEnableID,
				Name:               "Enable MEV-Boost",
				Description:        "Enable MEV-Boost, which connects your validators to one or more relays of your choice. The relays act as intermediaries between you and professional block builders that find and extract MEV opportunities. The builders will give you a healthy tip in return, which tends to be worth more than blocks you built on your own.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_BeaconNode, config.ContainerID_MevBoost},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]bool{
				config.Network_All: true,
			},
		},

		Mode: config.Parameter[config.ClientMode]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.MevBoostModeID,
				Name:               "MEV-Boost Mode",
				Description:        "Choose whether to let Hyperdrive manage your MEV-Boost instance (Locally Managed), or if you manage your own outside of Hyperdrive (Externally Managed).",
				AffectsContainers:  []config.ContainerID{config.ContainerID_BeaconNode, config.ContainerID_MevBoost},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: []*config.ParameterOption[config.ClientMode]{{
				ParameterOptionCommon: &config.ParameterOptionCommon{
					Name:        "Locally Managed",
					Description: "Allow Hyperdrive to manage the MEV-Boost client for you",
				},
				Value: config.ClientMode_Local,
			}, {
				ParameterOptionCommon: &config.ParameterOptionCommon{
					Name:        "Externally Managed",
					Description: "Use an existing MEV-Boost client that you manage on your own",
				},
				Value: config.ClientMode_External,
			}},
			Default: map[config.Network]config.ClientMode{
				config.Network_All: config.ClientMode_Local,
			},
		},

		SelectionMode: config.Parameter[MevSelectionMode]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.MevBoostSelectionModeID,
				Name:               "Selection Mode",
				Description:        "Select how the TUI shows you the options for which MEV relays to enable.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_MevBoost},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: []*config.ParameterOption[MevSelectionMode]{{
				ParameterOptionCommon: &config.ParameterOptionCommon{
					Name:        "Use All Relays",
					Description: "Use this if you simply want to enable all of the built-in relays without needing to read about each individual relay. If new relays get added to Hyperdrive, you'll automatically start using those too.\n\nNote that all of Hyperdrive's built-in relays support regional sanction lists (such as the US OFAC list) and are compliant with regulations. To learn more, please visit https://medium.com/coinmonks/understanding-the-impact-of-the-ofac-sanctions-on-block-builders-9c0e02b7e450.",
				},
				Value: MevSelectionMode_All,
			}, {
				ParameterOptionCommon: &config.ParameterOptionCommon{
					Name:        "Manual Mode",
					Description: "Each relay will be shown, and you can enable each one individually as you see fit.\nUse this if you already know about the relays and want to customize the ones you will use.\n\nNote that all of Hyperdrive's built-in relays support regional sanction lists (such as the US OFAC list) and are compliant with regulations. To learn more, please visit https://medium.com/coinmonks/understanding-the-impact-of-the-ofac-sanctions-on-block-builders-9c0e02b7e450.",
				},
				Value: MevSelectionMode_Manual,
			}},
			Default: map[config.Network]MevSelectionMode{
				config.Network_All: MevSelectionMode_All,
			},
		},

		// Explicit relay params
		FlashbotsRelay:          generateRelayParameter(ids.MevBoostFlashbotsID, relayMap[MevRelayID_Flashbots]),
		BloxRouteMaxProfitRelay: generateRelayParameter(ids.MevBoostBloxRouteMaxProfitID, relayMap[MevRelayID_BloxrouteMaxProfit]),
		BloxRouteRegulatedRelay: generateRelayParameter(ids.MevBoostBloxRouteRegulatedID, relayMap[MevRelayID_BloxrouteRegulated]),
		TitanRegionalRelay:      generateRelayParameter(ids.MevBoostTitanRegionalID, relayMap[MevRelayID_TitanRegional]),

		CustomRelays: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:          ids.MevBoostCustomRelaysID,
				Name:        "Custom Relays",
				Description: "Add custom relay URLs to MEV-Boost that aren't part of the built-in set. You can add multiple relays by separating each one with a comma. Any relay URLs can be used as long as they match your selected Ethereum network.\n\nFor a comprehensive list of available relays, we recommend the list maintained by ETHStaker:\nhttps://github.com/eth-educators/ethstaker-guides/blob/main/MEV-relay-list.md",
			},
			Default: map[config.Network]string{
				config.Network_All: "",
			},
		},

		Port: config.Parameter[uint16]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 nmc_ids.PortID,
				Name:               "Port",
				Description:        "The port that MEV-Boost should serve its API on.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_BeaconNode, config.ContainerID_MevBoost},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]uint16{
				config.Network_All: uint16(18550),
			},
		},

		OpenRpcPort: config.Parameter[config.RpcPortMode]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 nmc_ids.OpenPortID,
				Name:               "Expose API Port",
				Description:        "Expose the API port to other processes on your machine, or to your local network so other local machines can access MEV-Boost's API.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_MevBoost},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: rpcPortModes,
			Default: map[config.Network]config.RpcPortMode{
				config.Network_All: config.RpcPortMode_Closed,
			},
		},

		ContainerTag: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 nmc_ids.ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the MEV-Boost container you want to use on Docker Hub.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_MevBoost},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[config.Network]string{
				config.Network_All: mevBoostTag,
			},
		},

		AdditionalFlags: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 nmc_ids.AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass to MEV-Boost, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_MevBoost},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]string{
				config.Network_All: "",
			},
		},

		ExternalUrl: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.MevBoostExternalUrlID,
				Name:               "External URL",
				Description:        "The URL of the external MEV-Boost client or provider",
				AffectsContainers:  []config.ContainerID{config.ContainerID_BeaconNode},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]string{
				config.Network_All: "",
			},
		},

		relays:   relays,
		relayMap: relayMap,
	}
}

// The title for the config
func (cfg *MevBoostConfig) GetTitle() string {
	return "MEV-Boost"
}

// Get the Parameters for this config
func (cfg *MevBoostConfig) GetParameters() []config.IParameter {
	return []config.IParameter{
		&cfg.Enable,
		&cfg.Mode,
		&cfg.SelectionMode,
		&cfg.FlashbotsRelay,
		&cfg.BloxRouteMaxProfitRelay,
		&cfg.BloxRouteRegulatedRelay,
		&cfg.TitanRegionalRelay,
		&cfg.CustomRelays,
		&cfg.Port,
		&cfg.OpenRpcPort,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
		&cfg.ExternalUrl,
	}
}

// Get the sections underneath this one
func (cfg *MevBoostConfig) GetSubconfigs() map[string]config.IConfigSection {
	return map[string]config.IConfigSection{}
}

// Checks if any relays are available for the current network
func (cfg *MevBoostConfig) HasRelays() bool {
	networkName := cfg.parent.GetEthNetworkName()
	if networkName == "" {
		return false
	}

	// Check if any of the relays are available for that Eth network
	for _, relay := range cfg.relays {
		_, exists := relay.Urls[networkName]
		if !exists {
			continue
		}
		return true
	}

	return false
}

// Get the relays that are available for the current network
func (cfg *MevBoostConfig) GetAvailableRelays() []MevRelay {
	relays := []MevRelay{}
	networkName := cfg.parent.GetEthNetworkName()
	if networkName == "" {
		return relays
	}

	for _, relay := range cfg.relays {
		_, exists := relay.Urls[networkName]
		if !exists {
			continue
		}
		relays = append(relays, relay)
	}
	return relays
}

// Get which MEV-boost relays are enabled
func (cfg *MevBoostConfig) GetEnabledMevRelays() []MevRelay {
	relays := []MevRelay{}
	networkName := cfg.parent.GetEthNetworkName()
	if networkName == "" {
		return relays
	}

	switch cfg.SelectionMode.Value {
	case MevSelectionMode_All:
		for _, relay := range cfg.relays {
			_, exists := relay.Urls[networkName]
			if !exists {
				// Skip relays that don't exist on the current network
				continue
			}
			relays = append(relays, relay)
		}

	case MevSelectionMode_Manual:
		if cfg.FlashbotsRelay.Value {
			_, exists := cfg.relayMap[MevRelayID_Flashbots].Urls[networkName]
			if exists {
				relays = append(relays, cfg.relayMap[MevRelayID_Flashbots])
			}
		}
		if cfg.BloxRouteMaxProfitRelay.Value {
			_, exists := cfg.relayMap[MevRelayID_BloxrouteMaxProfit].Urls[networkName]
			if exists {
				relays = append(relays, cfg.relayMap[MevRelayID_BloxrouteMaxProfit])
			}
		}
		if cfg.BloxRouteRegulatedRelay.Value {
			_, exists := cfg.relayMap[MevRelayID_BloxrouteRegulated].Urls[networkName]
			if exists {
				relays = append(relays, cfg.relayMap[MevRelayID_BloxrouteRegulated])
			}
		}
		if cfg.TitanRegionalRelay.Value {
			_, exists := cfg.relayMap[MevRelayID_TitanRegional].Urls[networkName]
			if exists {
				relays = append(relays, cfg.relayMap[MevRelayID_TitanRegional])
			}
		}
	}

	return relays
}

func (cfg *MevBoostConfig) GetRelayString() string {
	relayUrls := []string{}
	networkName := cfg.parent.GetEthNetworkName()
	if networkName == "" {
		return ""
	}

	relays := cfg.GetEnabledMevRelays()
	for _, relay := range relays {
		relayUrls = append(relayUrls, relay.Urls[networkName])
	}
	if cfg.CustomRelays.Value != "" {
		relayUrls = append(relayUrls, cfg.CustomRelays.Value)
	}

	relayString := strings.Join(relayUrls, ",")
	return relayString
}

// Create the default MEV relays
func createDefaultRelays() []MevRelay {
	relays := []MevRelay{
		// Flashbots
		{
			ID:          MevRelayID_Flashbots,
			Name:        "Flashbots",
			Description: "Flashbots is the developer of MEV-Boost, and one of the best-known and most trusted relays in the space.",
			Urls: map[string]string{
				config.EthNetwork_Mainnet: "https://0xac6e77dfe25ecd6110b8e780608cce0dab71fdd5ebea22a16c0205200f2f8e2e3ad3b71d3499c54ad14d6c21b41a37ae@boost-relay.flashbots.net?id=hyperdrive",
				config.EthNetwork_Holesky: "https://0xafa4c6985aa049fb79dd37010438cfebeb0f2bd42b115b89dd678dab0670c1de38da0c4e9138c9290a398ecd9a0b3110@boost-relay-holesky.flashbots.net?id=hyperdrive",
			},
		},

		// bloXroute Max Profit
		{
			ID:          MevRelayID_BloxrouteMaxProfit,
			Name:        "bloXroute Max Profit",
			Description: "Select this to enable the \"max profit\" relay from bloXroute.",
			Urls: map[string]string{
				config.EthNetwork_Mainnet: "https://0x8b5d2e73e2a3a55c6c87b8b6eb92e0149a125c852751db1422fa951e42a09b82c142c3ea98d0d9930b056a3bc9896b8f@bloxroute.max-profit.blxrbdn.com?id=hyperdrive",
				config.EthNetwork_Holesky: "https://0x821f2a65afb70e7f2e820a925a9b4c80a159620582c1766b1b09729fec178b11ea22abb3a51f07b288be815a1a2ff516@bloxroute.holesky.blxrbdn.com",
			},
		},

		// bloXroute Regulated
		{
			ID:          MevRelayID_BloxrouteRegulated,
			Name:        "bloXroute Regulated",
			Description: "Select this to enable the \"regulated\" relay from bloXroute.",
			Urls: map[string]string{
				config.EthNetwork_Mainnet: "https://0xb0b07cd0abef743db4260b0ed50619cf6ad4d82064cb4fbec9d3ec530f7c5e6793d9f286c4e082c0244ffb9f2658fe88@bloxroute.regulated.blxrbdn.com?id=hyperdrive",
			},
		},

		// Titan Regional
		{
			ID:          MevRelayID_TitanRegional,
			Name:        "Titan Regional",
			Description: "Titan Relay is a neutral, Rust-based MEV-Boost Relay optimized for low latency through put, geographical distribution, and robustness. This is the regulated (censoring) version.",
			Urls: map[string]string{
				config.EthNetwork_Mainnet: "https://0x8c4ed5e24fe5c6ae21018437bde147693f68cda427cd1122cf20819c30eda7ed74f72dece09bb313f2a1855595ab677d@regional.titanrelay.xyz",
			},
		},
	}

	return relays
}

// Generate one of the relay parameters
func generateRelayParameter(id string, relay MevRelay) config.Parameter[bool] {
	description := fmt.Sprintf("[lime]NOTE: You can enable multiple options.\n\n[white]%s\n\n", relay.Description)

	return config.Parameter[bool]{
		ParameterCommon: &config.ParameterCommon{
			ID:                 id,
			Name:               fmt.Sprintf("Enable %s", relay.Name),
			Description:        description,
			AffectsContainers:  []config.ContainerID{config.ContainerID_MevBoost},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},
		Default: map[config.Network]bool{
			config.Network_All: false,
		},
	}
}
