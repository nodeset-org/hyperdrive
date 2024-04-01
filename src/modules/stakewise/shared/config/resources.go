package swconfig

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/rocket-pool/node-manager-core/config"
)

// A collection of network-specific resources and getters for them
type StakewiseResources struct {
	*config.NetworkResources

	// The address of the Stakewise vault
	Vault common.Address

	// The address of the NodeSet fee recipient
	FeeRecipient common.Address

	// The URL for the NodeSet API server
	NodesetApiUrl string

	// The address of the SplitMain contract
	Splitmain common.Address

	// The amount of ETH to claim
	ClaimEthAmount *big.Int

	// The list of token addresses to claim
	ClaimTokenList []common.Address
}

// Creates a new resource collection for the given network
func NewStakewiseResources(network config.Network) *StakewiseResources {
	// Mainnet
	mainnetResources := &StakewiseResources{
		NetworkResources: config.NewResources(config.Network_Mainnet),
		Vault:            common.HexToAddress(""),
		FeeRecipient:     common.HexToAddress(""),
		NodesetApiUrl:    "",
		Splitmain:        common.HexToAddress(""),
		ClaimEthAmount:   big.NewInt(0),      // 0 => claim all
		ClaimTokenList:   []common.Address{}, // TODO: Get list from Wander
	}

	// Holesky
	holeskyResources := &StakewiseResources{
		NetworkResources: config.NewResources(config.Network_Holesky),
		Vault:            common.HexToAddress("0x646F5285D195e08E309cF9A5aDFDF68D6Fcc51C4"),
		FeeRecipient:     common.HexToAddress("0xc98F25BcAA6B812a07460f18da77AF8385be7b56"),
		NodesetApiUrl:    "https://staging.nodeset.io/api",
		Splitmain:        common.HexToAddress("0x2ed6c4B5dA6378c7897AC67Ba9e43102Feb694EE"),
		ClaimEthAmount:   big.NewInt(0),      // 0 => claim all
		ClaimTokenList:   []common.Address{}, // TODO: Get list from Wander
	}

	// Holesky Dev
	holeskyDevResources := &StakewiseResources{
		NetworkResources: config.NewResources(config.Network_Holesky),
		Vault:            common.HexToAddress("0xf8763855473ce978232bBa37ef90fcFc8aAE10d1"),
		FeeRecipient:     common.HexToAddress("0xc98F25BcAA6B812a07460f18da77AF8385be7b56"),
		NodesetApiUrl:    "https://staging.nodeset.io/api",
		Splitmain:        common.HexToAddress("0x2ed6c4B5dA6378c7897AC67Ba9e43102Feb694EE"),
		ClaimEthAmount:   big.NewInt(0),      // 0 => claim all
		ClaimTokenList:   []common.Address{}, // TODO: Get list from Wander
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
