package with_ns_registered

import (
	"runtime/debug"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	hdtesting "github.com/nodeset-org/hyperdrive-daemon/testing"
	"github.com/nodeset-org/osha/keys"
	"github.com/rocket-pool/node-manager-core/utils"
	"github.com/stretchr/testify/require"
)

const (
	whitelistAddressString     string = "0x3fdc08D815cc4ED3B7F69Ee246716f2C8bCD6b07"
	expectedWhitelistSignature string = "0x38b8b7989b0f0695e8cb08253a3077b33c102c307a6b0a50c0a90dbb971ecd3f0e3fd38cd9ffecf696f1aa41151dceb86b405ed78259928c9be62f56d852e03c1b"
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
	whitelistAddress := common.HexToAddress(whitelistAddressString)
	nsMgr := testMgr.GetNodeSetMockServer().GetManager()
	nsMgr.SetConstellationAdminPrivateKey(adminKey)
	nsMgr.SetConstellationWhitelistAddress(whitelistAddress)

	// Get a whitelist signature
	hd := testMgr.GetApiClient()
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
