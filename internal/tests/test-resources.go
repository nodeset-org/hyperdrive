package tests

import (
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/nodeset-org/osha/beacon/db"
	"github.com/rocket-pool/node-manager-core/config"
)

func GetTestResources(beaconConfig *db.Config) *config.NetworkResources {
	return &config.NetworkResources{
		Network:            hdconfig.Network_LocalTest,
		EthNetworkName:     "local",
		ChainID:            uint(beaconConfig.ChainID),
		GenesisForkVersion: beaconConfig.GenesisForkVersion,
	}
}
