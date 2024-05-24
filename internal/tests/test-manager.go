package tests

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nodeset-org/beacon-mock/db"
	"github.com/nodeset-org/beacon-mock/manager"
	"github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive-daemon/internal/docker"
	"github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/rocket-pool/node-manager-core/beacon/client"
	"github.com/rocket-pool/node-manager-core/node/services"
)

const (
	// The environment variable for the locally running Hardhat instance
	HardhatEnvVar string = "HARDHAT_URL"
)

// TestManager provides bootstrapping and a test service provider, useful for testing
type TestManager struct {
	// The service provider for the test environment
	ServiceProvider *common.ServiceProvider

	// Logger for logging output messages during tests
	Logger *slog.Logger

	// Internal fields
	testingConfigDir string
}

// Creates a new TestManager instance
func NewTestManager() *TestManager {
	m := &TestManager{}

	// Make sure the Hardhat URL
	hardhatUrl, exists := os.LookupEnv(HardhatEnvVar)
	if !exists {
		m.Fail("Hardhat URL env var [%s] not set", HardhatEnvVar)
	}

	// Make a new logger
	m.Logger = slog.Default()

	// Create a temp folder
	var err error
	m.testingConfigDir, err = os.MkdirTemp("", "hd-tests-*")
	if err != nil {
		m.Fail("Error creating temp config dir: %v", err)
	}

	// Create the default Beacon config
	beaconCfg := db.NewDefaultConfig()

	// Make a new Hyperdrive config
	cfg := config.NewHyperdriveConfig(m.testingConfigDir)
	cfg.Network.Value = config.Network_LocalTest

	// Make test resources
	resources := GetTestResources(beaconCfg)

	// Make the Execution client manager
	clientTimeout := time.Duration(10) * time.Second
	primaryEc, err := ethclient.Dial(hardhatUrl)
	if err != nil {
		m.Fail("Error creating primary eth client with URL [%s]: %v", hardhatUrl, err)
	}
	ecManager, err := services.NewExecutionClientManager(primaryEc, nil, uint(beaconCfg.ChainID), clientTimeout)
	if err != nil {
		m.Fail("Error creating execution client manager: %v", err)
	}

	// Make the Beacon client manager
	primaryBnProvider := manager.NewBeaconMockManager(m.Logger, beaconCfg)
	primaryBn := client.NewStandardClient(primaryBnProvider)
	bnManager, err := services.NewBeaconClientManager(primaryBn, nil, uint(beaconCfg.ChainID), clientTimeout)
	if err != nil {
		m.Fail("Error creating beacon client manager: %v", err)
	}

	// Make a Docker client mock
	docker := docker.NewDockerClientMock()

	// Make a new service provider
	serviceProvider, err := common.NewServiceProviderFromCustomServices(
		cfg,
		resources,
		ecManager,
		bnManager,
		docker,
	)
	if err != nil {
		m.Fail("Error creating service provider: %v", err)
	}
	m.ServiceProvider = serviceProvider

	// Return
	return m
}

// Prints an error message to stderr and exits the program with an error code
func (m *TestManager) Fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	m.Cleanup()
	os.Exit(1)
}

// Cleans up the test environment, including the temporary folder to house any generated files
func (m *TestManager) Cleanup() {
	// Remove the temp folder
	if m.testingConfigDir == "" {
		return
	}
	err := os.RemoveAll(m.testingConfigDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "error removing temp config dir [%s]: %v", m.testingConfigDir, err)
	}
	m.testingConfigDir = ""
}
