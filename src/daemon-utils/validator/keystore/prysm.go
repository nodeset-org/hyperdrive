package keystore

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/nodeset-org/hyperdrive/shared/utils"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
	eth2ks "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
)

// Prysm keystore manager
type PrysmKeystoreManager struct {
	as                       *prysmAccountStore
	encryptor                *eth2ks.Encryptor
	keystoreDir              string
	walletDir                string
	accountsDir              string
	keystoreFileName         string
	configFileName           string
	keystorePasswordFileName string
}

type prysmAccountStore struct {
	PrivateKeys [][]byte `json:"private_keys"`
	PublicKeys  [][]byte `json:"public_keys"`
}

// Prysm direct wallet config
type prysmWalletConfig struct {
	DirectEIPVersion string `json:"direct_eip_version"`
}

// Create new prysm keystore manager
func NewPrysmKeystoreManager(keystorePath string) *PrysmKeystoreManager {
	return &PrysmKeystoreManager{
		encryptor:                eth2ks.New(),
		keystoreDir:              filepath.Join(keystorePath, "prysm-non-hd"),
		walletDir:                "direct",
		accountsDir:              "accounts",
		keystoreFileName:         "all-accounts.keystore.json",
		configFileName:           "keymanageropts.json",
		keystorePasswordFileName: "secret",
	}
}

// Get the keystore directory
func (ks *PrysmKeystoreManager) GetKeystoreDir() string {
	return ks.keystoreDir
}

// Store a validator key
func (ks *PrysmKeystoreManager) StoreValidatorKey(key *eth2types.BLSPrivateKey, derivationPath string) error {

	// Initialize the account store
	if err := ks.initialize(); err != nil {
		return err
	}

	// Cancel if validator key already exists in account store
	for ki := 0; ki < len(ks.as.PrivateKeys); ki++ {
		if bytes.Equal(key.Marshal(), ks.as.PrivateKeys[ki]) || bytes.Equal(key.PublicKey().Marshal(), ks.as.PublicKeys[ki]) {
			return nil
		}
	}

	// Add validator key to account store
	ks.as.PrivateKeys = append(ks.as.PrivateKeys, key.Marshal())
	ks.as.PublicKeys = append(ks.as.PublicKeys, key.PublicKey().Marshal())

	// Encode account store
	asBytes, err := json.Marshal(ks.as)
	if err != nil {
		return fmt.Errorf("Could not encode validator account store: %w", err)
	}

	// Get the keystore account password
	passwordFilePath := filepath.Join(ks.keystoreDir, ks.walletDir, ks.accountsDir, ks.keystorePasswordFileName)
	passwordBytes, err := os.ReadFile(passwordFilePath)
	if err != nil {
		return fmt.Errorf("Error reading account password file: %w", err)
	}
	password := string(passwordBytes)

	// Encrypt account store
	asEncrypted, err := ks.encryptor.Encrypt(asBytes, password)
	if err != nil {
		return fmt.Errorf("Could not encrypt validator account store: %w", err)
	}

	// Create new keystore
	keystore := types.ValidatorKeystore{
		Crypto:  asEncrypted,
		Name:    ks.encryptor.Name(),
		Version: ks.encryptor.Version(),
		UUID:    uuid.New(),
	}

	// Encode key store
	ksBytes, err := json.Marshal(keystore)
	if err != nil {
		return fmt.Errorf("Could not encode validator keystore: %w", err)
	}

	// Get file paths
	keystoreFilePath := filepath.Join(ks.keystoreDir, ks.walletDir, ks.accountsDir, ks.keystoreFileName)
	configFilePath := filepath.Join(ks.keystoreDir, ks.walletDir, ks.configFileName)

	// Create keystore dir
	if err := os.MkdirAll(filepath.Dir(keystoreFilePath), DirMode); err != nil {
		return fmt.Errorf("Could not create keystore folder: %w", err)
	}

	// Write keystore to disk
	if err := os.WriteFile(keystoreFilePath, ksBytes, FileMode); err != nil {
		return fmt.Errorf("Could not write keystore to disk: %w", err)
	}

	// Return if wallet config file exists
	if _, err := os.Stat(configFilePath); !os.IsNotExist(err) {
		return nil
	}

	// Create & encode wallet config
	configBytes, err := json.Marshal(prysmWalletConfig{
		DirectEIPVersion: DirectEIPVersion,
	})
	if err != nil {
		return fmt.Errorf("Could not encode wallet config: %w", err)
	}

	// Write wallet config to disk
	if err := os.WriteFile(configFilePath, configBytes, FileMode); err != nil {
		return fmt.Errorf("Could not write wallet config to disk: %w", err)
	}

	// Return
	return nil

}

// Initialize the account store
func (ks *PrysmKeystoreManager) initialize() error {

	// Cancel if already initialized
	if ks.as != nil {
		return nil
	}

	// Create the random keystore password if it doesn't exist
	var password string
	passwordFilePath := filepath.Join(ks.keystoreDir, ks.walletDir, ks.accountsDir, ks.keystorePasswordFileName)
	_, err := os.Stat(passwordFilePath)
	if os.IsNotExist(err) {
		// Create a new password
		password, err = utils.GenerateRandomPassword()
		if err != nil {
			return fmt.Errorf("Could not generate random password: %w", err)
		}

		// Encode it
		passwordBytes := []byte(password)

		// Write it
		err := os.MkdirAll(filepath.Dir(passwordFilePath), DirMode)
		if err != nil {
			return fmt.Errorf("Error creating account password directory: %w", err)
		}
		err = os.WriteFile(passwordFilePath, passwordBytes, FileMode)
		if err != nil {
			return fmt.Errorf("Error writing account password file: %w", err)
		}
	}

	// Get the random keystore password
	passwordBytes, err := os.ReadFile(passwordFilePath)
	if err != nil {
		return fmt.Errorf("Error opening account password file: %w", err)
	}
	password = string(passwordBytes)

	// Read keystore file; initialize empty account store if it doesn't exist
	ksBytes, err := os.ReadFile(filepath.Join(ks.keystoreDir, ks.walletDir, ks.accountsDir, ks.keystoreFileName))
	if err != nil {
		ks.as = &prysmAccountStore{}
		return nil
	}

	// Decode keystore
	keystore := &types.ValidatorKeystore{}
	if err = json.Unmarshal(ksBytes, keystore); err != nil {
		return fmt.Errorf("Could not decode validator keystore: %w", err)
	}

	// Decrypt account store
	asBytes, err := ks.encryptor.Decrypt(keystore.Crypto, password)
	if err != nil {
		return fmt.Errorf("Could not decrypt validator account store: %w", err)
	}

	// Decode account store
	as := &prysmAccountStore{}
	if err = json.Unmarshal(asBytes, as); err != nil {
		return fmt.Errorf("Could not decode validator account store: %w", err)
	}
	if len(as.PrivateKeys) != len(as.PublicKeys) {
		return errors.New("Validator account store private and public key counts do not match")
	}

	// Set account store & return
	ks.as = as
	return nil

}

// Load a private key
func (ks *PrysmKeystoreManager) LoadValidatorKey(pubkey beacon.ValidatorPubkey) (*eth2types.BLSPrivateKey, error) {

	// Initialize the account store
	err := ks.initialize()
	if err != nil {
		return nil, err
	}

	// Find the validator key in the account store
	for ki := 0; ki < len(ks.as.PrivateKeys); ki++ {
		if bytes.Equal(pubkey[:], ks.as.PublicKeys[ki]) {
			decryptedKey := ks.as.PrivateKeys[ki]
			privateKey, err := eth2types.BLSPrivateKeyFromBytes(decryptedKey)
			if err != nil {
				return nil, fmt.Errorf("error recreating private key for validator %s: %w", pubkey.HexWithPrefix(), err)
			}

			// Verify the private key matches the public key
			reconstructedPubkey := beacon.ValidatorPubkey(privateKey.PublicKey().Marshal())
			if reconstructedPubkey != pubkey {
				return nil, fmt.Errorf("Prysm's keystore has a key that claims to be for validator %s but it's for validator %s", pubkey.HexWithPrefix(), reconstructedPubkey.HexWithPrefix())
			}

			return privateKey, nil
		}
	}

	// Return nothing if the private key wasn't found
	return nil, nil

}
