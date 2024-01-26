package utils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

// A collection of network-specific resources and getters for them
type Resources struct {
	// The address of the multicall contract
	MulticallAddress common.Address

	// The chain ID for the current network
	ChainID uint

	// The network being used
	network types.Network
}

// Creates a new resource collection for the given network
func NewResources(network types.Network) *Resources {
	// Mainnet
	mainnetResources := &Resources{
		network:          network,
		MulticallAddress: common.HexToAddress("0x5BA1e12693Dc8F9c48aAD8770482f4739bEeD696"),
		ChainID:          1,
	}

	// Holesky
	holeskyResources := &Resources{
		network:          network,
		MulticallAddress: common.HexToAddress("0x0540b786f03c9491f3a2ab4b0e3ae4ecd4f63ce7"),
		ChainID:          17000,
	}

	switch network {
	case types.Network_Mainnet:
		return mainnetResources
	case types.Network_Holesky, types.Network_HoleskyDev:
		return holeskyResources
	}

	panic(fmt.Sprintf("network %s is not supported", network))
}
