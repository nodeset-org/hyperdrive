package swshared

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

// A collection of network-specific resources and getters for them
type StakewiseResources struct {
	// The Network being used
	Network types.Network

	// The address of the Stakewise vault
	Vault common.Address

	// The address of the NodeSet fee recipient
	FeeRecipient common.Address

	// The genesis fork version for the network according to the Beacon config for the network
	GenesisForkVersion []byte

	// The URL for uploading deposit data to NodeSet
	NodesetDepositUrl string
}

// Creates a new resource collection for the given network
func NewStakewiseResources(network types.Network) *StakewiseResources {
	// Mainnet
	mainnetResources := &StakewiseResources{
		Network:            network,
		Vault:              common.HexToAddress(""),
		FeeRecipient:       common.HexToAddress(""),
		GenesisForkVersion: common.FromHex("0x00000000"), // https://github.com/eth-clients/eth2-networks/tree/master/shared/mainnet#genesis-information
		NodesetDepositUrl:  "",
	}

	// Holesky
	holeskyResources := &StakewiseResources{
		Network:            network,
		Vault:              common.HexToAddress("0x646F5285D195e08E309cF9A5aDFDF68D6Fcc51C4"),
		FeeRecipient:       common.HexToAddress("0xc98F25BcAA6B812a07460f18da77AF8385be7b56"),
		GenesisForkVersion: common.FromHex("0x01017000"), // https://github.com/eth-clients/holesky
		NodesetDepositUrl:  "https://staging.nodeset.io/api/deposit-data",
	}

	// Holesky Dev
	holeskyDevResources := &StakewiseResources{
		Network:            network,
		Vault:              common.HexToAddress("0x01b353Abc66A65c4c0Ac9c2ecF82e693Ce0303Bc"),
		FeeRecipient:       common.HexToAddress("0xc98F25BcAA6B812a07460f18da77AF8385be7b56"),
		GenesisForkVersion: common.FromHex("0x01017000"), // https://github.com/eth-clients/holesky
		NodesetDepositUrl:  "https://staging.nodeset.io/api/deposit-data",
	}

	switch network {
	case types.Network_Mainnet:
		return mainnetResources
	case types.Network_Holesky:
		return holeskyResources
	case types.Network_HoleskyDev:
		return holeskyDevResources
	}

	panic(fmt.Sprintf("network %s is not supported", network))
}
