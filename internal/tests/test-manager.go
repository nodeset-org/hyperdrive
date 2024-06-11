package tests

import (
	"fmt"
	"time"

	"github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/nodeset-org/osha"
	"github.com/rocket-pool/node-manager-core/node/services"
)

const (
	// The environment variable for the locally running Hardhat instance
	HardhatEnvVar string = "HARDHAT_URL"
)

// TestManager provides bootstrapping and a test service provider, useful for testing
type TestManager struct {
	*osha.TestManager

	// The service provider for the test environment
	ServiceProvider *common.ServiceProvider
}

// Creates a new TestManager instance
func NewTestManager() (*TestManager, error) {
	tm, err := osha.NewTestManager()
	if err != nil {
		return nil, fmt.Errorf("error creating test manager: %w", err)
	}

	// Make a new Hyperdrive config
	testDir := tm.GetTestDir()
	cfg := config.NewHyperdriveConfig(testDir)
	cfg.Network.Value = config.Network_LocalTest

	// Make test resources
	beaconCfg := tm.GetBeaconMockManager().GetConfig()
	resources := GetTestResources(beaconCfg)

	// Make managers
	ecManager := services.NewExecutionClientManager(tm.GetExecutionClient(), uint(beaconCfg.ChainID), time.Minute)
	bnManager := services.NewBeaconClientManager(tm.GetBeaconClient(), uint(beaconCfg.ChainID), time.Minute)

	// Make a new service provider
	serviceProvider, err := common.NewServiceProviderFromCustomServices(
		cfg,
		resources,
		ecManager,
		bnManager,
		tm.GetDockerMockManager(),
	)
	if err != nil {
		err2 := tm.Close()
		if err2 != nil {
			tm.GetLogger().Error("Error closing test manager: %v", err2)
		}
		return nil, fmt.Errorf("error creating service provider: %v", err)
	}

	// Return
	m := &TestManager{
		TestManager:     tm,
		ServiceProvider: serviceProvider,
	}
	return m, nil
}
