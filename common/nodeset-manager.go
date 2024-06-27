package common

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/nodeset-org/hyperdrive-daemon/module-utils/services"
	"github.com/nodeset-org/hyperdrive-daemon/nodeset/api_v2"
	"github.com/nodeset-org/hyperdrive-daemon/nodeset/api_v2/types"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/nodeset-org/hyperdrive-stakewise/shared/keys"
	"github.com/rocket-pool/node-manager-core/log"
	"github.com/rocket-pool/node-manager-core/node/wallet"
)

// NodeSetServiceManager is a manager for interactions with the NodeSet service
type NodeSetServiceManager struct {
	// The node wallet
	wallet *wallet.Wallet

	// Client for the v2 API
	v2Client *api_v2.NodeSetClient

	// Active session token
	token string

	// The node wallet's registration status
	nodeRegistrationStatus types.NodeSetRegistrationStatus

	// Mutex for the registration status
	lock *sync.Mutex
}

// Creates a new NodeSet service manager
func NewNodeSetServiceManager(sp *ServiceProvider) *NodeSetServiceManager {
	wallet := sp.GetWallet()
	resources := sp.GetResources()

	client := &http.Client{
		Timeout: hdconfig.ClientTimeout,
	}

	return &NodeSetServiceManager{
		wallet:                 wallet,
		v2Client:               api_v2.NewNodeSetClient(wallet, resources, client),
		nodeRegistrationStatus: types.NodeSetRegistrationStatus_Unknown,
		lock:                   &sync.Mutex{},
	}
}

// Get the registration status of the node
func (m *NodeSetServiceManager) GetRegistrationStatus() types.NodeSetRegistrationStatus {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.nodeRegistrationStatus
}

// Log in to the NodeSet server
func (m *NodeSetServiceManager) Login(ctx context.Context) error {
	return m.loginImpl(ctx)
}

// ========================
// === Internal Methods ===
// ========================

// Implementation for logging in
func (m *NodeSetServiceManager) loginImpl(ctx context.Context) error {
	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}

	// Get the node wallet
	walletStatus, err := m.wallet.GetStatus()
	if err != nil {
		return fmt.Errorf("error getting wallet status for login: %w", err)
	}
	err = services.CheckIfWalletReady(walletStatus)
	if err != nil {
		m.nodeRegistrationStatus = types.NodeSetRegistrationStatus_NoWallet
		return fmt.Errorf("can't log into nodeset, hyperdrive wallet not initialized yet")
	}

	// Log the login attempt
	logger.Info("Not authenticated with the NodeSet server, logging in")

	// Get the nonce
	nonceData, err := m.v2Client.Nonce(ctx)
	if err != nil {
		m.setRegistrationStatus(types.NodeSetRegistrationStatus_Unknown)
		return fmt.Errorf("error getting nonce for login: %w", err)
	}
	logger.Debug("Got nonce for login",
		slog.String(keys.NonceKey, nonceData.Nonce),
	)
	m.token = nonceData.Token // Store this as a temp token for the login request

	// Attempt a login
	loginData, err := m.v2Client.Login(ctx, nonceData.Nonce, walletStatus.Wallet.WalletAddress)
	if err != nil {
		if errors.Is(err, wallet.ErrWalletNotLoaded) {
			m.setRegistrationStatus(types.NodeSetRegistrationStatus_NoWallet)
			return err
		}
		m.setRegistrationStatus(types.NodeSetRegistrationStatus_Unknown)
		return fmt.Errorf("error logging in: %w", err)
	}

	// Success
	m.token = loginData.Token // Save this as the persistent token for all other requests
	logger.Info("Logged into NodeSet server")
	m.setRegistrationStatus(types.NodeSetRegistrationStatus_Registered)

	return nil
}

// Sets the registration status of the node
func (m *NodeSetServiceManager) setRegistrationStatus(status types.NodeSetRegistrationStatus) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Only set to unknown if it hasn't already been figured out
	if status == types.NodeSetRegistrationStatus_Unknown &&
		(m.nodeRegistrationStatus == types.NodeSetRegistrationStatus_Unregistered ||
			m.nodeRegistrationStatus == types.NodeSetRegistrationStatus_Registered) {
		return
	}

	m.nodeRegistrationStatus = status
}
