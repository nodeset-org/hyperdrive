package services

import (
	"context"
	"errors"
	"fmt"
	"net"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	beaconutils "github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/hyperdrive/daemon-utils/beacon"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/nodeset-org/hyperdrive/shared/utils/log"
)

// This is a proxy for multiple Beacon clients, providing natural fallback support if one of them fails.
type BeaconClientManager struct {
	primaryBc       types.IBeaconClient
	fallbackBc      types.IBeaconClient
	logger          log.ColorLogger
	primaryReady    bool
	fallbackReady   bool
	ignoreSyncCheck bool
}

// This is a signature for a wrapped Beacon client function that only returns an error
type bcFunction0 func(types.IBeaconClient) error

// This is a signature for a wrapped Beacon client function that returns 1 var and an error
type bcFunction1 func(types.IBeaconClient) (interface{}, error)

// This is a signature for a wrapped Beacon client function that returns 2 vars and an error
type bcFunction2 func(types.IBeaconClient) (interface{}, interface{}, error)

// Creates a new BeaconClientManager instance based on the Hyperdrive config
func NewBeaconClientManager(cfg *config.HyperdriveConfig) (*BeaconClientManager, error) {
	// Primary BN
	var primaryProvider string
	if cfg.IsLocalMode() {
		primaryProvider = fmt.Sprintf("http://%s:%d", config.ContainerID_BeaconNode, cfg.LocalBeaconConfig.HttpPort.Value)
	} else if cfg.ClientMode.Value == config.ClientMode_External {
		primaryProvider = cfg.ExternalBeaconConfig.HttpUrl.Value
	} else {
		return nil, fmt.Errorf("unknown client mode '%v'", cfg.ClientMode.Value)
	}

	// Fallback BN
	var fallbackProvider string
	if cfg.Fallback.UseFallbackClients.Value {
		fallbackProvider = cfg.Fallback.BnHttpUrl.Value
	}

	var primaryBc types.IBeaconClient
	var fallbackBc types.IBeaconClient
	primaryBc = beacon.NewStandardHttpClient(primaryProvider, config.ClientTimeout)
	if fallbackProvider != "" {
		fallbackBc = beacon.NewStandardHttpClient(fallbackProvider, config.ClientTimeout)
	}

	return &BeaconClientManager{
		primaryBc:     primaryBc,
		fallbackBc:    fallbackBc,
		logger:        log.NewColorLogger(color.FgHiBlue),
		primaryReady:  true,
		fallbackReady: fallbackBc != nil,
	}, nil
}

func (m *BeaconClientManager) IsPrimaryReady() bool {
	return m.primaryReady
}

func (m *BeaconClientManager) IsFallbackReady() bool {
	return m.fallbackReady
}

/// ======================
/// BeaconClient Functions
/// ======================

// Get the client's sync status
func (m *BeaconClientManager) GetSyncStatus(ctx context.Context) (types.SyncStatus, error) {
	result, err := m.runFunction1(func(client types.IBeaconClient) (interface{}, error) {
		return client.GetSyncStatus(ctx)
	})
	if err != nil {
		return types.SyncStatus{}, err
	}
	return result.(types.SyncStatus), nil
}

// Get the Beacon configuration
func (m *BeaconClientManager) GetEth2Config(ctx context.Context) (types.Eth2Config, error) {
	result, err := m.runFunction1(func(client types.IBeaconClient) (interface{}, error) {
		return client.GetEth2Config(ctx)
	})
	if err != nil {
		return types.Eth2Config{}, err
	}
	return result.(types.Eth2Config), nil
}

// Get the Beacon configuration
func (m *BeaconClientManager) GetEth2DepositContract(ctx context.Context) (types.Eth2DepositContract, error) {
	result, err := m.runFunction1(func(client types.IBeaconClient) (interface{}, error) {
		return client.GetEth2DepositContract(ctx)
	})
	if err != nil {
		return types.Eth2DepositContract{}, err
	}
	return result.(types.Eth2DepositContract), nil
}

// Get the attestations in a Beacon chain block
func (m *BeaconClientManager) GetAttestations(ctx context.Context, blockId string) ([]types.AttestationInfo, bool, error) {
	result1, result2, err := m.runFunction2(func(client types.IBeaconClient) (interface{}, interface{}, error) {
		return client.GetAttestations(ctx, blockId)
	})
	if err != nil {
		return nil, false, err
	}
	return result1.([]types.AttestationInfo), result2.(bool), nil
}

// Get a Beacon chain block
func (m *BeaconClientManager) GetBeaconBlock(ctx context.Context, blockId string) (types.BeaconBlock, bool, error) {
	result1, result2, err := m.runFunction2(func(client types.IBeaconClient) (interface{}, interface{}, error) {
		return client.GetBeaconBlock(ctx, blockId)
	})
	if err != nil {
		return types.BeaconBlock{}, false, err
	}
	return result1.(types.BeaconBlock), result2.(bool), nil
}

// Get the Beacon chain's head information
func (m *BeaconClientManager) GetBeaconHead(ctx context.Context) (types.BeaconHead, error) {
	result, err := m.runFunction1(func(client types.IBeaconClient) (interface{}, error) {
		return client.GetBeaconHead(ctx)
	})
	if err != nil {
		return types.BeaconHead{}, err
	}
	return result.(types.BeaconHead), nil
}

// Get a validator's status by its index
func (m *BeaconClientManager) GetValidatorStatusByIndex(ctx context.Context, index string, opts *types.ValidatorStatusOptions) (types.ValidatorStatus, error) {
	result, err := m.runFunction1(func(client types.IBeaconClient) (interface{}, error) {
		return client.GetValidatorStatusByIndex(ctx, index, opts)
	})
	if err != nil {
		return types.ValidatorStatus{}, err
	}
	return result.(types.ValidatorStatus), nil
}

// Get a validator's status by its pubkey
func (m *BeaconClientManager) GetValidatorStatus(ctx context.Context, pubkey beaconutils.ValidatorPubkey, opts *types.ValidatorStatusOptions) (types.ValidatorStatus, error) {
	result, err := m.runFunction1(func(client types.IBeaconClient) (interface{}, error) {
		return client.GetValidatorStatus(ctx, pubkey, opts)
	})
	if err != nil {
		return types.ValidatorStatus{}, err
	}
	return result.(types.ValidatorStatus), nil
}

// Get the statuses of multiple validators by their pubkeys
func (m *BeaconClientManager) GetValidatorStatuses(ctx context.Context, pubkeys []beaconutils.ValidatorPubkey, opts *types.ValidatorStatusOptions) (map[beaconutils.ValidatorPubkey]types.ValidatorStatus, error) {
	result, err := m.runFunction1(func(client types.IBeaconClient) (interface{}, error) {
		return client.GetValidatorStatuses(ctx, pubkeys, opts)
	})
	if err != nil {
		return nil, err
	}
	return result.(map[beaconutils.ValidatorPubkey]types.ValidatorStatus), nil
}

// Get a validator's index
func (m *BeaconClientManager) GetValidatorIndex(ctx context.Context, pubkey beaconutils.ValidatorPubkey) (string, error) {
	result, err := m.runFunction1(func(client types.IBeaconClient) (interface{}, error) {
		return client.GetValidatorIndex(ctx, pubkey)
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

// Get a validator's sync duties
func (m *BeaconClientManager) GetValidatorSyncDuties(ctx context.Context, indices []string, epoch uint64) (map[string]bool, error) {
	result, err := m.runFunction1(func(client types.IBeaconClient) (interface{}, error) {
		return client.GetValidatorSyncDuties(ctx, indices, epoch)
	})
	if err != nil {
		return nil, err
	}
	return result.(map[string]bool), nil
}

// Get a validator's proposer duties
func (m *BeaconClientManager) GetValidatorProposerDuties(ctx context.Context, indices []string, epoch uint64) (map[string]uint64, error) {
	result, err := m.runFunction1(func(client types.IBeaconClient) (interface{}, error) {
		return client.GetValidatorProposerDuties(ctx, indices, epoch)
	})
	if err != nil {
		return nil, err
	}
	return result.(map[string]uint64), nil
}

// Get the Beacon chain's domain data
func (m *BeaconClientManager) GetDomainData(ctx context.Context, domainType []byte, epoch uint64, useGenesisFork bool) ([]byte, error) {
	result, err := m.runFunction1(func(client types.IBeaconClient) (interface{}, error) {
		return client.GetDomainData(ctx, domainType, epoch, useGenesisFork)
	})
	if err != nil {
		return nil, err
	}
	return result.([]byte), nil
}

// Voluntarily exit a validator
func (m *BeaconClientManager) ExitValidator(ctx context.Context, validatorIndex string, epoch uint64, signature beaconutils.ValidatorSignature) error {
	err := m.runFunction0(func(client types.IBeaconClient) error {
		return client.ExitValidator(ctx, validatorIndex, epoch, signature)
	})
	return err
}

// Close the connection to the Beacon client
func (m *BeaconClientManager) Close(ctx context.Context) error {
	err := m.runFunction0(func(client types.IBeaconClient) error {
		return client.Close(ctx)
	})
	return err
}

// Get the EL data for a CL block
func (m *BeaconClientManager) GetEth1DataForEth2Block(ctx context.Context, blockId string) (types.Eth1Data, bool, error) {
	result1, result2, err := m.runFunction2(func(client types.IBeaconClient) (interface{}, interface{}, error) {
		return client.GetEth1DataForEth2Block(ctx, blockId)
	})
	if err != nil {
		return types.Eth1Data{}, false, err
	}
	return result1.(types.Eth1Data), result2.(bool), nil
}

// Change the withdrawal credentials for a validator
func (m *BeaconClientManager) ChangeWithdrawalCredentials(ctx context.Context, validatorIndex string, fromBlsPubkey beaconutils.ValidatorPubkey, toExecutionAddress common.Address, signature beaconutils.ValidatorSignature) error {
	err := m.runFunction0(func(client types.IBeaconClient) error {
		return client.ChangeWithdrawalCredentials(ctx, validatorIndex, fromBlsPubkey, toExecutionAddress, signature)
	})
	if err != nil {
		return err
	}
	return nil
}

/// ==================
/// Internal Functions
/// ==================

func (m *BeaconClientManager) CheckStatus(ctx context.Context) *api.ClientManagerStatus {

	status := &api.ClientManagerStatus{
		FallbackEnabled: m.fallbackBc != nil,
	}

	// Ignore the sync check and just use the predefined settings if requested
	if m.ignoreSyncCheck {
		status.PrimaryClientStatus.IsWorking = m.primaryReady
		status.PrimaryClientStatus.IsSynced = m.primaryReady
		if status.FallbackEnabled {
			status.FallbackClientStatus.IsWorking = m.fallbackReady
			status.FallbackClientStatus.IsSynced = m.fallbackReady
		}
		return status
	}

	// Get the primary BC status
	status.PrimaryClientStatus = checkBcStatus(ctx, m.primaryBc)

	// Get the fallback BC status if applicable
	if status.FallbackEnabled {
		status.FallbackClientStatus = checkBcStatus(ctx, m.fallbackBc)
	}

	// Flag the ready clients
	m.primaryReady = (status.PrimaryClientStatus.IsWorking && status.PrimaryClientStatus.IsSynced)
	m.fallbackReady = (status.FallbackEnabled && status.FallbackClientStatus.IsWorking && status.FallbackClientStatus.IsSynced)

	return status

}

// Check the client status
func checkBcStatus(ctx context.Context, client types.IBeaconClient) api.ClientStatus {

	status := api.ClientStatus{}

	// Get the fallback's sync progress
	syncStatus, err := client.GetSyncStatus(ctx)
	if err != nil {
		status.Error = fmt.Sprintf("Sync progress check failed with [%s]", err.Error())
		status.IsSynced = false
		status.IsWorking = false
		return status
	}

	// Return the sync status
	if !syncStatus.Syncing {
		status.IsWorking = true
		status.IsSynced = true
		status.SyncProgress = 1
	} else {
		status.IsWorking = true
		status.IsSynced = false
		status.SyncProgress = syncStatus.Progress
	}
	return status

}

// Attempts to run a function progressively through each client until one succeeds or they all fail.
func (m *BeaconClientManager) runFunction0(function bcFunction0) error {

	// Check if we can use the primary
	if m.primaryReady {
		// Try to run the function on the primary
		err := function(m.primaryBc)
		if err != nil {
			if m.isDisconnected(err) {
				// If it's disconnected, log it and try the fallback
				m.logger.Printlnf("WARNING: Primary Beacon client disconnected (%s), using fallback...", err.Error())
				m.primaryReady = false
				return m.runFunction0(function)
			}
			// If it's a different error, just return it
			return err
		}
		// If there's no error, return the result
		return nil
	}

	if m.fallbackReady {
		// Try to run the function on the fallback
		err := function(m.fallbackBc)
		if err != nil {
			if m.isDisconnected(err) {
				// If it's disconnected, log it and try the fallback
				m.logger.Printlnf("WARNING: Fallback Beacon client disconnected (%s)", err.Error())
				m.fallbackReady = false
				return fmt.Errorf("all Beacon clients failed")
			}

			// If it's a different error, just return it
			return err
		}
		// If there's no error, return the result
		return nil
	}

	return fmt.Errorf("no Beacon clients were ready")
}

// Attempts to run a function progressively through each client until one succeeds or they all fail.
func (m *BeaconClientManager) runFunction1(function bcFunction1) (interface{}, error) {

	// Check if we can use the primary
	if m.primaryReady {
		// Try to run the function on the primary
		result, err := function(m.primaryBc)
		if err != nil {
			if m.isDisconnected(err) {
				// If it's disconnected, log it and try the fallback
				m.logger.Printlnf("WARNING: Primary Beacon client disconnected (%s), using fallback...", err.Error())
				m.primaryReady = false
				return m.runFunction1(function)
			}
			// If it's a different error, just return it
			return nil, err
		}
		// If there's no error, return the result
		return result, nil
	}

	if m.fallbackReady {
		// Try to run the function on the fallback
		result, err := function(m.fallbackBc)
		if err != nil {
			if m.isDisconnected(err) {
				// If it's disconnected, log it and try the fallback
				m.logger.Printlnf("WARNING: Fallback Beacon client disconnected (%s)", err.Error())
				m.fallbackReady = false
				return nil, fmt.Errorf("all Beacon clients failed")
			}
			// If it's a different error, just return it
			return nil, err
		}
		// If there's no error, return the result
		return result, nil
	}

	return nil, fmt.Errorf("no Beacon clients were ready")

}

// Attempts to run a function progressively through each client until one succeeds or they all fail.
func (m *BeaconClientManager) runFunction2(function bcFunction2) (interface{}, interface{}, error) {

	// Check if we can use the primary
	if m.primaryReady {
		// Try to run the function on the primary
		result1, result2, err := function(m.primaryBc)
		if err != nil {
			if m.isDisconnected(err) {
				// If it's disconnected, log it and try the fallback
				m.logger.Printlnf("WARNING: Primary Beacon client request failed (%s), using fallback...", err.Error())
				m.primaryReady = false
				return m.runFunction2(function)
			}
			// If it's a different error, just return it
			return nil, nil, err
		}
		// If there's no error, return the result
		return result1, result2, nil
	}

	if m.fallbackReady {
		// Try to run the function on the fallback
		result1, result2, err := function(m.fallbackBc)
		if err != nil {
			if m.isDisconnected(err) {
				// If it's disconnected, log it and try the fallback
				m.logger.Printlnf("WARNING: Fallback Beacon client request failed (%s)", err.Error())
				m.fallbackReady = false
				return nil, nil, fmt.Errorf("all Beacon clients failed")
			}
			// If it's a different error, just return it
			return nil, nil, err
		}
		// If there's no error, return the result
		return result1, result2, nil
	}

	return nil, nil, fmt.Errorf("no Beacon clients were ready")

}

// Returns true if the error was a connection failure and a backup client is available
func (m *BeaconClientManager) isDisconnected(err error) bool {
	var sysErr syscall.Errno
	if errors.As(err, &sysErr) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr)
}
