package testing

import (
	"github.com/ethereum/go-ethereum/common"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/nodeset-org/osha/beacon/db"
	"github.com/rocket-pool/node-manager-core/config"
)

const (
	multicallAddressString      string = "0x05Aa229Aec102f78CE0E852A812a388F076Aa555"
	balanceBatcherAddressString string = "0x0b48aF34f4c854F5ae1A3D587da471FeA45bAD52"
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
