package constconfig

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/rocket-pool/node-manager-core/config"
)

// A collection of network-specific resources and getters for them
type ConstellationResources struct {
	*config.NetworkResources

	// The address of the NodeSet fee recipient
	FeeRecipient common.Address
}

// Creates a new resource collection for the given network
func NewConstellationResources(network config.Network) *ConstellationResources {
	// Mainnet
	mainnetResources := &ConstellationResources{
		NetworkResources: config.NewResources(config.Network_Mainnet),
		FeeRecipient:     common.HexToAddress(""),
	}

	// Holesky
	holeskyResources := &ConstellationResources{
		NetworkResources: config.NewResources(config.Network_Holesky),
		FeeRecipient:     common.HexToAddress("0xc98F25BcAA6B812a07460f18da77AF8385be7b56"),
	}

	// Holesky Dev
	holeskyDevResources := &ConstellationResources{
		NetworkResources: config.NewResources(config.Network_Holesky),
		FeeRecipient:     common.HexToAddress("0xc98F25BcAA6B812a07460f18da77AF8385be7b56"),
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
