package config

import (
	"github.com/rocket-pool/node-manager-core/config"
)

const (
	// The NodeSet dev network on Holesky
	Network_HoleskyDev config.Network = "holesky-dev"

	// Local test network for development
	Network_LocalTest config.Network = "hd-local-test"
)

type MevRelayID string

// Enum to identify MEV-boost relays
const (
	MevRelayID_Unknown            MevRelayID = ""
	MevRelayID_Flashbots          MevRelayID = "flashbots"
	MevRelayID_BloxrouteMaxProfit MevRelayID = "bloxrouteMaxProfit"
	MevRelayID_BloxrouteRegulated MevRelayID = "bloxrouteRegulated"
	MevRelayID_Eden               MevRelayID = "eden"
	MevRelayID_TitanRegional      MevRelayID = "titanRegional"
)

type MevSelectionMode string

// Enum to describe MEV-Boost relay selection mode
const (
	MevSelectionMode_All    MevSelectionMode = "all"
	MevSelectionMode_Manual MevSelectionMode = "manual"
)
