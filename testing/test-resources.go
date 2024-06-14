package testing

import (
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/nodeset-org/osha/beacon/db"
	"github.com/rocket-pool/node-manager-core/config"
)

// Returns a network resources instance with local testing network values
func GetTestResources(beaconConfig *db.Config) *config.NetworkResources {
	return &config.NetworkResources{
		Network:            hdconfig.Network_LocalTest,
		EthNetworkName:     "localtest",
		ChainID:            uint(beaconConfig.ChainID),
		GenesisForkVersion: beaconConfig.GenesisForkVersion,
	}
}
