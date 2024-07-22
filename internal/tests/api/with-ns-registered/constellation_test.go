package with_ns_registered

import (
	"runtime/debug"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	hdtesting "github.com/nodeset-org/hyperdrive-daemon/testing"
	"github.com/nodeset-org/osha/keys"
	"github.com/rocket-pool/node-manager-core/utils"
	"github.com/stretchr/testify/require"
)

const (
	whitelistTimestamp         int64  = 1721417393
	whitelistAddressString     string = "0x1E3b98102e19D3a164d239BdD190913C2F02E756"
	expectedWhitelistSignature string = "0xdd45a03d896d93e4fd2ee947bed23fb4f87a24d528cd5ecfe847f4c521cba8c1519f4fbc74d9a12d40fa64244a0616370ae709394a0217659d028351bb8dc3c21b"
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
	manualTime := time.Unix(whitelistTimestamp, 0)
	whitelistAddress := common.HexToAddress(whitelistAddressString)
	nsMgr := testMgr.GetNodeSetMockServer().GetManager()
	nsMgr.SetConstellationAdminPrivateKey(adminKey)
	nsMgr.SetManualSignatureTimestamp(&manualTime)

	// Get a whitelist signature
	hd := testMgr.GetApiClient()
	response, err := hd.NodeSet_Constellation.GetRegistrationSignature(whitelistAddress)
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
	nsmock := testMgr.GetNodeSetMockServer()
	nsmock.GetManager().SetAvailableConstellationMinipoolCount(nodeAddress, expectedMinipoolCount)

	// Get the minipool count and assert
	minipoolCountResponse, err := testMgr.GetApiClient().NodeSet_Constellation.GetAvailableMinipoolCount()
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
