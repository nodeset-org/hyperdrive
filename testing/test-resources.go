package testing

import (
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/nodeset-org/osha/beacon/db"
	"github.com/rocket-pool/node-manager-core/config"
)

const (
	TestNetworkEthName string = "hardhat"
)

// Creates a new set of network settings designed for usage in local testing with Hardat.
// The settings are incomplete; things like the multicall address and balance batcher address should be set by the user if needed.
func GetDefaultTestNetworkSettings(beaconConfig *db.Config) *config.NetworkSettings {
	return &config.NetworkSettings{
		Key:         hdconfig.Network_LocalTest,
		Name:        "Local Test Network",
		Description: "Local test network for development and testing",
		NetworkResources: &config.NetworkResources{
			EthNetworkName:     TestNetworkEthName,
			ChainID:            uint(beaconConfig.ChainID),
			GenesisForkVersion: beaconConfig.GenesisForkVersion,
		},
		DefaultConfigSettings: map[string]any{},
	}
}

// Returns a network resources instance with local testing network values
func getTestResources(networkResources *config.NetworkResources, nodesetUrl string, deploymentName string) *hdconfig.MergedResources {
	return &hdconfig.MergedResources{
		NetworkResources: networkResources,
		HyperdriveResources: &hdconfig.HyperdriveResources{
			NodeSetApiUrl:  nodesetUrl,
			DeploymentName: deploymentName,
		},
	}
}
