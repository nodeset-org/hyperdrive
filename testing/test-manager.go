package testing

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/nodeset-org/hyperdrive-daemon/client"
	"github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive-daemon/server"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	nsserver "github.com/nodeset-org/nodeset-client-go/server-mock/server"
	"github.com/nodeset-org/osha"
	"github.com/rocket-pool/node-manager-core/log"
	"github.com/rocket-pool/node-manager-core/node/services"
)

// HyperdriveTestManager provides bootstrapping and a test service provider, useful for testing
type HyperdriveTestManager struct {
	*osha.TestManager

	// The service provider for the test environment
	serviceProvider *common.ServiceProvider

	// The mock for the nodeset.io service
	nodesetMock *nsserver.NodeSetMockServer

	// The Hyperdrive Daemon server
	serverMgr *server.ServerManager

	// The Hyperdrive Daemon client
	apiClient *client.ApiClient

	// Snapshot ID from the baseline - the initial state of the nodeset.io service prior to running any tests
	baselineSnapshotID string

	// Map of which services were captured during a snapshot
	snapshotServiceMap map[string]Service

	// Wait groups for graceful shutdown
	hdWg *sync.WaitGroup
	nsWg *sync.WaitGroup
}

// Creates a new HyperdriveTestManager instance. Requires management of your own nodeset.io server mock.
// `address` is the address to bind the Hyperdrive daemon to.
func NewHyperdriveTestManager(address string, cfg *hdconfig.HyperdriveConfig, resources *hdconfig.HyperdriveResources, nsServer *nsserver.NodeSetMockServer) (*HyperdriveTestManager, error) {
	tm, err := osha.NewTestManager()
	if err != nil {
		return nil, fmt.Errorf("error creating test manager: %w", err)
	}
	return newHyperdriveTestManagerImpl(address, tm, cfg, resources, nsServer, nil)
}

// Creates a new HyperdriveTestManager instance with default test artifacts.
// `hyperdriveAddress` is the address to bind the Hyperdrive daemon to.
// `nodesetAddress` is the address to bind the nodeset.io server to.
func NewHyperdriveTestManagerWithDefaults(hyperdriveAddress string, nodesetAddress string) (*HyperdriveTestManager, error) {
	tm, err := osha.NewTestManager()
	if err != nil {
		return nil, fmt.Errorf("error creating test manager: %w", err)
	}

	// Make the nodeset.io mock server
	nsWg := &sync.WaitGroup{}
	nodesetMock, err := nsserver.NewNodeSetMockServer(tm.GetLogger(), nodesetAddress, 0)
	if err != nil {
		closeTestManager(tm)
		return nil, fmt.Errorf("error creating nodeset mock server: %v", err)
	}
	err = nodesetMock.Start(nsWg)
	if err != nil {
		closeTestManager(tm)
		return nil, fmt.Errorf("error starting nodeset mock server: %v", err)
	}

	// Make a new Hyperdrive config
	testDir := tm.GetTestDir()
	beaconCfg := tm.GetBeaconMockManager().GetConfig()
	resources := GetTestResources(beaconCfg, fmt.Sprintf("http://%s:%d/api/", nodesetAddress, nodesetMock.GetPort()))
	cfg := hdconfig.NewHyperdriveConfigForNetwork(testDir, hdconfig.Network_LocalTest, resources)
	cfg.Network.Value = hdconfig.Network_LocalTest

	// Make the test manager
	return newHyperdriveTestManagerImpl(hyperdriveAddress, tm, cfg, resources, nodesetMock, nsWg)
}

// Implementation for creating a new HyperdriveTestManager
func newHyperdriveTestManagerImpl(address string, tm *osha.TestManager, cfg *hdconfig.HyperdriveConfig, resources *hdconfig.HyperdriveResources, nsServer *nsserver.NodeSetMockServer, nsWaitGroup *sync.WaitGroup) (*HyperdriveTestManager, error) {
	// Make managers
	beaconCfg := tm.GetBeaconMockManager().GetConfig()
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
		closeTestManager(tm)
		return nil, fmt.Errorf("error creating service provider: %v", err)
	}

	// Make sure the data and modules directories exist
	dataDir := cfg.UserDataPath.Value
	moduleDir := filepath.Join(dataDir, hdconfig.ModulesName)
	err = os.MkdirAll(moduleDir, 0755)
	if err != nil {
		closeTestManager(tm)
		return nil, fmt.Errorf("error creating data and modules directories [%s]: %v", moduleDir, err)
	}

	// Create the server
	hdWg := &sync.WaitGroup{}
	serverMgr, err := server.NewServerManager(serviceProvider, address, 0, hdWg)
	if err != nil {
		closeTestManager(tm)
		return nil, fmt.Errorf("error creating hyperdrive server: %v", err)
	}

	// Create the client
	urlString := fmt.Sprintf("http://%s:%d/%s", address, serverMgr.GetPort(), hdconfig.HyperdriveApiClientRoute)
	url, err := url.Parse(urlString)
	if err != nil {
		closeTestManager(tm)
		return nil, fmt.Errorf("error parsing client URL [%s]: %v", urlString, err)
	}
	apiClient := client.NewApiClient(url, tm.GetLogger(), nil)

	// Return
	m := &HyperdriveTestManager{
		TestManager:        tm,
		serviceProvider:    serviceProvider,
		nodesetMock:        nsServer,
		serverMgr:          serverMgr,
		apiClient:          apiClient,
		hdWg:               hdWg,
		nsWg:               nsWaitGroup,
		snapshotServiceMap: map[string]Service{},
	}

	// Create the baseline snapshot
	baselineSnapshotID, err := m.takeSnapshot(Service_All)
	if err != nil {
		return nil, fmt.Errorf("error creating baseline snapshot: %w", err)
	}
	m.baselineSnapshotID = baselineSnapshotID

	return m, nil
}

// Closes the Hyperdrive test manager, shutting down the daemon
func (m *HyperdriveTestManager) Close() error {
	if m.nodesetMock != nil {
		// Check if we're managing the service - if so just stop it
		if m.nsWg != nil {
			err := m.nodesetMock.Stop()
			if err != nil {
				m.GetLogger().Warn("WARNING: nodeset server mock didn't shutdown cleanly", log.Err(err))
			}
			m.nsWg.Wait()
			m.TestManager.GetLogger().Info("Stopped nodeset.io mock server")
		} else {
			err := m.nodesetMock.GetManager().RevertToSnapshot(m.baselineSnapshotID)
			if err != nil {
				m.GetLogger().Warn("WARNING: error reverting nodeset server mock to baseline", log.Err(err))
			} else {
				m.TestManager.GetLogger().Info("Reverted nodeset.io mock server to baseline snapshot")
			}
		}
		m.nodesetMock = nil
	}
	if m.serverMgr != nil {
		m.serverMgr.Stop()
		m.hdWg.Wait()
		m.TestManager.GetLogger().Info("Stopped daemon API server")
		m.serverMgr = nil
	}
	if m.TestManager != nil {
		err := m.TestManager.Close()
		m.TestManager = nil
		return err
	}
	return nil
}

// ===============
// === Getters ===
// ===============

// Returns the service provider for the test environment
func (m *HyperdriveTestManager) GetServiceProvider() *common.ServiceProvider {
	return m.serviceProvider
}

// Get the nodeset.io mock server
func (m *HyperdriveTestManager) GetNodeSetMockServer() *nsserver.NodeSetMockServer {
	return m.nodesetMock
}

// Returns the Hyperdrive Daemon server
func (m *HyperdriveTestManager) GetServerManager() *server.ServerManager {
	return m.serverMgr
}

// Returns the Hyperdrive Daemon client
func (m *HyperdriveTestManager) GetApiClient() *client.ApiClient {
	return m.apiClient
}

// ====================
// === Snapshotting ===
// ====================

// Reverts the services to the baseline snapshot
func (m *HyperdriveTestManager) RevertToBaseline() error {
	err := m.TestManager.RevertToBaseline()
	if err != nil {
		return fmt.Errorf("error reverting to baseline snapshot: %w", err)
	}

	// Regenerate the baseline snapshot since Hardhat can't revert to it multiple times
	baselineSnapshotID, err := m.takeSnapshot(Service_All)
	if err != nil {
		return fmt.Errorf("error creating baseline snapshot: %w", err)
	}
	m.baselineSnapshotID = baselineSnapshotID
	return nil
}

// Takes a snapshot of the service states
func (m *HyperdriveTestManager) CreateCustomSnapshot(services Service) (string, error) {
	return m.takeSnapshot(services)
}

// Revert the services to a snapshot state
func (m *HyperdriveTestManager) RevertToCustomSnapshot(snapshotID string) error {
	return m.revertToSnapshot(snapshotID)
}

// ==========================
// === Internal Functions ===
// ==========================

// Takes a snapshot of the service states
func (m *HyperdriveTestManager) takeSnapshot(services Service) (string, error) {
	// Run the parent snapshotter
	parentServices := osha.Service(services)
	snapshotName, err := m.TestManager.CreateCustomSnapshot(parentServices)
	if err != nil {
		return "", fmt.Errorf("error taking snapshot: %w", err)
	}

	// Snapshot the nodeset.io mock
	if services.Contains(Service_NodeSet) {
		m.nodesetMock.GetManager().TakeSnapshot(snapshotName)
	}

	// Store the services that were captured
	m.snapshotServiceMap[snapshotName] = services
	return snapshotName, nil
}

// Revert the services to a snapshot state
func (m *HyperdriveTestManager) revertToSnapshot(snapshotID string) error {
	services, exists := m.snapshotServiceMap[snapshotID]
	if !exists {
		return fmt.Errorf("snapshot with ID [%s] does not exist", snapshotID)
	}

	// Revert the nodeset.io mock
	if services.Contains(Service_NodeSet) {
		err := m.nodesetMock.GetManager().RevertToSnapshot(snapshotID)
		if err != nil {
			return fmt.Errorf("error reverting the nodeset.io mock to snapshot %s: %w", snapshotID, err)
		}
	}

	return m.TestManager.RevertToCustomSnapshot(snapshotID)
}

// Closes the OSHA test manager, logging any errors
func closeTestManager(tm *osha.TestManager) {
	err := tm.Close()
	if err != nil {
		tm.GetLogger().Error("Error closing test manager", log.Err(err))
	}
}
