package wallet

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-json"
	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/hyperdrive/modules/common/services"
	"github.com/nodeset-org/hyperdrive/shared/types"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
)

const (
	walletDataFilename string = "wallet_data"
)

// Data relating to Stakewise's wallet
type stakewiseWalletData struct {
	// The next account to generate the key for
	NextAccount uint64 `json:"nextAccount"`
}

// Wallet manager for the Stakewise daemon
type Wallet struct {
	validatorManager         *services.ValidatorManager
	stakewiseKeystoreManager *stakewiseKeystoreManager
	data                     *stakewiseWalletData
	sp                       *services.ServiceProvider
}

// Create a new wallet
func NewWallet(sp *services.ServiceProvider) (*Wallet, error) {
	moduleDir := sp.GetModuleDir()
	wallet := &Wallet{
		sp:               sp,
		validatorManager: services.NewValidatorManager(moduleDir),
	}

	// Check if the wallet data exists
	dataPath := filepath.Join(moduleDir, walletDataFilename)
	_, err := os.Stat(dataPath)
	if errors.Is(err, fs.ErrNotExist) {
		// No data yet, so make some
		wallet.data = &stakewiseWalletData{}
		return wallet, nil
	}

	// Read it
	bytes, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, fmt.Errorf("error loading wallet data: %w", err)
	}
	err = json.Unmarshal(bytes, wallet.data)
	if err != nil {
		return nil, fmt.Errorf("error deserializing wallet data: %w", err)
	}

	// Make the Stakewise keystore manager
	stakewiseKeystoreMgr, err := newStakewiseKeystoreManager(moduleDir)
	if err != nil {
		return nil, fmt.Errorf("error creating Stakewise keystore manager: %w", err)
	}
	wallet.stakewiseKeystoreManager = stakewiseKeystoreMgr

	return wallet, nil
}

// Generate a new validator key and save it
func (w *Wallet) GenerateNewValidatorKey() (*eth2types.BLSPrivateKey, error) {
	// Get the path for the next validator key
	path := fmt.Sprintf(types.StakewiseValidatorPath, w.data.NextAccount)

	// Ask the HD daemon to generate the key
	client := w.sp.GetClient()
	response, err := client.Wallet.GenerateValidatorKey(path)
	if err != nil {
		return nil, fmt.Errorf("error generating validator key for path [%s]: %w", path, err)
	}

	// Increment the next account index first for safety
	w.data.NextAccount++
	err = w.saveData()
	if err != nil {
		return nil, err
	}

	// Save the key to the VC stores
	key, err := eth2types.BLSPrivateKeyFromBytes(response.Data.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error converting BLS private key for path %s: %w", path, err)
	}
	err = w.validatorManager.StoreKey(key, path)
	if err != nil {
		return nil, fmt.Errorf("error saving validator key: %w", err)
	}

	// Save the key to the Stakewise folder
	err = w.stakewiseKeystoreManager.StoreValidatorKey(key, path)
	if err != nil {
		return nil, fmt.Errorf("error saving validator key to the Stakewise store: %w", err)
	}
	return key, nil
}

// Gets all of the validator private keys that are stored in the Stakewise keystore folder
func (w *Wallet) GetAllPrivateKeys() ([]*eth2types.BLSPrivateKey, error) {
	dir := w.stakewiseKeystoreManager.GetKeystoreDir()
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error enumerating Stakewise keystore folder [%s]: %w", dir, err)
	}

	// Go through each file
	keys := []*eth2types.BLSPrivateKey{}
	for _, file := range files {
		filename := file.Name()
		if !strings.HasSuffix(filename, KeystoreSuffix) {
			continue
		}

		// Get the pubkey from the filename
		trimmed := strings.TrimSuffix(filename, KeystoreSuffix)
		pubkey, err := beacon.HexToValidatorPubkey(trimmed)
		if err != nil {
			return nil, fmt.Errorf("error getting pubkey for keystore file [%s]: %w", filename, err)
		}

		// Load it
		key, err := w.stakewiseKeystoreManager.LoadValidatorKey(pubkey)
		if err != nil {
			return nil, fmt.Errorf("error loading validator keystore file [%s]: %w", filename, err)
		}
		keys = append(keys, key)
	}

	return keys, nil
}

// Write the wallet data to disk
func (w *Wallet) saveData() error {
	// Serialize it
	dataPath := filepath.Join(w.sp.GetModuleDir(), walletDataFilename)
	bytes, err := json.Marshal(w.data)
	if err != nil {
		return fmt.Errorf("error serializing wallet data: %w", err)
	}

	// Save it
	err = os.WriteFile(dataPath, bytes, 0600)
	if err != nil {
		return fmt.Errorf("error saving wallet data: %w", err)
	}
	return nil
}
