package utils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

// A collection of network-specific resources and getters for them
type Resources struct {
	// The network being used
	network types.Network
}

// Creates a new resource collection for the given network
func NewResources(network types.Network) *Resources {
	return &Resources{
		network: network,
	}
}

// The address of the multicall contract
func (r *Resources) GetMulticallAddress() common.Address {
	switch r.network {
	case types.Network_Mainnet:
		return common.HexToAddress("0x5BA1e12693Dc8F9c48aAD8770482f4739bEeD696")
	case types.Network_Holesky:
		return common.HexToAddress("0x0540b786f03c9491f3a2ab4b0e3ae4ecd4f63ce7")
	case types.Network_HoleskyDev:
		return common.HexToAddress("0x0540b786f03c9491f3a2ab4b0e3ae4ecd4f63ce7")
	default:
		panic(fmt.Sprintf("network %s doesn't have a multicall address", r.network))
	}
}
