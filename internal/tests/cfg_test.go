package tests

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime/debug"
	"sync"
	"testing"

	hdcommon "github.com/nodeset-org/hyperdrive-daemon/common"
	modservices "github.com/nodeset-org/hyperdrive-daemon/module-utils/services"
	hdserver "github.com/nodeset-org/hyperdrive-daemon/server"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	swcommon "github.com/nodeset-org/hyperdrive-stakewise/common"
	swconfig "github.com/nodeset-org/hyperdrive-stakewise/shared/config"
	hdclient "github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/context"
	"github.com/nodeset-org/osha"
	"github.com/rocket-pool/node-manager-core/config"
	"github.com/stretchr/testify/require"
)

func TestNewConfig_Holesky(t *testing.T) {
	// Take a snapshot, revert at the end
	snapshotName, err := testMgr.CreateCustomSnapshot(osha.Service_Filesystem)
	if err != nil {
		fail("Error creating custom snapshot: %v", err)
	}
	defer cfg_cleanup(snapshotName)

	cfgPath := testMgr.GetTestDir()

	// Make a new Hyperdrive context
	hdCtx := context.HyperdriveContext{
		ConfigPath: cfgPath,
		InstallationInfo: &context.InstallationInfo{
			ScriptsDir:        "../../install/deploy/scripts",
			TemplatesDir:      "../../install/deploy/templates",
			OverrideSourceDir: "../../install/deploy/override",
			NetworksDir:       "../../install/deploy/networks",
		},
	}

	// Load the network settings
	hdNetworkSettings, err := hdconfig.LoadSettingsFiles(hdCtx.NetworksDir)
	require.NoError(t, err)
	swNetSettingsDir := filepath.Join(hdCtx.NetworksDir, hdconfig.ModulesName, swconfig.ModuleName)
	swNetworkSettings, err := swconfig.LoadSettingsFiles(swNetSettingsDir)
	require.NoError(t, err)
	hdCtx.HyperdriveNetworkSettings = hdNetworkSettings
	hdCtx.StakeWiseNetworkSettings = swNetworkSettings

	// Make a new Hyperdrive client
	hdClient, err := hdclient.NewHyperdriveClientFromHyperdriveCtx(&hdCtx)
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
	hdSp, err := hdcommon.NewHyperdriveServiceProvider(hdCtx.ConfigPath, hdCtx.NetworksDir)
	require.NoError(t, err)
	hdWg := &sync.WaitGroup{}
	hdServer, err := hdserver.NewServerManager(hdSp, "localhost", 0, hdWg)
	require.NoError(t, err)
	hdServerPort := hdServer.GetPort()
	defer hdServer.Stop()
	t.Logf("Hyperdrive daemon server started on port %d", hdServerPort)

	// Make a new StakeWise daemon server
	hdApiUrl, _ := url.Parse(fmt.Sprintf("http://localhost:%d", hdServerPort))
	swModDir := filepath.Join(cfg.Hyperdrive.UserDataPath.Value, hdconfig.ModulesName, swconfig.ModuleName)
	err = os.MkdirAll(swModDir, 0755)
	require.NoError(t, err)
	swSettings, err := hdclient.LoadStakeWiseSettings(hdCtx.NetworksDir)
	require.NoError(t, err)
	modSp, err := modservices.NewModuleServiceProvider(hdApiUrl, swModDir, swconfig.ModuleName, swconfig.ClientLogName, func(hdCfg *hdconfig.HyperdriveConfig) (*swconfig.StakeWiseConfig, error) {
		return swconfig.NewStakeWiseConfig(hdCfg, swSettings)
	}, hdconfig.ClientTimeout)
	require.NoError(t, err)
	swSp, err := swcommon.NewStakeWiseServiceProvider(modSp, swSettings)
	require.NoError(t, err)
	t.Log("StakeWise service provider created")

	expectedVaultAddress := swconfig.HoleskyResourcesReference.Vault
	swRes := swSp.GetResources()
	require.Equal(t, expectedVaultAddress, swRes.Vault)
	t.Logf("StakeWise vault address was correct: %s", swRes.Vault.Hex())

	expectedForkVersion := config.HoleskyResourcesReference.GenesisForkVersion
	require.Equal(t, expectedForkVersion, swRes.MergedResources.GenesisForkVersion)
	t.Logf("Genesis fork version was correct: %x", swRes.MergedResources.GenesisForkVersion)
}

func cfg_cleanup(snapshotName string) {
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
