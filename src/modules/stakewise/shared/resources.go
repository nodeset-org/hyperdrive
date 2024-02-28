package swshared

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/shared/config"
)

// A collection of network-specific resources and getters for them
type StakewiseResources struct {
	// The Network being used
	Network config.Network

	// The address of the Stakewise vault
	Vault common.Address

	// The address of the NodeSet fee recipient
	FeeRecipient common.Address

	// The genesis fork version for the network according to the Beacon config for the network
	GenesisForkVersion []byte

	// The URL for the NodeSet API server
	NodesetApiUrl string

	// The string to put in requests for the network param
	NodesetNetwork string
}

// Creates a new resource collection for the given network
func NewStakewiseResources(network config.Network) *StakewiseResources {
	// Mainnet
	mainnetResources := &StakewiseResources{
		Network:            network,
		Vault:              common.HexToAddress(""),
		FeeRecipient:       common.HexToAddress(""),
		GenesisForkVersion: common.FromHex("0x00000000"), // https://github.com/eth-clients/eth2-networks/tree/master/shared/mainnet#genesis-information
		NodesetApiUrl:      "",
		NodesetNetwork:     "mainnet",
	}

	// Holesky
	holeskyResources := &StakewiseResources{
		Network:            network,
		Vault:              common.HexToAddress("0x646F5285D195e08E309cF9A5aDFDF68D6Fcc51C4"),
		FeeRecipient:       common.HexToAddress("0xc98F25BcAA6B812a07460f18da77AF8385be7b56"),
		GenesisForkVersion: common.FromHex("0x01017000"), // https://github.com/eth-clients/holesky
		NodesetApiUrl:      "https://staging.nodeset.io/api",
		NodesetNetwork:     "holesky",
	}

	// Holesky Dev
	holeskyDevResources := &StakewiseResources{
		Network:            network,
		Vault:              common.HexToAddress("0xf8763855473ce978232bBa37ef90fcFc8aAE10d1"),
		FeeRecipient:       common.HexToAddress("0xc98F25BcAA6B812a07460f18da77AF8385be7b56"),
		GenesisForkVersion: common.FromHex("0x01017000"), // https://github.com/eth-clients/holesky
		NodesetApiUrl:      "https://staging.nodeset.io/api",
		NodesetNetwork:     "holesky",
	}

	switch network {
	case config.Network_Mainnet:
		return mainnetResources
	case config.Network_Holesky:
		return holeskyResources
	case config.Network_HoleskyDev:
		return holeskyDevResources
	}

	panic(fmt.Sprintf("network %s is not supported", network))
}
