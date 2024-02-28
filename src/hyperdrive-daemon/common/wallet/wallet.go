package wallet

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	sharedtypes "github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/nodeset-org/hyperdrive/shared/utils/log"
	"github.com/tyler-smith/go-bip39"
)

// Config
const (
	EntropyBits                          = 256
	FileMode                             = 0600
	DefaultNodeKeyPath                   = "m/44'/60'/0'/0/%d"
	LedgerLiveNodeKeyPath                = "m/44'/60'/%d/0/0"
	MyEtherWalletNodeKeyPath             = "m/44'/60'/0'/%d"
	walletFileMode           fs.FileMode = 0600
)

// Wallet
type Wallet struct {
	// Managers
	walletManager   IWalletManager
	addressManager  *AddressManager
	passwordManager *PasswordManager

	// Misc cache
	chainID        uint
	walletDataPath string
}

// Create new wallet
func NewWallet(log *log.ColorLogger, walletDataPath string, walletAddressPath string, passwordFilePath string, chainID uint) (*Wallet, error) {
	// Create the wallet
	w := &Wallet{
		// Create managers
		addressManager:  NewAddressManager(walletAddressPath),
		passwordManager: NewPasswordManager(passwordFilePath),

		// Initialize other fields
		chainID:        chainID,
		walletDataPath: walletDataPath,
	}

	// Load the password
	password, isPasswordSaved, err := w.passwordManager.GetPasswordFromDisk()
	if err != nil {
		return nil, fmt.Errorf("error loading password: %w", err)
	}

	// Load the wallet
	if isPasswordSaved {
		walletMgr, err := w.loadWalletData(password)
		if err != nil {
			log.Printlnf("[WALLET] Loading wallet with stored node password failed: %s", err.Error())
		} else if walletMgr != nil {
			w.walletManager = walletMgr
		}
	}

	// Load the node address
	_, _, err = w.addressManager.LoadAddress()
	if err != nil {
		return nil, fmt.Errorf("error loading node address: %w", err)
	}
	return w, nil
}

// Gets the status of the wallet and its artifacts
func (w *Wallet) GetStatus() (sharedtypes.WalletStatus, error) {
	// Make a status wrapper
	status := sharedtypes.WalletStatus{}

	// Get the password details
	var err error
	_, status.Password.IsPasswordSaved, err = w.passwordManager.GetPasswordFromDisk()
	if err != nil {
		return status, fmt.Errorf("error checking password manager status: %w", err)
	}

	// Get the wallet details
	if w.walletManager != nil {
		status.Wallet.IsLoaded = true
		status.Wallet.Type = w.walletManager.GetType()
		status.Wallet.IsOnDisk = true
		status.Wallet.WalletAddress, err = w.walletManager.GetAddress()
		if err != nil {
			return status, fmt.Errorf("error getting wallet address: %w", err)
		}
	} else {
		status.Wallet.IsOnDisk, err = w.isWalletDataOnDisk()
		if err != nil {
			return status, fmt.Errorf("error checking if wallet data is on disk: %w", err)
		}
	}

	// Get the address details
	status.Address.NodeAddress, status.Address.HasAddress = w.addressManager.GetAddress()
	return status, nil
}

// Get the node address, if one is loaded
func (w *Wallet) GetAddress() (common.Address, bool) {
	return w.addressManager.GetAddress()
}

// Get the transactor that can sign transactions
func (w *Wallet) GetTransactor() (*bind.TransactOpts, error) {
	if w.walletManager == nil {
		return nil, fmt.Errorf("wallet is not loaded")
	}
	return w.walletManager.GetTransactor()
}

// Sign a message with the wallet's private key
func (w *Wallet) SignMessage(message []byte) ([]byte, error) {
	if w.walletManager == nil {
		return nil, fmt.Errorf("wallet is not loaded")
	}
	return w.walletManager.SignMessage(message)
}

// Sign a transaction with the wallet's private key
func (w *Wallet) SignTransaction(serializedTx []byte) ([]byte, error) {
	if w.walletManager == nil {
		return nil, fmt.Errorf("wallet is not loaded")
	}
	return w.walletManager.SignTransaction(serializedTx)
}

// Masquerade as another node address, running all node functions as that address (in read only moe)
func (w *Wallet) MasqueradeAsAddress(newAddress common.Address) error {
	return w.addressManager.SetAndSaveAddress(newAddress)
}

// End masquerading as another node address, and use the wallet's address (returning to read/write mode)
func (w *Wallet) RestoreAddressToWallet() error {
	if w.addressManager == nil {
		return fmt.Errorf("wallet is not loaded")
	}

	walletAddress, err := w.walletManager.GetAddress()
	if err != nil {
		return fmt.Errorf("error getting wallet address: %w", err)
	}

	return w.MasqueradeAsAddress(walletAddress)
}

// Initialize the wallet from a random seed
func (w *Wallet) CreateNewLocalWallet(derivationPath string, walletIndex uint, password string, savePassword bool) (string, error) {
	if w.walletManager != nil {
		return "", fmt.Errorf("wallet keystore is already present - please delete it before creating a new wallet")
	}

	// Make a mnemonic
	mnemonic, err := GenerateNewMnemonic()
	if err != nil {
		return "", err
	}

	// Initialize the wallet with it
	err = w.buildLocalWallet(derivationPath, walletIndex, mnemonic, password, savePassword, false)
	if err != nil {
		return "", fmt.Errorf("error initializing new wallet keystore: %w", err)
	}
	return mnemonic, nil
}

// Recover a local wallet from a mnemonic
func (w *Wallet) Recover(derivationPath string, walletIndex uint, mnemonic string, password string, savePassword bool, testMode bool) error {
	if w.walletManager != nil {
		return fmt.Errorf("wallet keystore is already present - please delete it before recovering an existing wallet")
	}

	// Check the mnemonic
	if !bip39.IsMnemonicValid(mnemonic) {
		return fmt.Errorf("invalid mnemonic '%s'", mnemonic)
	}

	return w.buildLocalWallet(derivationPath, walletIndex, mnemonic, password, savePassword, testMode)
}

// Attempts to load the wallet keystore with the provided password if not set
func (w *Wallet) SetPassword(password string, save bool) error {
	if w.walletManager != nil {
		if !save {
			return fmt.Errorf("wallet is already loaded, nothing to do")
		}

		switch w.walletManager.GetType() {
		case sharedtypes.WalletType_Local:
			// Make sure the password is correct
			localMgr := w.walletManager.(*LocalWalletManager)
			isValid, err := localMgr.VerifyPassword(password)
			if err != nil {
				return fmt.Errorf("error setting password: %w", err)
			}
			if !isValid {
				return fmt.Errorf("provided password is not correct for the loaded wallet")
			}

			// Save and exit
			return w.passwordManager.SavePassword(password)
		default:
			return fmt.Errorf("loaded wallet is not local and does not use a password")
		}
	}

	// Try to load the wallet with the new password
	isWalletOnDisk, err := w.isWalletDataOnDisk()
	if err != nil {
		return fmt.Errorf("error checking if wallet data is on disk: %w", err)
	}
	if !isWalletOnDisk {
		return fmt.Errorf("keystore not present, wallet must be initialized or recovered first")
	}
	mgr, err := w.loadWalletData(password)
	if err != nil {
		return fmt.Errorf("error loading wallet with provided password: %w", err)
	}

	// Save if requested
	if save {
		err := w.passwordManager.SavePassword(password)
		if err != nil {
			return err
		}
	}

	// Set the wallet manager
	w.walletManager = mgr
	return nil
}

// Retrieves the wallet's password
func (w *Wallet) GetPassword() (string, bool, error) {
	return w.passwordManager.GetPasswordFromDisk()
}

// Delete the wallet password from disk, but retain it in memory if a local keystore is already loaded
func (w *Wallet) DeletePassword() error {
	err := w.passwordManager.DeletePassword()
	if err != nil {
		return fmt.Errorf("error deleting wallet password: %w", err)
	}
	return nil
}

// Get the node account private key bytes
func (w *Wallet) GetNodePrivateKeyBytes() ([]byte, error) {
	if w.walletManager == nil {
		return nil, fmt.Errorf("wallet is not loaded")
	}

	switch w.walletManager.GetType() {
	case sharedtypes.WalletType_Local:
		localMgr := w.walletManager.(*LocalWalletManager)
		return crypto.FromECDSA(localMgr.GetPrivateKey()), nil
	default:
		return nil, fmt.Errorf("loaded wallet is not local")
	}
}

// Get the node account private key bytes
func (w *Wallet) GetEthKeystore(password string) ([]byte, error) {
	if w.walletManager == nil {
		return nil, fmt.Errorf("wallet is not loaded")
	}

	switch w.walletManager.GetType() {
	case sharedtypes.WalletType_Local:
		localMgr := w.walletManager.(*LocalWalletManager)
		return localMgr.GetEthKeystore(password)
	default:
		return nil, fmt.Errorf("loaded wallet is not local")
	}
}

// Serialize the wallet data as JSON
func (w *Wallet) SerializeData() (string, error) {
	if w.walletManager == nil {
		return "", fmt.Errorf("wallet is not loaded")
	}
	return w.walletManager.SerializeData()
}

// Generate a BLS validator key from the provided path, using the node wallet's seed as a basis
func (w *Wallet) GenerateValidatorKey(path string) ([]byte, error) {
	if w.walletManager == nil {
		return nil, fmt.Errorf("wallet is not loaded")
	}

	switch w.walletManager.GetType() {
	case sharedtypes.WalletType_Local:
		localMgr := w.walletManager.(*LocalWalletManager)
		return localMgr.GenerateValidatorKey(path)
	default:
		return nil, fmt.Errorf("loaded wallet is not local")
	}
}

// Builds a local wallet keystore and saves its artifacts to disk
func (w *Wallet) buildLocalWallet(derivationPath string, walletIndex uint, mnemonic string, password string, savePassword bool, testMode bool) error {
	// Initialize the wallet with it
	localMgr := NewLocalWalletManager(w.chainID)
	localData, err := localMgr.InitializeKeystore(derivationPath, walletIndex, mnemonic, password)
	if err != nil {
		return fmt.Errorf("error initializing wallet keystore with recovered data: %w", err)
	}

	// Get the wallet address
	walletAddress, _ := localMgr.GetAddress()

	if !testMode {
		// Create data
		data := &sharedtypes.WalletData{
			Type:      sharedtypes.WalletType_Local,
			LocalData: *localData,
		}

		// Save the wallet data
		err = w.saveWalletData(data)
		if err != nil {
			return fmt.Errorf("error saving wallet data: %w", err)
		}
		// Update the address file
		err = w.addressManager.SetAndSaveAddress(walletAddress)
		if err != nil {
			return fmt.Errorf("error saving wallet address to node address file: %w", err)
		}

		if savePassword {
			err := w.passwordManager.SavePassword(password)
			if err != nil {
				return fmt.Errorf("error saving password: %w", err)
			}
		}
	} else {
		w.addressManager.SetAddress(walletAddress)
	}

	w.walletManager = localMgr
	return nil
}

// Check if the wallet file is saved to disk
func (w *Wallet) isWalletDataOnDisk() (bool, error) {
	// Read the file
	_, err := os.Stat(w.walletDataPath)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("error checking if wallet file [%s] exists: %w", w.walletDataPath, err)
	}
	return true, nil
}

// Load the wallet data from disk
func (w *Wallet) loadWalletData(password string) (IWalletManager, error) {
	// Read the file
	bytes, err := os.ReadFile(w.walletDataPath)
	if err != nil {
		return nil, fmt.Errorf("error reading wallet data at [%s]: %w", w.walletDataPath, err)
	}

	// Deserialize it
	data := new(sharedtypes.WalletData)
	err = json.Unmarshal(bytes, data)
	if err != nil {
		return nil, fmt.Errorf("error deserializing wallet data at [%s]: %w", w.walletDataPath, err)
	}

	// Load the proper type
	var manager IWalletManager
	switch data.Type {
	case sharedtypes.WalletType_Local:
		localMgr := NewLocalWalletManager(w.chainID)
		err = localMgr.LoadWallet(&data.LocalData, password)
		if err != nil {
			return nil, fmt.Errorf("error loading local wallet data at %s: %w", w.walletDataPath, err)
		}
		manager = localMgr
	default:
		return nil, fmt.Errorf("unsupported wallet type: %s", data.Type)
	}

	// Data loaded!
	return manager, nil
}

// Save the wallet data to disk
func (w *Wallet) saveWalletData(data *sharedtypes.WalletData) error {
	// Serialize it
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error serializing wallet data: %w", err)
	}

	// Write the file
	err = os.WriteFile(w.walletDataPath, bytes, walletFileMode)
	if err != nil {
		return fmt.Errorf("error writing wallet data to [%s]: %w", w.walletDataPath, err)
	}
	return nil
}

// =============
// === Utils ===
// =============

// Generate a new random mnemonic and seed
func GenerateNewMnemonic() (string, error) {
	// Generate random entropy for the mnemonic
	entropy, err := bip39.NewEntropy(EntropyBits)
	if err != nil {
		return "", fmt.Errorf("error generating wallet mnemonic entropy bytes: %w", err)
	}

	// Generate a new mnemonic
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("error generating wallet mnemonic: %w", err)
	}
	return mnemonic, nil
}
