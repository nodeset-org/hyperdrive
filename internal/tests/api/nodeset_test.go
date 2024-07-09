package api_test

import (
	"runtime/debug"
	"testing"

	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	hdtesting "github.com/nodeset-org/hyperdrive-daemon/testing"
	"github.com/nodeset-org/osha/keys"
	"github.com/rocket-pool/node-manager-core/wallet"
	"github.com/stretchr/testify/require"
)

const (
	nsEmail string = "test@nodeset.io"
)

// Test registration with nodeset.io if the node doesn't have a wallet yet
func TestNodeSetRegistration_NoWallet(t *testing.T) {
	// Take a snapshot, revert at the end
	snapshotName, err := testMgr.CreateCustomSnapshot(hdtesting.Service_Filesystem | hdtesting.Service_NodeSet)
	if err != nil {
		fail("Error creating custom snapshot: %v", err)
	}
	defer nodeset_cleanup(snapshotName)

	// Run the round-trip test
	hd := testMgr.GetApiClient()
	response, err := hd.NodeSet.GetRegistrationStatus()
	require.NoError(t, err)
	require.Equal(t, api.NodeSetRegistrationStatus_NoWallet, response.Data.Status)
	t.Logf("Node has no wallet, registration status is correct")
}

// Test registration with nodeset.io if the node has a wallet but hasn't been registered yet
func TestNodeSetRegistration_NoRegistration(t *testing.T) {
	// Take a snapshot, revert at the end
	snapshotName, err := testMgr.CreateCustomSnapshot(hdtesting.Service_Filesystem | hdtesting.Service_NodeSet)
	if err != nil {
		fail("Error creating custom snapshot: %v", err)
	}
	defer nodeset_cleanup(snapshotName)

	// Recover a wallet
	derivationPath := string(wallet.DerivationPath_Default)
	index := uint64(0)
	recoverResponse, err := testMgr.GetApiClient().Wallet.Recover(&derivationPath, keys.DefaultMnemonic, &index, goodPassword, true)
	require.NoError(t, err)
	t.Log("Recover called")

	// Check the response
	require.Equal(t, expectedWalletAddress, recoverResponse.Data.AccountAddress)
	t.Log("Received correct wallet address")

	// Run the round-trip test
	hd := testMgr.GetApiClient()
	registrationResponse, err := hd.NodeSet.GetRegistrationStatus()
	require.NoError(t, err)
	require.Equal(t, api.NodeSetRegistrationStatus_Unregistered, registrationResponse.Data.Status)
	t.Logf("Node has a wallet but isn't registered, registration status is correct")
}

// Test registration with nodeset.io if the node has a wallet and has been registered
func TestNodeSetRegistration_Registered(t *testing.T) {
	// Take a snapshot, revert at the end
	snapshotName, err := testMgr.CreateCustomSnapshot(hdtesting.Service_Filesystem | hdtesting.Service_NodeSet)
	if err != nil {
		fail("Error creating custom snapshot: %v", err)
	}
	defer nodeset_cleanup(snapshotName)

	// Recover a wallet
	derivationPath := string(wallet.DerivationPath_Default)
	index := uint64(0)
	recoverResponse, err := testMgr.GetApiClient().Wallet.Recover(&derivationPath, keys.DefaultMnemonic, &index, goodPassword, true)
	require.NoError(t, err)
	t.Log("Recover called")

	// Check the response
	require.Equal(t, expectedWalletAddress, recoverResponse.Data.AccountAddress)
	t.Log("Received correct wallet address")

	// Register the node with nodeset.io
	hd := testMgr.GetApiClient()
	nsMgr := testMgr.GetNodeSetMockServer().GetManager()
	err = nsMgr.AddUser(nsEmail)
	require.NoError(t, err)
	err = nsMgr.WhitelistNodeAccount(nsEmail, expectedWalletAddress)
	require.NoError(t, err)
	registerResponse, err := hd.NodeSet.RegisterNode(nsEmail)
	require.NoError(t, err)
	require.True(t, registerResponse.Data.Success)

	// Run the round-trip test
	registrationResponse, err := hd.NodeSet.GetRegistrationStatus()
	require.NoError(t, err)
	require.Equal(t, api.NodeSetRegistrationStatus_Registered, registrationResponse.Data.Status)
	t.Logf("Node is registered with nodeset.io")
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

	// Reload the wallet to undo any changes made during the test
	err := testMgr.GetServiceProvider().GetWallet().Reload(testMgr.GetLogger())
	if err != nil {
		fail("Error reloading wallet: %v", err)
	}
}
