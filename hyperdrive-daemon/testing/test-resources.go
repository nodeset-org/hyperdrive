package testing

import (
	"filippo.io/age"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/osha/beacon/db"
	"github.com/rocket-pool/node-manager-core/config"
)

const (
	// ETH network name to use for test networks
	TestNetworkEthName string = "hardhat"

	// Serialized nodeset.io deployment dummy encryption ID
	EncryptionIdentityString string = "AGE-SECRET-KEY-19N32FTRU5JJ66DNTVE8NTTE04CUQC3R3FC5QD9QKA97AZWCUW74ST78LD3"
)

var (
	// Nodeset.io deployment dummy encryption ID
	EncryptionIdentity, _ = age.ParseX25519Identity(EncryptionIdentityString)
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
func getTestResources(networkResources *config.NetworkResources, nodesetUrl string) *hdconfig.MergedResources {
	return &hdconfig.MergedResources{
		NetworkResources: networkResources,
		HyperdriveResources: &hdconfig.HyperdriveResources{
			NodeSetApiUrl:    nodesetUrl,
			EncryptionPubkey: EncryptionIdentity.Recipient().String(),
		},
	}
}
