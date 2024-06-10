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
	mevBoostTag                 string = "flashbots/mev-boost:1.7"
	RegulatedRelayDescription   string = "Select this to enable the relays that comply with government regulations (e.g. OFAC sanctions), "
	UnregulatedRelayDescription string = "Select this to enable the relays that do not follow any sanctions lists (do not censor transactions), "
	NoSandwichRelayDescription  string = "and do not allow front-running or sandwich attacks."
	AllMevRelayDescription      string = "and allow for all types of MEV (including sandwich attacks)."
)

// A MEV relay
type MevRelay struct {
	ID          MevRelayID
	Name        string
	Description string
	Urls        map[config.Network]string
	Regulated   bool
}

// Configuration for MEV-Boost
type MevBoostConfig struct {
	// Toggle to enable / disable
	Enable config.Parameter[bool]

	// Ownership mode
	Mode config.Parameter[config.ClientMode]

	// The mode for relay selection
	SelectionMode config.Parameter[MevSelectionMode]

	// Regulated, all types
	EnableRegulatedAllMev config.Parameter[bool]

	// Unregulated, all types
	EnableUnregulatedAllMev config.Parameter[bool]

	// Flashbots relay
	FlashbotsRelay config.Parameter[bool]

	// bloXroute max profit relay
	BloxRouteMaxProfitRelay config.Parameter[bool]

	// bloXroute regulated relay
	BloxRouteRegulatedRelay config.Parameter[bool]

	// Eden relay
	EdenRelay config.Parameter[bool]

	// Ultra sound relay
	UltrasoundRelay config.Parameter[bool]

	// Aestus relay
	AestusRelay config.Parameter[bool]

	// Titan global relay
	TitanGlobalRelay config.Parameter[bool]

	// Titan regional relay
	TitanRegionalRelay config.Parameter[bool]

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
					Name:        "Profile Mode",
					Description: "Relays will be bundled up based on whether or not they're regulated, and whether or not they allow sandwich attacks.\nUse this if you simply want to specify which type of relay you want to use without needing to read about each individual relay.",
				},
				Value: MevSelectionMode_Profile,
			}, {
				ParameterOptionCommon: &config.ParameterOptionCommon{
					Name:        "Relay Mode",
					Description: "Each relay will be shown, and you can enable each one individually as you see fit.\nUse this if you already know about the relays and want to customize the ones you will use.",
				},
				Value: MevSelectionMode_Relay,
			}},
			Default: map[config.Network]MevSelectionMode{
				config.Network_All: MevSelectionMode_Profile,
			},
		},

		EnableRegulatedAllMev:   generateProfileParameter(ids.MevBoostEnableRegulatedAllID, relays, true),
		EnableUnregulatedAllMev: generateProfileParameter(ids.MevBoostEnableUnregulatedAllID, relays, false),

		// Explicit relay params
		FlashbotsRelay:          generateRelayParameter(ids.MevBoostFlashbotsID, relayMap[MevRelayID_Flashbots]),
		BloxRouteMaxProfitRelay: generateRelayParameter(ids.MevBoostBloxRouteMaxProfitID, relayMap[MevRelayID_BloxrouteMaxProfit]),
		BloxRouteRegulatedRelay: generateRelayParameter(ids.MevBoostBloxRouteRegulatedID, relayMap[MevRelayID_BloxrouteRegulated]),
		EdenRelay:               generateRelayParameter(ids.MevBoostEdenID, relayMap[MevRelayID_Eden]),
		UltrasoundRelay:         generateRelayParameter(ids.MevBoostUltrasoundID, relayMap[MevRelayID_Ultrasound]),
		AestusRelay:             generateRelayParameter(ids.MevBoostAestusID, relayMap[MevRelayID_Aestus]),
		TitanGlobalRelay:        generateRelayParameter(ids.MevBoostTitanGlobalID, relayMap[MevRelayID_TitanGlobal]),
		TitanRegionalRelay:      generateRelayParameter(ids.MevBoostTitanRegionalID, relayMap[MevRelayID_TitanRegional]),

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
		&cfg.EnableRegulatedAllMev,
		&cfg.EnableUnregulatedAllMev,
		&cfg.FlashbotsRelay,
		&cfg.BloxRouteMaxProfitRelay,
		&cfg.BloxRouteRegulatedRelay,
		&cfg.EdenRelay,
		&cfg.UltrasoundRelay,
		&cfg.AestusRelay,
		&cfg.TitanGlobalRelay,
		&cfg.TitanRegionalRelay,
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

// Get the profiles that are available for the current network
func (cfg *MevBoostConfig) GetAvailableProfiles() (bool, bool) {
	regulatedAllMev := false
	unregulatedAllMev := false

	currentNetwork := cfg.parent.Network.Value
	for _, relay := range cfg.relays {
		_, exists := relay.Urls[currentNetwork]
		if !exists {
			continue
		}
		regulatedAllMev = regulatedAllMev || relay.Regulated
		unregulatedAllMev = unregulatedAllMev || !relay.Regulated
	}

	return regulatedAllMev, unregulatedAllMev
}

// Get the relays that are available for the current network
func (cfg *MevBoostConfig) GetAvailableRelays() []MevRelay {
	relays := []MevRelay{}
	currentNetwork := cfg.parent.Network.Value
	for _, relay := range cfg.relays {
		_, exists := relay.Urls[currentNetwork]
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

	currentNetwork := cfg.parent.Network.Value
	switch cfg.SelectionMode.Value {
	case MevSelectionMode_Profile:
		for _, relay := range cfg.relays {
			_, exists := relay.Urls[currentNetwork]
			if !exists {
				// Skip relays that don't exist on the current network
				continue
			}
			if relay.Regulated {
				if cfg.EnableRegulatedAllMev.Value {
					relays = append(relays, relay)
				}
			} else {
				if cfg.EnableUnregulatedAllMev.Value {
					relays = append(relays, relay)
				}
			}
		}

	case MevSelectionMode_Relay:
		if cfg.FlashbotsRelay.Value {
			_, exists := cfg.relayMap[MevRelayID_Flashbots].Urls[currentNetwork]
			if exists {
				relays = append(relays, cfg.relayMap[MevRelayID_Flashbots])
			}
		}
		if cfg.BloxRouteMaxProfitRelay.Value {
			_, exists := cfg.relayMap[MevRelayID_BloxrouteMaxProfit].Urls[currentNetwork]
			if exists {
				relays = append(relays, cfg.relayMap[MevRelayID_BloxrouteMaxProfit])
			}
		}
		if cfg.BloxRouteRegulatedRelay.Value {
			_, exists := cfg.relayMap[MevRelayID_BloxrouteRegulated].Urls[currentNetwork]
			if exists {
				relays = append(relays, cfg.relayMap[MevRelayID_BloxrouteRegulated])
			}
		}
		if cfg.EdenRelay.Value {
			_, exists := cfg.relayMap[MevRelayID_Eden].Urls[currentNetwork]
			if exists {
				relays = append(relays, cfg.relayMap[MevRelayID_Eden])
			}
		}
		if cfg.UltrasoundRelay.Value {
			_, exists := cfg.relayMap[MevRelayID_Ultrasound].Urls[currentNetwork]
			if exists {
				relays = append(relays, cfg.relayMap[MevRelayID_Ultrasound])
			}
		}
		if cfg.AestusRelay.Value {
			_, exists := cfg.relayMap[MevRelayID_Aestus].Urls[currentNetwork]
			if exists {
				relays = append(relays, cfg.relayMap[MevRelayID_Aestus])
			}
		}
		if cfg.TitanGlobalRelay.Value {
			_, exists := cfg.relayMap[MevRelayID_TitanGlobal].Urls[currentNetwork]
			if exists {
				relays = append(relays, cfg.relayMap[MevRelayID_TitanGlobal])
			}
		}
		if cfg.TitanRegionalRelay.Value {
			_, exists := cfg.relayMap[MevRelayID_TitanRegional].Urls[currentNetwork]
			if exists {
				relays = append(relays, cfg.relayMap[MevRelayID_TitanRegional])
			}
		}
	}

	return relays
}

func (cfg *MevBoostConfig) GetRelayString() string {
	relayUrls := []string{}
	currentNetwork := cfg.parent.Network.Value

	relays := cfg.GetEnabledMevRelays()
	for _, relay := range relays {
		relayUrls = append(relayUrls, relay.Urls[currentNetwork])
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
			Urls: map[config.Network]string{
				config.Network_Mainnet: "https://0xac6e77dfe25ecd6110b8e780608cce0dab71fdd5ebea22a16c0205200f2f8e2e3ad3b71d3499c54ad14d6c21b41a37ae@boost-relay.flashbots.net?id=hyperdrive",
				config.Network_Holesky: "https://0xafa4c6985aa049fb79dd37010438cfebeb0f2bd42b115b89dd678dab0670c1de38da0c4e9138c9290a398ecd9a0b3110@boost-relay-holesky.flashbots.net?id=hyperdrive",
			},
			Regulated: true,
		},

		// bloXroute Max Profit
		{
			ID:          MevRelayID_BloxrouteMaxProfit,
			Name:        "bloXroute Max Profit",
			Description: "Select this to enable the \"max profit\" relay from bloXroute.",
			Urls: map[config.Network]string{
				config.Network_Mainnet: "https://0x8b5d2e73e2a3a55c6c87b8b6eb92e0149a125c852751db1422fa951e42a09b82c142c3ea98d0d9930b056a3bc9896b8f@bloxroute.max-profit.blxrbdn.com?id=hyperdrive",
				config.Network_Holesky: "https://0x821f2a65afb70e7f2e820a925a9b4c80a159620582c1766b1b09729fec178b11ea22abb3a51f07b288be815a1a2ff516@bloxroute.holesky.blxrbdn.com?id=hyperdrive",
			},
			Regulated: false,
		},

		// bloXroute Regulated
		{
			ID:          MevRelayID_BloxrouteRegulated,
			Name:        "bloXroute Regulated",
			Description: "Select this to enable the \"regulated\" relay from bloXroute.",
			Urls: map[config.Network]string{
				config.Network_Mainnet: "https://0xb0b07cd0abef743db4260b0ed50619cf6ad4d82064cb4fbec9d3ec530f7c5e6793d9f286c4e082c0244ffb9f2658fe88@bloxroute.regulated.blxrbdn.com?id=hyperdrive",
			},
			Regulated: true,
		},

		// Eden
		{
			ID:          MevRelayID_Eden,
			Name:        "Eden Network",
			Description: "Eden Network is the home of Eden Relay, a block building hub focused on optimising block rewards for validators.",
			Urls: map[config.Network]string{
				config.Network_Mainnet: "https://0xb3ee7afcf27f1f1259ac1787876318c6584ee353097a50ed84f51a1f21a323b3736f271a895c7ce918c038e4265918be@relay.edennetwork.io?id=hyperdrive",
				config.Network_Holesky: "https://0xb1d229d9c21298a87846c7022ebeef277dfc321fe674fa45312e20b5b6c400bfde9383f801848d7837ed5fc449083a12@relay-holesky.edennetwork.io?id=hyperdrive",
			},
			Regulated: true,
		},

		// Ultrasound
		{
			ID:          MevRelayID_Ultrasound,
			Name:        "Ultra Sound",
			Description: "The ultra sound relay is a credibly-neutral and permissionless relay — a public good from the ultrasound.money team.",
			Urls: map[config.Network]string{
				config.Network_Mainnet: "https://0xa1559ace749633b997cb3fdacffb890aeebdb0f5a3b6aaa7eeeaf1a38af0a8fe88b9e4b1f61f236d2e64d95733327a62@relay.ultrasound.money?id=hyperdrive",
			},
			Regulated: false,
		},

		// Aestus
		{
			ID:          MevRelayID_Aestus,
			Name:        "Aestus",
			Description: "The Aestus MEV-Boost Relay is an independent and non-censoring relay. It is committed to neutrality and the development of a healthy MEV-Boost ecosystem.",
			Urls: map[config.Network]string{
				config.Network_Mainnet: "https://0xa15b52576bcbf1072f4a011c0f99f9fb6c66f3e1ff321f11f461d15e31b1cb359caa092c71bbded0bae5b5ea401aab7e@aestus.live?id=hyperdrive",
				config.Network_Holesky: "https://0xab78bf8c781c58078c3beb5710c57940874dd96aef2835e7742c866b4c7c0406754376c2c8285a36c630346aa5c5f833@holesky.aestus.live?id=hyperdrive",
			},
			Regulated: false,
		},

		// Titan Global
		{
			ID:          MevRelayID_TitanGlobal,
			Name:        "Titan Global (Unregulated)",
			Description: "Titan Relay is a neutral, Rust-based MEV-Boost Relay optimized for low latency through put, geographical distribution, and robustness. This is the unregulated (non-censoring) version.",
			Urls: map[config.Network]string{
				config.Network_Mainnet: "https://0x8c4ed5e24fe5c6ae21018437bde147693f68cda427cd1122cf20819c30eda7ed74f72dece09bb313f2a1855595ab677d@global.titanrelay.xyz",
				config.Network_Holesky: "https://0xaa58208899c6105603b74396734a6263cc7d947f444f396a90f7b7d3e65d102aec7e5e5291b27e08d02c50a050825c2f@holesky.titanrelay.xyz",
			},
			Regulated: false,
		},

		// Titan Regional
		{
			ID:          MevRelayID_TitanRegional,
			Name:        "Titan Regional (Regulated)",
			Description: "Titan Relay is a neutral, Rust-based MEV-Boost Relay optimized for low latency through put, geographical distribution, and robustness. This is the regulated (censoring) version.",
			Urls: map[config.Network]string{
				config.Network_Mainnet: "https://0x8c4ed5e24fe5c6ae21018437bde147693f68cda427cd1122cf20819c30eda7ed74f72dece09bb313f2a1855595ab677d@regional.titanrelay.xyz",
			},
			Regulated: true,
		},
	}

	return relays
}

// Generate one of the profile parameters
func generateProfileParameter(id string, relays []MevRelay, regulated bool) config.Parameter[bool] {
	name := "Enable "
	description := "[lime]NOTE: You can enable multiple options.\n\n"

	if regulated {
		name += "Regulated "
		description += RegulatedRelayDescription
	} else {
		name += "Unregulated "
		description += UnregulatedRelayDescription
	}

	// Generate the Mainnet description
	mainnetRelays := []string{}
	mainnetDescription := description + "\n\nRelays: "
	for _, relay := range relays {
		_, exists := relay.Urls[config.Network_Mainnet]
		if !exists {
			continue
		}
		if relay.Regulated == regulated {
			mainnetRelays = append(mainnetRelays, relay.Name)
		}
	}
	mainnetDescription += strings.Join(mainnetRelays, ", ")

	// Generate the Holesky description
	holeskyRelays := []string{}
	holeskyDescription := description + "\n\nRelays: "
	for _, relay := range relays {
		_, exists := relay.Urls[config.Network_Holesky]
		if !exists {
			continue
		}
		if relay.Regulated == regulated {
			holeskyRelays = append(holeskyRelays, relay.Name)
		}
	}
	holeskyDescription += strings.Join(holeskyRelays, ", ")

	return config.Parameter[bool]{
		ParameterCommon: &config.ParameterCommon{
			ID:                 id,
			Name:               name,
			Description:        mainnetDescription,
			AffectsContainers:  []config.ContainerID{config.ContainerID_MevBoost},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
			DescriptionsByNetwork: map[config.Network]string{
				config.Network_Mainnet: mainnetDescription,
				config.Network_Holesky: holeskyDescription,
				Network_HoleskyDev:     holeskyDescription,
			},
		},
		Default: map[config.Network]bool{
			config.Network_All: false,
		},
	}
}

// Generate one of the relay parameters
func generateRelayParameter(id string, relay MevRelay) config.Parameter[bool] {
	description := fmt.Sprintf("[lime]NOTE: You can enable multiple options.\n\n[white]%s\n\n", relay.Description)

	if relay.Regulated {
		description += "Complies with Regulations: YES\n"
	} else {
		description += "Complies with Regulations: NO\n"
	}

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
