package swcommon

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive/daemon-utils/services"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/rocket-pool/node-manager-core/node/validator"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
)

const (
	walletDataFilename string = "wallet_data"
)

// Data relating to Stakewise's wallet
type stakewiseWalletData struct {
	// The next account to generate the key for
	NextAccount uint64 `json:"nextAccount"`

	// The ID of the nodeset deposit data stored on disk
	NodeSetDepositDataVersion int `json:"nodeSetDepositDataVersion"`
}

// Wallet manager for the Stakewise daemon
type Wallet struct {
	validatorManager         *validator.ValidatorManager
	stakewiseKeystoreManager *stakewiseKeystoreManager
	data                     stakewiseWalletData
	sp                       *services.ServiceProvider
}

// Create a new wallet
func NewWallet(sp *services.ServiceProvider) (*Wallet, error) {
	moduleDir := sp.GetModuleDir()
	validatorPath := filepath.Join(moduleDir, config.ValidatorsDirectory)
	wallet := &Wallet{
		sp:               sp,
		validatorManager: validator.NewValidatorManager(validatorPath),
	}

	// Check if the wallet data exists
	dataPath := filepath.Join(moduleDir, walletDataFilename)
	_, err := os.Stat(dataPath)
	if errors.Is(err, fs.ErrNotExist) {
		// No data yet, so make some
		wallet.data = stakewiseWalletData{
			NextAccount:               0,
			NodeSetDepositDataVersion: 0,
		}

		// Save it
		err = wallet.saveData()
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, fmt.Errorf("error checking status of wallet file [%s]: %w", dataPath, err)
	} else {
		// Read it
		bytes, err := os.ReadFile(dataPath)
		if err != nil {
			return nil, fmt.Errorf("error loading wallet data: %w", err)
		}
		var data stakewiseWalletData
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			return nil, fmt.Errorf("error deserializing wallet data: %w", err)
		}
		wallet.data = data
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
	path := fmt.Sprintf(shared.StakewiseValidatorPath, w.data.NextAccount)

	// Ask the HD daemon to generate the key
	client := w.sp.GetHyperdriveClient()
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

// Get the private validator key with the corresponding pubkey
func (w *Wallet) GetPrivateKeyForPubkey(pubkey beacon.ValidatorPubkey) (*eth2types.BLSPrivateKey, error) {
	return w.stakewiseKeystoreManager.LoadValidatorKey(pubkey)
}

// Get the private validator key with the corresponding pubkey
func (w *Wallet) DerivePubKeys(privateKeys []*eth2types.BLSPrivateKey) ([]beacon.ValidatorPubkey, error) {
	publicKeys := make([]beacon.ValidatorPubkey, 0, len(privateKeys))

	for i, privateKey := range privateKeys {
		if privateKey == nil {
			return nil, fmt.Errorf("nil private key encountered at index %d", i)
		}

		validatorPubkey := beacon.ValidatorPubkey(privateKey.PublicKey().Marshal())
		publicKeys = append(publicKeys, validatorPubkey)
	}

	return publicKeys, nil
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
		if !strings.HasPrefix(filename, keystorePrefix) || !strings.HasSuffix(filename, keystoreSuffix) {
			continue
		}

		// Get the pubkey from the filename
		trimmed := strings.TrimPrefix(filename, keystorePrefix)
		trimmed = strings.TrimSuffix(trimmed, keystoreSuffix)
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

// Get the version of the aggregated deposit data from the NodeSet server that's stored on disk
func (w *Wallet) GetLatestDepositDataVersion() int {
	return w.data.NodeSetDepositDataVersion
}

// Set the latest deposit data version and save the wallet data
func (w *Wallet) SetLatestDepositDataVersion(version int) error {
	w.data.NodeSetDepositDataVersion = version
	return w.saveData()
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
