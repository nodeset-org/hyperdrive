package tests

import (
	"os"
	"sync"
	"testing"

	hdclient "github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/context"
	hdcommon "github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common"
	hdserver "github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server"
	"github.com/nodeset-org/hyperdrive/shared/auth"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/osha"
	"github.com/rocket-pool/node-manager-core/config"
	"github.com/stretchr/testify/require"
)

const (
	hdTestApiKey string = "hd-api-key"
)

func TestNewConfig_Holesky(t *testing.T) {
	// Take a snapshot, revert at the end
	snapshotName, err := testMgr.CreateCustomSnapshot(osha.Service_Filesystem)
	if err != nil {
		fail("Error creating custom snapshot: %v", err)
	}
	defer basicTestCleanup(snapshotName)

	cfgPath := testMgr.GetTestDir()

	// Make a new Hyperdrive context
	os.Setenv(context.TestSystemDirEnvVar, "../../install/deploy")
	hdCtx := context.NewHyperdriveContext(cfgPath, nil)

	// Load the network settings
	hdNetworkSettings, err := hdconfig.LoadSettingsFiles(hdCtx.NetworksDir)
	require.NoError(t, err)
	hdCtx.HyperdriveNetworkSettings = hdNetworkSettings

	// Make a new Hyperdrive client
	hdClient, err := hdclient.NewHyperdriveClientFromHyperdriveCtx(hdCtx)
	require.NoError(t, err)

	// Make a new config
	cfg, isNewCfg, err := hdClient.LoadConfig()
	require.NoError(t, err)
	require.True(t, isNewCfg)
	t.Log("Config initialized")

	// Check the list of networks
	networkOptions := cfg.Hyperdrive.Network.Options
	foundMainnet := false
	foundHolesky := false
	for _, network := range networkOptions {
		switch network.Value {
		case config.Network_Mainnet:
			foundMainnet = true
		case config.Network_Holesky:
			foundHolesky = true
		}
	}
	require.True(t, foundMainnet)
	require.True(t, foundHolesky)
	t.Log("Network options loaded successfully")

	// Set the network to Holesky and save the config
	cfg.ChangeNetwork(config.Network_Holesky)
	err = hdClient.SaveConfig(cfg)
	require.NoError(t, err)
	t.Log("Config saved successfully")

	// Make a Hyperdrive daemon server
	hdSp, err := hdcommon.NewHyperdriveServiceProvider(hdCtx.UserDirPath, hdCtx.NetworksDir)
	require.NoError(t, err)
	hdWg := &sync.WaitGroup{}
	serverAuthMgr := auth.NewAuthorizationManager("", "server", auth.DefaultRequestLifespan)
	serverAuthMgr.SetKey([]byte(hdTestApiKey))
	hdServer, err := hdserver.NewServerManager(hdSp, "localhost", 0, hdWg, serverAuthMgr)
	require.NoError(t, err)
	hdServerPort := hdServer.GetPort()
	defer hdServer.Stop()
	t.Logf("Hyperdrive daemon server started on port %d", hdServerPort)
}
