package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/node-manager-core/log"
	"github.com/rocket-pool/node-manager-core/node/services"
	"github.com/rocket-pool/node-manager-core/utils"
	"github.com/rocket-pool/node-manager-core/wallet"
)

const (
	// Log keys
	PrimarySyncProgressKey  string = "primarySyncProgress"
	FallbackSyncProgressKey string = "fallbackSyncProgress"
	SyncProgressKey         string = "syncProgress"
	PrimaryErrorKey         string = "primaryError"
	FallbackErrorKey        string = "fallbackError"

	ethClientStatusRefreshInterval time.Duration = 60 * time.Second
	ethClientSyncPollInterval      time.Duration = 5 * time.Second
	beaconClientSyncPollInterval   time.Duration = 5 * time.Second
	walletReadyCheckInterval       time.Duration = 15 * time.Second
)

func (sp *ServiceProvider) RequireNodeAddress(status wallet.WalletStatus) error {
	if !status.Address.HasAddress {
		return errors.New("The node currently does not have an address set. Please run 'hyperdrive wallet init' and try again.")
	}
	return nil
}

func (sp *ServiceProvider) RequireWalletReady(status wallet.WalletStatus) error {
	return CheckIfWalletReady(status)
}

func (sp *ServiceProvider) RequireEthClientSynced(ctx context.Context) error {
	synced, _, err := sp.checkExecutionClientStatus(ctx)
	if err != nil {
		return err
	}
	if synced {
		return nil
	}
	return errors.New("The Execution client is currently syncing. Please try again later.")
}

func (sp *ServiceProvider) RequireBeaconClientSynced(ctx context.Context) error {
	synced, err := sp.checkBeaconClientStatus(ctx)
	if err != nil {
		return err
	}
	if synced {
		return nil
	}
	return errors.New("The Beacon client is currently syncing. Please try again later.")
}

// Wait for the Executon client to sync; timeout of 0 indicates no timeout
func (sp *ServiceProvider) WaitEthClientSynced(ctx context.Context, verbose bool) error {
	_, err := sp.waitEthClientSynced(ctx, verbose)
	return err
}

// Wait for the Beacon client to sync; timeout of 0 indicates no timeout
func (sp *ServiceProvider) WaitBeaconClientSynced(ctx context.Context, verbose bool) error {
	_, err := sp.waitBeaconClientSynced(ctx, verbose)
	return err
}

// Wait for the Hyperdrive wallet to be ready
func (sp *ServiceProvider) WaitForWallet(ctx context.Context) error {
	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}

	for {
		hdWalletStatus, err := sp.GetHyperdriveClient().Wallet.Status()
		if err != nil {
			return fmt.Errorf("error getting Hyperdrive wallet status: %w", err)
		}

		if CheckIfWalletReady(hdWalletStatus.Data.WalletStatus) == nil {
			return nil
		}

		logger.Info("Hyperdrive wallet not ready yet",
			slog.Duration("retry", walletReadyCheckInterval),
		)
		if utils.SleepWithCancel(ctx, walletReadyCheckInterval) {
			return nil
		}
	}
}

// Check if the primary and fallback Execution clients are synced
// TODO: Move this into ec-manager and stop exposing the primary and fallback directly...
func (sp *ServiceProvider) checkExecutionClientStatus(ctx context.Context) (bool, eth.IExecutionClient, error) {
	// Check the EC status
	ecMgr := sp.GetEthClient()
	mgrStatus := ecMgr.CheckStatus(ctx, true) // Always check the chain ID for now
	if ecMgr.IsPrimaryReady() {
		return true, nil, nil
	}

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}

	// If the primary isn't synced but there's a fallback and it is, return true
	if ecMgr.IsFallbackReady() {
		if mgrStatus.PrimaryClientStatus.Error != "" {
			logger.Warn("Primary execution client is unavailable using fallback execution client...", slog.String(log.ErrorKey, mgrStatus.PrimaryClientStatus.Error))
		} else {
			logger.Warn("Primary execution client is still syncing, using fallback execution client...", slog.Float64(PrimarySyncProgressKey, mgrStatus.PrimaryClientStatus.SyncProgress*100))
		}
		return true, nil, nil
	}

	// If neither is synced, go through the status to figure out what to do

	// Is the primary working and syncing? If so, wait for it
	if mgrStatus.PrimaryClientStatus.IsWorking && mgrStatus.PrimaryClientStatus.Error == "" {
		logger.Error("Fallback execution client is not configured or unavailable, waiting for primary execution client to finish syncing", slog.Float64(PrimarySyncProgressKey, mgrStatus.PrimaryClientStatus.SyncProgress*100))
		return false, ecMgr.GetPrimaryClient(), nil
	}

	// Is the fallback working and syncing? If so, wait for it
	if mgrStatus.FallbackEnabled && mgrStatus.FallbackClientStatus.IsWorking && mgrStatus.FallbackClientStatus.Error == "" {
		logger.Error("Primary execution client is unavailable, waiting for the fallback execution client to finish syncing", slog.String(PrimaryErrorKey, mgrStatus.PrimaryClientStatus.Error), slog.Float64(FallbackSyncProgressKey, mgrStatus.FallbackClientStatus.SyncProgress*100))
		return false, ecMgr.GetFallbackClient(), nil
	}

	// If neither client is working, report the errors
	if mgrStatus.FallbackEnabled {
		return false, nil, fmt.Errorf("Primary execution client is unavailable (%s) and fallback execution client is unavailable (%s), no execution clients are ready.", mgrStatus.PrimaryClientStatus.Error, mgrStatus.FallbackClientStatus.Error)
	}

	return false, nil, fmt.Errorf("Primary execution client is unavailable (%s) and no fallback execution client is configured.", mgrStatus.PrimaryClientStatus.Error)
}

// Check if the primary and fallback Beacon clients are synced
func (sp *ServiceProvider) checkBeaconClientStatus(ctx context.Context) (bool, error) {
	// Check the BC status
	bcMgr := sp.GetBeaconClient()
	mgrStatus := bcMgr.CheckStatus(ctx, true) // Always check the chain ID for now
	if bcMgr.IsPrimaryReady() {
		return true, nil
	}

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}

	// If the primary isn't synced but there's a fallback and it is, return true
	if bcMgr.IsFallbackReady() {
		if mgrStatus.PrimaryClientStatus.Error != "" {
			logger.Warn("Primary Beacon Node is unavailable, using fallback Beacon Node...", slog.String(PrimaryErrorKey, mgrStatus.PrimaryClientStatus.Error))
		} else {
			logger.Warn("Primary Beacon Node is still syncing, using fallback Beacon Node...", slog.Float64(PrimarySyncProgressKey, mgrStatus.PrimaryClientStatus.SyncProgress*100))
		}
		return true, nil
	}

	// If neither is synced, go through the status to figure out what to do

	// Is the primary working and syncing? If so, wait for it
	if mgrStatus.PrimaryClientStatus.IsWorking && mgrStatus.PrimaryClientStatus.Error == "" {
		logger.Error("Fallback Beacon Node is not configured or unavailable, waiting for primary Beacon Node to finish syncing...", slog.Float64(PrimarySyncProgressKey, mgrStatus.PrimaryClientStatus.SyncProgress*100))
		return false, nil
	}

	// Is the fallback working and syncing? If so, wait for it
	if mgrStatus.FallbackEnabled && mgrStatus.FallbackClientStatus.IsWorking && mgrStatus.FallbackClientStatus.Error == "" {
		logger.Error("Primary Beacon Node is unavailable, waiting for the fallback Beacon Node to finish syncing...", slog.String(PrimaryErrorKey, mgrStatus.PrimaryClientStatus.Error), slog.Float64(FallbackSyncProgressKey, mgrStatus.FallbackClientStatus.SyncProgress*100))
		return false, nil
	}

	// If neither client is working, report the errors
	if mgrStatus.FallbackEnabled {
		return false, fmt.Errorf("Primary Beacon Node is unavailable (%s) and fallback Beacon Node is unavailable (%s), no Beacon Nodes are ready.", mgrStatus.PrimaryClientStatus.Error, mgrStatus.FallbackClientStatus.Error)
	}

	return false, fmt.Errorf("Primary Beacon Node is unavailable (%s) and no fallback Beacon Node is configured.", mgrStatus.PrimaryClientStatus.Error)
}

// Wait for the primary or fallback Execution client to be synced
func (sp *ServiceProvider) waitEthClientSynced(ctx context.Context, verbose bool) (bool, error) {
	synced, clientToCheck, err := sp.checkExecutionClientStatus(ctx)
	if err != nil {
		return false, err
	}
	if synced {
		return true, nil
	}

	// Get EC status refresh time
	ecRefreshTime := time.Now()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}

	// Wait for sync
	for {
		// Check if the EC status needs to be refreshed
		if time.Since(ecRefreshTime) > ethClientStatusRefreshInterval {
			logger.Info("Refreshing primary / fallback execution client status...")
			ecRefreshTime = time.Now()
			synced, clientToCheck, err = sp.checkExecutionClientStatus(ctx)
			if err != nil {
				return false, err
			}
			if synced {
				return true, nil
			}
		}

		// Get sync progress
		progress, err := clientToCheck.SyncProgress(ctx)
		if err != nil {
			return false, err
		}

		// Check sync progress
		if progress != nil {
			if verbose {
				p := float64(progress.CurrentBlock-progress.StartingBlock) / float64(progress.HighestBlock-progress.StartingBlock)
				if p > 1 {
					logger.Info("Execution client syncing...")
				} else {
					logger.Info("Execution client syncing...", slog.Float64(SyncProgressKey, p*100))
				}
			}
		} else {
			// Eth 1 client is not in "syncing" state but may be behind head
			// Get the latest block it knows about and make sure it's recent compared to system clock time
			isUpToDate, _, err := services.IsSyncWithinThreshold(clientToCheck)
			if err != nil {
				return false, err
			}
			// Only return true if the last reportedly known block is within our defined threshold
			if isUpToDate {
				return true, nil
			}
		}

		// Pause before next poll
		time.Sleep(ethClientSyncPollInterval)
	}
}

// Wait for the primary or fallback Beacon client to be synced
func (sp *ServiceProvider) waitBeaconClientSynced(ctx context.Context, verbose bool) (bool, error) {
	synced, err := sp.checkBeaconClientStatus(ctx)
	if err != nil {
		return false, err
	}
	if synced {
		return true, nil
	}

	// Get BC status refresh time
	bcRefreshTime := time.Now()

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}

	// Wait for sync
	for {
		// Check if the BC status needs to be refreshed
		if time.Since(bcRefreshTime) > ethClientStatusRefreshInterval {
			logger.Info("Refreshing primary / fallback Beacon Node status...")
			bcRefreshTime = time.Now()
			synced, err = sp.checkBeaconClientStatus(ctx)
			if err != nil {
				return false, err
			}
			if synced {
				return true, nil
			}
		}

		// Get sync status
		syncStatus, err := sp.GetBeaconClient().GetSyncStatus(ctx)
		if err != nil {
			return false, err
		}

		// Check sync status
		if syncStatus.Syncing {
			if verbose {
				logger.Info("Beacon Node syncing...", slog.Float64(SyncProgressKey, syncStatus.Progress*100))
			}
		} else {
			return true, nil
		}

		// Pause before next poll
		time.Sleep(beaconClientSyncPollInterval)
	}
}
