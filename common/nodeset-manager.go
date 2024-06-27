package common

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive-daemon/module-utils/services"
	"github.com/nodeset-org/hyperdrive-daemon/nodeset/api_v2"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	"github.com/nodeset-org/hyperdrive-stakewise/shared/keys"
	"github.com/rocket-pool/node-manager-core/log"
	"github.com/rocket-pool/node-manager-core/node/wallet"
)

// NodeSetServiceManager is a manager for interactions with the NodeSet service
type NodeSetServiceManager struct {
	// The node wallet
	wallet *wallet.Wallet

	// Resources for the current network
	resources *hdconfig.HyperdriveResources

	// Client for the v2 API
	v2Client *api_v2.NodeSetClient

	// The current session token
	sessionToken string

	// The node wallet's registration status
	nodeRegistrationStatus api.NodeSetRegistrationStatus

	// Mutex for the registration status
	lock *sync.Mutex
}

// Creates a new NodeSet service manager
func NewNodeSetServiceManager(sp *ServiceProvider) *NodeSetServiceManager {
	wallet := sp.GetWallet()
	resources := sp.GetResources()

	return &NodeSetServiceManager{
		wallet:                 wallet,
		resources:              resources,
		v2Client:               api_v2.NewNodeSetClient(resources, hdconfig.ClientTimeout),
		nodeRegistrationStatus: api.NodeSetRegistrationStatus_Unknown,
		lock:                   &sync.Mutex{},
	}
}

// Get the registration status of the node
func (m *NodeSetServiceManager) GetRegistrationStatus() api.NodeSetRegistrationStatus {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.nodeRegistrationStatus
}

// Log in to the NodeSet server
func (m *NodeSetServiceManager) Login(ctx context.Context) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.loginImpl(ctx)
}

// Register the node with the NodeSet server
func (m *NodeSetServiceManager) RegisterNode(ctx context.Context, email string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}
	logger.Debug("Registering node with NodeSet")

	// Make sure there's a wallet
	walletStatus, err := m.wallet.GetStatus()
	if err != nil {
		return fmt.Errorf("error getting wallet status: %w", err)
	}
	if !walletStatus.Wallet.IsLoaded {
		return fmt.Errorf("can't register node with NodeSet, wallet not loaded")
	}

	// Create the signature
	message := fmt.Sprintf(api_v2.NodeAddressMessageFormat, email, walletStatus.Wallet.WalletAddress.Hex())
	sigBytes, err := m.wallet.SignMessage([]byte(message))
	if err != nil {
		m.setRegistrationStatus(api.NodeSetRegistrationStatus_Unknown)
		return fmt.Errorf("error signing registration message: %w", err)
	}

	// Run the request
	return m.v2Client.NodeAddress(ctx, email, walletStatus.Wallet.WalletAddress, sigBytes)
}

// Get the version of the latest deposit data set from the server
func (m *NodeSetServiceManager) StakeWise_GetServerDepositDataVersion(ctx context.Context, vault common.Address) (int, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}
	logger.Debug("Getting server deposit data version")

	// Run the request
	var data api_v2.DepositDataMetaData
	err := m.runRequest(ctx, func(ctx context.Context) error {
		var err error
		data, err = m.v2Client.DepositDataMeta(ctx, vault, m.resources.EthNetworkName)
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("error getting deposit data version: %w", err)
	}
	return data.Version, nil
}

// Get the deposit data set from the server
func (m *NodeSetServiceManager) StakeWise_GetServerDepositData(ctx context.Context, vault common.Address) (int, []types.ExtendedDepositData, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}
	logger.Debug("Getting deposit data")

	// Run the request
	var data api_v2.DepositDataData
	err := m.runRequest(ctx, func(ctx context.Context) error {
		var err error
		data, err = m.v2Client.DepositData_Get(ctx, vault, m.resources.EthNetworkName)
		return err
	})
	if err != nil {
		return 0, nil, fmt.Errorf("error getting deposit data: %w", err)
	}
	return data.Version, data.DepositData, nil
}

// Uploads local deposit data set to the server
func (m *NodeSetServiceManager) StakeWise_UploadDepositData(ctx context.Context, depositData []*types.ExtendedDepositData) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}
	logger.Debug("Uploading deposit data")

	// Run the request
	err := m.runRequest(ctx, func(ctx context.Context) error {
		return m.v2Client.DepositData_Post(ctx, depositData)
	})
	if err != nil {
		return fmt.Errorf("error uploading deposit data: %w", err)
	}
	return nil
}

// Get the version of the latest deposit data set from the server
func (m *NodeSetServiceManager) StakeWise_GetRegisteredValidators(ctx context.Context, vault common.Address) ([]api.ValidatorStatus, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}
	logger.Debug("Getting registered validators")

	// Run the request
	var data api_v2.ValidatorsData
	err := m.runRequest(ctx, func(ctx context.Context) error {
		var err error
		data, err = m.v2Client.Validators_Get(ctx, m.resources.EthNetworkName)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("error getting registered validators: %w", err)
	}
	return data.Validators, nil
}

// Uploads signed exit messages set to the server
func (m *NodeSetServiceManager) StakeWise_UploadSignedExitMessages(ctx context.Context, exitData []api.ExitData) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}
	logger.Debug("Uploading signed exit messages")

	// Run the request
	err := m.runRequest(ctx, func(ctx context.Context) error {
		return m.v2Client.Validators_Patch(ctx, exitData, m.resources.EthNetworkName)
	})
	if err != nil {
		return fmt.Errorf("error uploading signed exit messages: %w", err)
	}
	return nil
}

// ========================
// === Internal Methods ===
// ========================

// Runs a request to the NodeSet server, re-logging in if necessary
func (m *NodeSetServiceManager) runRequest(ctx context.Context, request func(ctx context.Context) error) error {
	// Run the request
	err := request(ctx)
	if err != nil {
		if errors.Is(err, api_v2.ErrInvalidSession) {
			// Session expired so log in again
			err = m.loginImpl(ctx)
			if err != nil {
				return err
			}

			// Re-run the request
			return request(ctx)
		} else {
			return err
		}
	}
	return nil
}

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
		m.nodeRegistrationStatus = api.NodeSetRegistrationStatus_NoWallet
		return fmt.Errorf("can't log into nodeset, hyperdrive wallet not initialized yet")
	}

	// Log the login attempt
	logger.Info("Not authenticated with the NodeSet server, logging in")

	// Get the nonce
	nonceData, err := m.v2Client.Nonce(ctx)
	if err != nil {
		m.setRegistrationStatus(api.NodeSetRegistrationStatus_Unknown)
		return fmt.Errorf("error getting nonce for login: %w", err)
	}
	logger.Debug("Got nonce for login",
		slog.String(keys.NonceKey, nonceData.Nonce),
	)

	// Create a new session
	m.setSessionToken(nonceData.Token)

	// Create the signature
	message := fmt.Sprintf(api_v2.LoginMessageFormat, nonceData.Nonce, walletStatus.Wallet.WalletAddress.Hex())
	sigBytes, err := m.wallet.SignMessage([]byte(message))
	if err != nil {
		m.setRegistrationStatus(api.NodeSetRegistrationStatus_Unknown)
		return fmt.Errorf("error signing login message: %w", err)
	}

	// Attempt a login
	_, err = m.v2Client.Login(ctx, nonceData.Nonce, walletStatus.Wallet.WalletAddress, sigBytes)
	if err != nil {
		if errors.Is(err, wallet.ErrWalletNotLoaded) {
			m.setRegistrationStatus(api.NodeSetRegistrationStatus_NoWallet)
			return err
		}
		m.setRegistrationStatus(api.NodeSetRegistrationStatus_Unknown)
		return fmt.Errorf("error logging in: %w", err)
	}

	// Success
	m.setSessionToken(nonceData.Token)
	logger.Info("Logged into NodeSet server")
	m.setRegistrationStatus(api.NodeSetRegistrationStatus_Registered)

	return nil
}

// Sets the session token for the client after logging in
func (m *NodeSetServiceManager) setSessionToken(sessionToken string) {
	m.sessionToken = sessionToken
	m.v2Client.SetSessionToken(sessionToken)
}

// Sets the registration status of the node
func (m *NodeSetServiceManager) setRegistrationStatus(status api.NodeSetRegistrationStatus) {
	// Only set to unknown if it hasn't already been figured out
	if status == api.NodeSetRegistrationStatus_Unknown &&
		(m.nodeRegistrationStatus == api.NodeSetRegistrationStatus_Unregistered ||
			m.nodeRegistrationStatus == api.NodeSetRegistrationStatus_Registered) {
		return
	}

	m.nodeRegistrationStatus = status
}
