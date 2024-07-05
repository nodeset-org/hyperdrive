package config

import (
	"fmt"

	"github.com/rocket-pool/node-manager-core/config"
)

// A collection of network-specific resources and getters for them
type HyperdriveResources struct {
	*config.NetworkResources

	// The URL for the NodeSet API server
	NodeSetApiUrl string
}

// Creates a new resource collection for the given network
func NewHyperdriveResources(network config.Network) *HyperdriveResources {
	// Mainnet
	mainnetResources := &HyperdriveResources{
		NetworkResources: config.NewResources(config.Network_Mainnet),
		NodeSetApiUrl:    "https://nodeset.io/api",
	}

	// Holesky
	holeskyResources := &HyperdriveResources{
		NetworkResources: config.NewResources(config.Network_Holesky),
		NodeSetApiUrl:    "https://nodeset.io/api",
	}

	// Holesky Dev
	holeskyDevResources := &HyperdriveResources{
		NetworkResources: config.NewResources(config.Network_Holesky),
		NodeSetApiUrl:    "https://staging.nodeset.io/api",
	}
	holeskyDevResources.Network = Network_HoleskyDev

	switch network {
	case config.Network_Mainnet:
		return mainnetResources
	case config.Network_Holesky:
		return holeskyResources
	case Network_HoleskyDev:
		return holeskyDevResources
	}

	panic(fmt.Sprintf("network %s is not supported", network))
}
