package common

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	apiv1 "github.com/nodeset-org/nodeset-client-go/api-v1"
	apiv2 "github.com/nodeset-org/nodeset-client-go/api-v2"
	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/rocket-pool/node-manager-core/log"
	"github.com/rocket-pool/node-manager-core/node/wallet"
	"github.com/rocket-pool/node-manager-core/utils"
)

// NodeSetServiceManager is a manager for interactions with the NodeSet service
type NodeSetServiceManager struct {
	// The node wallet
	wallet *wallet.Wallet

	// Resources for the current network
	resources *hdconfig.HyperdriveResources

	// Client for the v1 API
	v1Client *apiv1.NodeSetClient

	// Client for the v2 API
	v2Client *apiv2.NodeSetClient

	// The current session token
	sessionToken string

	// The node wallet's registration status
	nodeRegistrationStatus api.NodeSetRegistrationStatus

	// Mutex for the registration status
	lock *sync.Mutex
}

// Creates a new NodeSet service manager
func NewNodeSetServiceManager(sp IHyperdriveServiceProvider) *NodeSetServiceManager {
	wallet := sp.GetWallet()
	resources := sp.GetResources()

	return &NodeSetServiceManager{
		wallet:                 wallet,
		resources:              resources,
		v1Client:               apiv1.NewNodeSetClient(resources.NodeSetApiUrl, hdconfig.ClientTimeout),
		v2Client:               apiv2.NewNodeSetClient(resources.NodeSetApiUrl, hdconfig.ClientTimeout),
		nodeRegistrationStatus: api.NodeSetRegistrationStatus_Unknown,
		lock:                   &sync.Mutex{},
	}
}

// Get the registration status of the node
func (m *NodeSetServiceManager) GetRegistrationStatus(ctx context.Context) (api.NodeSetRegistrationStatus, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Force refresh the registration status if it hasn't been determined yet
	if m.nodeRegistrationStatus == api.NodeSetRegistrationStatus_Unknown ||
		m.nodeRegistrationStatus == api.NodeSetRegistrationStatus_NoWallet {
		err := m.loginImpl(ctx)
		return m.nodeRegistrationStatus, err
	}
	return m.nodeRegistrationStatus, nil
}

// Log in to the NodeSet server
func (m *NodeSetServiceManager) Login(ctx context.Context) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.loginImpl(ctx)
}

// Result of RegisterNode
type RegistrationResult int

const (
	RegistrationResult_Unknown RegistrationResult = iota
	RegistrationResult_Success
	RegistrationResult_AlreadyRegistered
	RegistrationResult_NotWhitelisted
)

// Register the node with the NodeSet server
func (m *NodeSetServiceManager) RegisterNode(ctx context.Context, email string) (RegistrationResult, error) {
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
		return RegistrationResult_Unknown, fmt.Errorf("error getting wallet status: %w", err)
	}
	if !walletStatus.Wallet.IsLoaded {
		return RegistrationResult_Unknown, fmt.Errorf("can't register node with NodeSet, wallet not loaded")
	}

	// Create the signature
	message := fmt.Sprintf(apiv1.NodeAddressMessageFormat, email, walletStatus.Wallet.WalletAddress.Hex())
	sigBytes, err := m.wallet.SignMessage([]byte(message))
	if err != nil {
		m.setRegistrationStatus(api.NodeSetRegistrationStatus_Unknown)
		return RegistrationResult_Unknown, fmt.Errorf("error signing registration message: %w", err)
	}

	// Run the request
	err = m.v1Client.NodeAddress(ctx, email, walletStatus.Wallet.WalletAddress, sigBytes)
	if err != nil {
		m.setRegistrationStatus(api.NodeSetRegistrationStatus_Unknown)
		if errors.Is(err, apiv1.ErrAlreadyRegistered) {
			return RegistrationResult_AlreadyRegistered, nil
		} else if errors.Is(err, apiv1.ErrNotWhitelisted) {
			return RegistrationResult_NotWhitelisted, nil
		}
		return RegistrationResult_Unknown, fmt.Errorf("error registering node: %w", err)
	}
	return RegistrationResult_Success, nil
}

// =========================
// === StakeWise Methods ===
// =========================

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
	var data apiv1.DepositDataMetaData
	err := m.runRequest(ctx, func(ctx context.Context) error {
		var err error
		data, err = m.v1Client.DepositDataMeta(ctx, vault, m.resources.EthNetworkName)
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("error getting deposit data version: %w", err)
	}
	return data.Version, nil
}

// Get the deposit data set from the server
func (m *NodeSetServiceManager) StakeWise_GetServerDepositData(ctx context.Context, vault common.Address) (int, []beacon.ExtendedDepositData, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}
	logger.Debug("Getting deposit data")

	// Run the request
	var data apiv1.DepositDataData
	err := m.runRequest(ctx, func(ctx context.Context) error {
		var err error
		data, err = m.v1Client.DepositData_Get(ctx, vault, m.resources.EthNetworkName)
		return err
	})
	if err != nil {
		return 0, nil, fmt.Errorf("error getting deposit data: %w", err)
	}
	return data.Version, data.DepositData, nil
}

// Uploads local deposit data set to the server
func (m *NodeSetServiceManager) StakeWise_UploadDepositData(ctx context.Context, depositData []beacon.ExtendedDepositData) error {
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
		return m.v1Client.DepositData_Post(ctx, depositData)
	})
	if err != nil {
		return fmt.Errorf("error uploading deposit data: %w", err)
	}
	return nil
}

// Get the version of the latest deposit data set from the server
func (m *NodeSetServiceManager) StakeWise_GetRegisteredValidators(ctx context.Context, vault common.Address) ([]apiv1.ValidatorStatus, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}
	logger.Debug("Getting registered validators")

	// Run the request
	var data apiv1.ValidatorsData
	err := m.runRequest(ctx, func(ctx context.Context) error {
		var err error
		data, err = m.v1Client.Validators_Get(ctx, m.resources.EthNetworkName)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("error getting registered validators: %w", err)
	}
	return data.Validators, nil
}

// Uploads signed exit messages set to the server
func (m *NodeSetServiceManager) StakeWise_UploadSignedExitMessages(ctx context.Context, exitData []apiv1.ExitData) error {
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
		return m.v1Client.Validators_Patch(ctx, exitData, m.resources.EthNetworkName)
	})
	if err != nil {
		return fmt.Errorf("error uploading signed exit messages: %w", err)
	}
	return nil
}

// =============================
// === Constellation Methods ===
// =============================

// Gets a signature for registering / whitelisting the node with the Constellation contracts and the timestamp it was created
func (m *NodeSetServiceManager) Constellation_GetRegistrationSignatureAndTime(ctx context.Context) (time.Time, []byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}
	logger.Debug("Registering with the Constellation contracts")

	// Run the request
	var data apiv2.WhitelistData
	err := m.runRequest(ctx, func(ctx context.Context) error {
		var err error
		data, err = m.v2Client.Whitelist(ctx, big.NewInt(int64(m.resources.ChainID)))
		return err
	})
	if err != nil {
		return time.Time{}, nil, fmt.Errorf("error registering with Constellation: %w", err)
	}

	// Decode the signature
	sig, err := utils.DecodeHex(data.Signature)
	if err != nil {
		return time.Time{}, nil, fmt.Errorf("error decoding signature from server: %w", err)
	}

	// Get the time from the timestamp
	timestamp := time.Unix(data.Time, 0)
	return timestamp, sig, nil
}

// Gets the available minipool count for the node from the Constellation contracts
func (m *NodeSetServiceManager) Constellation_GetAvailableMinipoolCount(ctx context.Context) (int, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}
	logger.Debug("Getting available minipool count")

	// Run the request
	var data apiv2.MinipoolAvailableData
	err := m.runRequest(ctx, func(ctx context.Context) error {
		var err error
		data, err = m.v2Client.MinipoolAvailable(ctx)
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("error getting available minipool count: %w", err)
	}
	return data.Count, nil
}

// Gets the deposit signature for a minipool from the Constellation contracts and the timestamp it was created
func (m *NodeSetServiceManager) Constellation_GetDepositSignatureAndTime(ctx context.Context, minipoolAddress common.Address, salt []byte) (time.Time, []byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}
	logger.Debug("Getting minipool deposit signature")
	// Run the request
	var data apiv2.MinipoolDepositSignatureData
	err := m.runRequest(ctx, func(ctx context.Context) error {
		var err error
		data, err = m.v2Client.MinipoolDepositSignature(ctx, minipoolAddress, salt, big.NewInt(int64(m.resources.ChainID)))
		return err
	})
	if err != nil {
		return time.Time{}, nil, fmt.Errorf("error getting deposit signature: %w", err)
	}

	// Decode the signature
	sig, err := utils.DecodeHex(data.Signature)
	if err != nil {
		return time.Time{}, nil, fmt.Errorf("error decoding signature from server: %w", err)
	}

	// Get the time from the timestamp
	timestamp := time.Unix(data.Time, 0)
	return timestamp, sig, nil
}

// ========================
// === Internal Methods ===
// ========================

// Runs a request to the NodeSet server, re-logging in if necessary
func (m *NodeSetServiceManager) runRequest(ctx context.Context, request func(ctx context.Context) error) error {
	// Run the request
	err := request(ctx)
	if err != nil {
		if errors.Is(err, apiv1.ErrInvalidSession) {
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
	err = CheckIfWalletReady(walletStatus)
	if err != nil {
		m.nodeRegistrationStatus = api.NodeSetRegistrationStatus_NoWallet
		return fmt.Errorf("can't log into nodeset, hyperdrive wallet not initialized yet")
	}

	// Log the login attempt
	logger.Info("Not authenticated with the NodeSet server, logging in")

	// Get the nonce
	nonceData, err := m.v1Client.Nonce(ctx)
	if err != nil {
		m.setRegistrationStatus(api.NodeSetRegistrationStatus_Unknown)
		return fmt.Errorf("error getting nonce for login: %w", err)
	}
	logger.Debug("Got nonce for login",
		slog.String("nonce", nonceData.Nonce),
	)

	// Create a new session
	m.setSessionToken(nonceData.Token)

	// Create the signature
	message := fmt.Sprintf(apiv1.LoginMessageFormat, nonceData.Nonce, walletStatus.Wallet.WalletAddress.Hex())
	sigBytes, err := m.wallet.SignMessage([]byte(message))
	if err != nil {
		m.setRegistrationStatus(api.NodeSetRegistrationStatus_Unknown)
		return fmt.Errorf("error signing login message: %w", err)
	}

	// Attempt a login
	loginData, err := m.v1Client.Login(ctx, nonceData.Nonce, walletStatus.Wallet.WalletAddress, sigBytes)
	if err != nil {
		if errors.Is(err, wallet.ErrWalletNotLoaded) {
			m.setRegistrationStatus(api.NodeSetRegistrationStatus_NoWallet)
			return err
		}
		if errors.Is(err, apiv1.ErrUnregisteredNode) {
			m.setRegistrationStatus(api.NodeSetRegistrationStatus_Unregistered)
			return nil
		}
		m.setRegistrationStatus(api.NodeSetRegistrationStatus_Unknown)
		return fmt.Errorf("error logging in: %w", err)
	}

	// Success
	m.setSessionToken(loginData.Token)
	logger.Info("Logged into NodeSet server")
	m.setRegistrationStatus(api.NodeSetRegistrationStatus_Registered)

	return nil
}

// Sets the session token for the client after logging in
func (m *NodeSetServiceManager) setSessionToken(sessionToken string) {
	m.sessionToken = sessionToken
	m.v1Client.SetSessionToken(sessionToken)
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
