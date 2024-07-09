package testing

import (
	"github.com/ethereum/go-ethereum/common"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/nodeset-org/osha/beacon/db"
	"github.com/rocket-pool/node-manager-core/config"
)

const (
	multicallAddressString      string = "0x59F2f1fCfE2474fD5F0b9BA1E73ca90b143Eb8d0"
	balanceBatcherAddressString string = "0xC6bA8C3233eCF65B761049ef63466945c362EdD2"
)

// Returns a network resources instance with local testing network values
func GetTestResources(beaconConfig *db.Config, nodesetUrl string) *hdconfig.HyperdriveResources {
	return &hdconfig.HyperdriveResources{
		NetworkResources: &config.NetworkResources{
			Network:               hdconfig.Network_LocalTest,
			EthNetworkName:        "localtest",
			ChainID:               uint(beaconConfig.ChainID),
			GenesisForkVersion:    beaconConfig.GenesisForkVersion,
			MulticallAddress:      common.HexToAddress(multicallAddressString),
			BalanceBatcherAddress: common.HexToAddress(balanceBatcherAddressString),
		},
		NodeSetApiUrl: nodesetUrl,
	}
}
