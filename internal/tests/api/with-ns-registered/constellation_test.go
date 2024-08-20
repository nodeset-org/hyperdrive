package with_ns_registered

import (
	"math/big"
	"runtime/debug"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	hdtesting "github.com/nodeset-org/hyperdrive-daemon/testing"
	"github.com/nodeset-org/nodeset-client-go/server-mock/db"
	"github.com/nodeset-org/osha/keys"
	"github.com/rocket-pool/node-manager-core/utils"
	"github.com/stretchr/testify/require"
)

const (
	whitelistTimestamp         int64  = 1721417393
	whitelistAddressString     string = "0xA9e6Bfa2BF53dE88FEb19761D9b2eE2e821bF1Bf"
	expectedWhitelistSignature string = "0x8d6779cdc17bbfd0416fce5af7e4bde2b106ea5904d4c532eee8dfd73e60019b08c35f86b5dd94713b1dc30fa2fc8f91dd1bd32ab2592c22bfade08bfab3817d1b"
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
	manualTime := time.Unix(whitelistTimestamp, 0)
	nsMgr := testMgr.GetNodeSetMockServer().GetManager()
	nsMgr.SetConstellationAdminPrivateKey(adminKey)
	nsMgr.SetManualSignatureTimestamp(&manualTime)
	nsMgr.SetDeployment(&db.Deployment{
		DeploymentID:     res.DeploymentName,
		WhitelistAddress: common.HexToAddress(whitelistAddressString),
		ChainID:          new(big.Int).SetUint64(uint64(res.ChainID)),
	})

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

func TestGetMinipoolAvailabilityCount(t *testing.T) {
	// Take a snapshot, revert at the end
	snapshotName, err := testMgr.CreateCustomSnapshot(hdtesting.Service_EthClients | hdtesting.Service_Filesystem | hdtesting.Service_NodeSet)
	if err != nil {
		fail("Error creating custom snapshot: %v", err)
	}
	defer nodeset_cleanup(snapshotName)

	// Set the minipool count
	res := testMgr.GetNode().GetServiceProvider().GetResources()
	nsmock := testMgr.GetNodeSetMockServer()
	nsMgr := nsmock.GetManager()
	nsMgr.SetAvailableConstellationMinipoolCount(nsEmail, expectedMinipoolCount)
	nsMgr.SetDeployment(&db.Deployment{
		DeploymentID: res.DeploymentName,
		ChainID:      new(big.Int).SetUint64(uint64(res.ChainID)),
	})

	// Get the minipool count and assert
	minipoolCountResponse, err := hdNode.GetApiClient().NodeSet_Constellation.GetAvailableMinipoolCount()
	require.NoError(t, err)
	require.Equal(t, expectedMinipoolCount, minipoolCountResponse.Data.Count)

	t.Log("Minipool availability count is correct")

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
