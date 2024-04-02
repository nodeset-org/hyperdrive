package constconfig

import (
	"fmt"

	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/rocket-pool/node-manager-core/config"
)

// A collection of network-specific resources and getters for them
type ConstellationResources struct {
	*config.NetworkResources
}

// Creates a new resource collection for the given network
func NewConstellationResources(network config.Network) *ConstellationResources {
	// Mainnet
	mainnetResources := &ConstellationResources{
		NetworkResources: config.NewResources(config.Network_Mainnet),
	}

	// Holesky
	holeskyResources := &ConstellationResources{
		NetworkResources: config.NewResources(config.Network_Holesky),
	}

	// Holesky Dev
	holeskyDevResources := &ConstellationResources{
		NetworkResources: config.NewResources(config.Network_Holesky),
	}
	holeskyDevResources.Network = hdconfig.Network_HoleskyDev

	switch network {
	case config.Network_Mainnet:
		return mainnetResources
	case config.Network_Holesky:
		return holeskyResources
	case hdconfig.Network_HoleskyDev:
		return holeskyDevResources
	}

	panic(fmt.Sprintf("network %s is not supported", network))
}
