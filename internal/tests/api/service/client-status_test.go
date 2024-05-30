package service

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test getting the client status of synced clients
func TestClientStatus_Synced(t *testing.T) {
	// Set up panic handling and cleanup
	defer func() {
		r := recover()
		if r != nil {
			debug.PrintStack()
			fail("Recovered from panic: %v", r)
		} else {
			testMgr.RevertToBaseline()
		}
	}()

	// Take a snapshot, revert at the end
	snapshotName, err := testMgr.CreateCustomSnapshot()
	if err != nil {
		t.Fatalf("Error creating custom snapshot: %v", err)
	}
	defer func() {
		err := testMgr.RevertToCustomSnapshot(snapshotName)
		if err != nil {
			t.Fatalf("Error reverting to custom snapshot: %v", err)
		}
	}()

	// Commit a block just so the latest block is fresh - otherwise the sync progress check will
	// error out because the block is too old and it thinks the client just can't find any peers
	err = testMgr.CommitBlock()
	if err != nil {
		t.Fatalf("Error committing block: %v", err)
	}

	// Run the round-trip test
	response, err := apiClient.Service.ClientStatus()
	require.NoError(t, err)
	require.True(t, response.Data.EcManagerStatus.PrimaryClientStatus.IsSynced)
	require.True(t, response.Data.EcManagerStatus.PrimaryClientStatus.IsWorking)
	require.Equal(t, 1.0, response.Data.EcManagerStatus.PrimaryClientStatus.SyncProgress)
	require.Equal(t, "", response.Data.EcManagerStatus.PrimaryClientStatus.Error)
	require.False(t, response.Data.EcManagerStatus.FallbackEnabled)

	require.True(t, response.Data.BcManagerStatus.PrimaryClientStatus.IsSynced)
	require.True(t, response.Data.BcManagerStatus.PrimaryClientStatus.IsWorking)
	require.Equal(t, 1.0, response.Data.BcManagerStatus.PrimaryClientStatus.SyncProgress)
	require.Equal(t, "", response.Data.BcManagerStatus.PrimaryClientStatus.Error)
	require.False(t, response.Data.BcManagerStatus.FallbackEnabled)

	require.Equal(t, response.Data.EcManagerStatus.PrimaryClientStatus.ChainId, response.Data.BcManagerStatus.PrimaryClientStatus.ChainId)
	t.Logf("Received correct client status - both clients synced on chain %d", response.Data.EcManagerStatus.PrimaryClientStatus.ChainId)
}
