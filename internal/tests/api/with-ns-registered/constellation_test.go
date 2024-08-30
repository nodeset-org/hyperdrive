package with_ns_registered

import (
	"math/big"
	"runtime/debug"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	hdtesting "github.com/nodeset-org/hyperdrive-daemon/testing"
	"github.com/nodeset-org/osha/keys"
	"github.com/rocket-pool/node-manager-core/utils"
	"github.com/stretchr/testify/require"
)

const (
	whitelistAddressString     string = "0xA9e6Bfa2BF53dE88FEb19761D9b2eE2e821bF1Bf"
	expectedWhitelistSignature string = "0xf2b73cd729a9b15e8f17ce0189c4ddfe63ad35917f63e2b1ffa7ea1dc527bdf535ba05ba44d2dce733096b8c389472e81a4548b1d75a600633c4ac4bcb8e7c6f1b"
	expectedMinipoolCount      int    = 10
)

// Test getting a signature for whitelisting a node
func TestConstellationWhitelistSignature(t *testing.T) {
	// Take a snapshot, revert at the end
	snapshotName, err := testMgr.CreateCustomSnapshot(hdtesting.Service_EthClients | hdtesting.Service_Filesystem | hdtesting.Service_NodeSet)
	if err != nil {
		fail("Error creating custom snapshot: %v", err)
	}
	defer nodeset_cleanup(snapshotName)

	// Get the private key for the Constellation deployer (the admin)
	keygen, err := keys.NewKeyGeneratorWithDefaults()
	require.NoError(t, err)
	adminKey, err := keygen.GetEthPrivateKey(0)
	require.NoError(t, err)

	// Set up the nodeset.io mock
	res := testMgr.GetNode().GetServiceProvider().GetResources()
	nsMgr := testMgr.GetNodeSetMockServer().GetManager()
	nsDB := nsMgr.GetDatabase()
	deployment := nsDB.Constellation.AddDeployment(
		res.DeploymentName,
		new(big.Int).SetUint64(uint64(res.ChainID)),
		common.HexToAddress(whitelistAddressString),
		common.Address{},
	)
	deployment.SetAdminPrivateKey(adminKey)

	// Get a whitelist signature
	hd := hdNode.GetApiClient()
	response, err := hd.NodeSet_Constellation.GetRegistrationSignature()
	require.NoError(t, err)
	require.False(t, response.Data.NotAuthorized)
	require.False(t, response.Data.NotRegistered)
	sigHex := utils.EncodeHexWithPrefix(response.Data.Signature)
	require.Equal(t, expectedWhitelistSignature, sigHex)
	t.Logf("Whitelist signature is correct")
}

// Cleanup after a unit test
func nodeset_cleanup(snapshotName string) {
	// Handle panics
	r := recover()
	if r != nil {
		debug.PrintStack()
		fail("Recovered from panic: %v", r)
	}

	// Revert to the snapshot taken at the start of the test
	if snapshotName != "" {
		err := testMgr.RevertToCustomSnapshot(snapshotName)
		if err != nil {
			fail("Error reverting to custom snapshot: %v", err)
		}
	}
}
