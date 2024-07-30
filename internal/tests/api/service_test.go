package api_test

import (
	"runtime/debug"
	"testing"
	"time"

	dtypes "github.com/docker/docker/api/types"
	"github.com/nodeset-org/hyperdrive-daemon/shared"
	hdtesting "github.com/nodeset-org/hyperdrive-daemon/testing"
	"github.com/stretchr/testify/require"
)

// Test getting the client status of synced clients
func TestClientStatus_Synced(t *testing.T) {
	// Take a snapshot, revert at the end
	snapshotName, err := testMgr.CreateCustomSnapshot(hdtesting.Service_EthClients)
	if err != nil {
		fail("Error creating custom snapshot: %v", err)
	}
	defer service_cleanup(snapshotName)

	// Commit a block just so the latest block is fresh - otherwise the sync progress check will
	// error out because the block is too old and it thinks the client just can't find any peers
	err = testMgr.CommitBlock()
	if err != nil {
		t.Fatalf("Error committing block: %v", err)
	}

	// Run the round-trip test
	response, err := hdNode.GetApiClient().Service.ClientStatus()
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

// Test getting the server version
func TestServerVersion(t *testing.T) {
	defer service_cleanup("")

	version := shared.HyperdriveVersion

	// Run the round-trip test
	response, err := hdNode.GetApiClient().Service.Version()
	require.NoError(t, err)
	require.Equal(t, version, response.Data.Version)
	t.Logf("Received correct version: %s", version)
}

func TestRestartContainer(t *testing.T) {
	// Take a snapshot, revert at the end
	snapshotName, err := testMgr.CreateCustomSnapshot(hdtesting.Service_Docker)
	if err != nil {
		fail("Error creating custom snapshot: %v", err)
	}
	defer service_cleanup(snapshotName)

	// Get some services
	sp := hdNode.GetServiceProvider()
	cfg := sp.GetConfig()
	ctx := sp.GetBaseContext()

	// Create a fake VC directly
	containerName := "mock_vc"
	fullName := cfg.GetDockerArtifactName(containerName)
	oneMinuteAgo := time.Now().Add(-1 * time.Minute)
	oneMinuteAgoStr := oneMinuteAgo.Format(time.RFC3339Nano)
	vc := dtypes.ContainerJSON{
		ContainerJSONBase: &dtypes.ContainerJSONBase{
			Name:    fullName,
			Created: oneMinuteAgoStr,
			State: &dtypes.ContainerState{
				Running:    false,
				StartedAt:  oneMinuteAgoStr,
				FinishedAt: oneMinuteAgoStr,
			},
		},
	}
	docker := testMgr.GetDockerMockManager()
	err = docker.Mock_AddContainer(vc)
	if err != nil {
		t.Fatalf("Error creating mock VC: %v", err)
	}
	t.Log("Created mock VC")

	// Run the client call
	_, err = hdNode.GetApiClient().Service.RestartContainer(containerName)
	require.NoError(t, err)
	t.Log("Restart called")

	// Make sure the restart actually worked
	vc, err = docker.ContainerInspect(ctx, fullName)
	if err != nil {
		t.Fatalf("Error inspecting mock VC: %v", err)
	}
	require.True(t, vc.State.Running)
	require.Greater(t, vc.State.StartedAt, oneMinuteAgoStr)
	t.Logf("VC restart was successful - original start = %s, new start = %s", oneMinuteAgoStr, vc.State.StartedAt)
}

func service_cleanup(snapshotName string) {
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
