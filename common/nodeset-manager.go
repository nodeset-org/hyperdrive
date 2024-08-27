package common

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	apiv2 "github.com/nodeset-org/nodeset-client-go/api-v2"
	v2constellation "github.com/nodeset-org/nodeset-client-go/api-v2/constellation"
	nscommon "github.com/nodeset-org/nodeset-client-go/common"
	"github.com/nodeset-org/nodeset-client-go/common/core"
	"github.com/nodeset-org/nodeset-client-go/common/stakewise"
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
	resources *hdconfig.MergedResources

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
	message := fmt.Sprintf(core.NodeAddressMessageFormat, email, walletStatus.Wallet.WalletAddress.Hex())
	sigBytes, err := m.wallet.SignMessage([]byte(message))
	if err != nil {
		m.setRegistrationStatus(api.NodeSetRegistrationStatus_Unknown)
		return RegistrationResult_Unknown, fmt.Errorf("error signing registration message: %w", err)
	}

	// Run the request
	err = m.v2Client.Core.NodeAddress(ctx, email, walletStatus.Wallet.WalletAddress, sigBytes)
	if err != nil {
		m.setRegistrationStatus(api.NodeSetRegistrationStatus_Unknown)
		if errors.Is(err, core.ErrAlreadyRegistered) {
			return RegistrationResult_AlreadyRegistered, nil
		} else if errors.Is(err, core.ErrNotWhitelisted) {
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
	var data stakewise.DepositDataMetaData
	err := m.runRequest(ctx, func(ctx context.Context) error {
		var err error
		data, err = m.v2Client.StakeWise.DepositDataMeta(ctx, m.resources.DeploymentName, vault)
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
	var data stakewise.DepositDataData
	err := m.runRequest(ctx, func(ctx context.Context) error {
		var err error
		data, err = m.v2Client.StakeWise.DepositData_Get(ctx, m.resources.DeploymentName, vault)
		return err
	})
	if err != nil {
		return 0, nil, fmt.Errorf("error getting deposit data: %w", err)
	}
	return data.Version, data.DepositData, nil
}

// Uploads local deposit data set to the server
func (m *NodeSetServiceManager) StakeWise_UploadDepositData(ctx context.Context, vault common.Address, depositData []beacon.ExtendedDepositData) error {
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
		return m.v2Client.StakeWise.DepositData_Post(ctx, m.resources.DeploymentName, vault, depositData)
	})
	if err != nil {
		return fmt.Errorf("error uploading deposit data: %w", err)
	}
	return nil
}

// Get the version of the latest deposit data set from the server
func (m *NodeSetServiceManager) StakeWise_GetRegisteredValidators(ctx context.Context, vault common.Address) ([]stakewise.ValidatorStatus, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}
	logger.Debug("Getting registered validators")

	// Run the request
	var data stakewise.ValidatorsData
	err := m.runRequest(ctx, func(ctx context.Context) error {
		var err error
		data, err = m.v2Client.StakeWise.Validators_Get(ctx, m.resources.DeploymentName, vault)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("error getting registered validators: %w", err)
	}
	return data.Validators, nil
}

// Uploads signed exit messages set to the server
func (m *NodeSetServiceManager) StakeWise_UploadSignedExitMessages(ctx context.Context, vault common.Address, exitData []nscommon.ExitData) error {
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
		return m.v2Client.StakeWise.Validators_Patch(ctx, m.resources.DeploymentName, vault, exitData)
	})
	if err != nil {
		return fmt.Errorf("error uploading signed exit messages: %w", err)
	}
	return nil
}

// =============================
// === Constellation Methods ===
// =============================

// Gets a signature for registering / whitelisting the node with the Constellation contracts
func (m *NodeSetServiceManager) Constellation_GetRegistrationSignature(ctx context.Context) ([]byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}
	logger.Debug("Registering with the Constellation contracts")

	// Run the request
	var data v2constellation.WhitelistData
	err := m.runRequest(ctx, func(ctx context.Context) error {
		var err error
		data, err = m.v2Client.Constellation.Whitelist(ctx, m.resources.HyperdriveResources.DeploymentName)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("error registering with Constellation: %w", err)
	}

	// Decode the signature
	sig, err := utils.DecodeHex(data.Signature)
	if err != nil {
		return nil, fmt.Errorf("error decoding signature from server: %w", err)
	}
	return sig, nil
}

// Gets the deposit signature for a minipool from the Constellation contracts
func (m *NodeSetServiceManager) Constellation_GetDepositSignature(ctx context.Context, minipoolAddress common.Address, salt *big.Int) ([]byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}

	// Run the request
	var data v2constellation.MinipoolDepositSignatureData
	logger.Debug("Getting minipool deposit signature")
	err := m.runRequest(ctx, func(ctx context.Context) error {
		var err error
		data, err = m.v2Client.Constellation.MinipoolDepositSignature(ctx, m.resources.HyperdriveResources.DeploymentName, minipoolAddress, salt)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("error getting deposit signature: %w", err)
	}

	// Decode the signature
	sig, err := utils.DecodeHex(data.Signature)
	if err != nil {
		return nil, fmt.Errorf("error decoding signature from server: %w", err)
	}
	return sig, nil
}

// Get the validators that NodeSet has on record for this node
func (m *NodeSetServiceManager) Constellation_GetValidators(ctx context.Context) ([]v2constellation.ValidatorStatus, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}

	// Run the request
	var data v2constellation.ValidatorsData
	logger.Debug("Getting validators for node")
	err := m.runRequest(ctx, func(ctx context.Context) error {
		var err error
		data, err = m.v2Client.Constellation.Validators_Get(ctx, m.resources.HyperdriveResources.DeploymentName)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("error getting validators for node: %w", err)
	}
	return data.Validators, nil
}

// Upload signed exit messages for Constellation minipools to the NodeSet service
func (m *NodeSetServiceManager) Constellation_UploadSignedExitMessages(ctx context.Context, exitMessages []nscommon.ExitData) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}

	// Run the request
	logger.Debug("Submitting signed exit messages to nodeset")
	err := m.runRequest(ctx, func(ctx context.Context) error {
		return m.v2Client.Constellation.Validators_Patch(ctx, m.resources.HyperdriveResources.DeploymentName, exitMessages)
	})
	if err != nil {
		return fmt.Errorf("error submitting signed exit messages: %w", err)
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
		if errors.Is(err, nscommon.ErrInvalidSession) {
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
	nonceData, err := m.v2Client.Core.Nonce(ctx)
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
	message := fmt.Sprintf(core.LoginMessageFormat, nonceData.Nonce, walletStatus.Wallet.WalletAddress.Hex())
	sigBytes, err := m.wallet.SignMessage([]byte(message))
	if err != nil {
		m.setRegistrationStatus(api.NodeSetRegistrationStatus_Unknown)
		return fmt.Errorf("error signing login message: %w", err)
	}

	// Attempt a login
	loginData, err := m.v2Client.Core.Login(ctx, nonceData.Nonce, walletStatus.Wallet.WalletAddress, sigBytes)
	if err != nil {
		if errors.Is(err, wallet.ErrWalletNotLoaded) {
			m.setRegistrationStatus(api.NodeSetRegistrationStatus_NoWallet)
			return err
		}
		if errors.Is(err, core.ErrUnregisteredNode) {
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
