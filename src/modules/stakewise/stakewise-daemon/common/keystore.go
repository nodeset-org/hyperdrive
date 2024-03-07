package swcommon

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/nodeset-org/eth-utils/beacon"
	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/types"
	nmc_utils "github.com/rocket-pool/node-manager-core/utils"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
	eth2ks "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
)

const (
	keystorePrefix string = "keystore-"
	keystoreSuffix string = ".json"
)

// Keystore manager for the Stakewise operator
type stakewiseKeystoreManager struct {
	encryptor   *eth2ks.Encryptor
	keystoreDir string
	password    string
}

// Create new Stakewise keystore manager
func newStakewiseKeystoreManager(moduleDir string) (*stakewiseKeystoreManager, error) {
	keystoreDir := filepath.Join(moduleDir, config.ValidatorsDirectory, swconfig.ModuleName)
	passwordPath := filepath.Join(keystoreDir, swconfig.KeystorePasswordFile)

	// Read the password file
	var password string
	_, err := os.Stat(passwordPath)
	if errors.Is(err, fs.ErrNotExist) {
		// Make a new one
		password, err = initializeKeystorePassword(passwordPath)
		if err != nil {
			return nil, fmt.Errorf("error generating initial random validator keystore password: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("error reading keystore password from [%s]: %w", passwordPath, err)
	} else {
		bytes, err := os.ReadFile(passwordPath)
		if err != nil {
			return nil, fmt.Errorf("error reading keystore password from [%s]: %w", passwordPath, err)
		}
		password = string(bytes)
	}

	return &stakewiseKeystoreManager{
		encryptor:   eth2ks.New(eth2ks.WithCipher("scrypt")),
		keystoreDir: keystoreDir,
		password:    password,
	}, nil
}

// Get the keystore directory
func (ks *stakewiseKeystoreManager) GetKeystoreDir() string {
	return ks.keystoreDir
}

// Store a validator key
func (ks *stakewiseKeystoreManager) StoreValidatorKey(key *eth2types.BLSPrivateKey, derivationPath string) error {
	// Get validator pubkey
	pubkey := beacon.ValidatorPubkey(key.PublicKey().Marshal())

	// Encrypt key
	encryptedKey, err := ks.encryptor.Encrypt(key.Marshal(), ks.password)
	if err != nil {
		return fmt.Errorf("Could not encrypt validator key: %w", err)
	}

	// Create key store
	keyStore := types.ValidatorKeystore{
		Crypto:  encryptedKey,
		Version: ks.encryptor.Version(),
		UUID:    uuid.New(),
		Path:    derivationPath,
		Pubkey:  pubkey,
	}

	// Encode key store
	keyStoreBytes, err := json.Marshal(keyStore)
	if err != nil {
		return fmt.Errorf("Could not encode validator key: %w", err)
	}

	// Get key file path
	keyFilePath := filepath.Join(ks.keystoreDir, keystorePrefix+pubkey.HexWithPrefix()+keystoreSuffix)

	// Write key store to disk
	if err := os.WriteFile(keyFilePath, keyStoreBytes, fileMode); err != nil {
		return fmt.Errorf("Could not write validator key to disk: %w", err)
	}

	// Return
	return nil
}

// Load a private key
func (ks *stakewiseKeystoreManager) LoadValidatorKey(pubkey beacon.ValidatorPubkey) (*eth2types.BLSPrivateKey, error) {
	// Get key file path
	keyFilePath := filepath.Join(ks.keystoreDir, keystorePrefix+pubkey.HexWithPrefix()+keystoreSuffix)

	// Read the key file
	_, err := os.Stat(keyFilePath)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("couldn't open the Stakewise keystore for pubkey %s: %w", pubkey.HexWithPrefix(), err)
	}
	bytes, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read the Stakewise keystore for pubkey %s: %w", pubkey.HexWithPrefix(), err)
	}

	// Unmarshal the keystore
	var keystore types.ValidatorKeystore
	err = json.Unmarshal(bytes, &keystore)
	if err != nil {
		return nil, fmt.Errorf("error deserializing Stakewise keystore for pubkey %s: %w", pubkey.HexWithPrefix(), err)
	}

	// Decrypt key
	decryptedKey, err := ks.encryptor.Decrypt(keystore.Crypto, ks.password)
	if err != nil {
		return nil, fmt.Errorf("couldn't decrypt keystore for pubkey %s: %w", pubkey.HexWithPrefix(), err)
	}
	privateKey, err := eth2types.BLSPrivateKeyFromBytes(decryptedKey)
	if err != nil {
		return nil, fmt.Errorf("error recreating private key for validator %s: %w", keystore.Pubkey.HexWithPrefix(), err)
	}

	// Verify the private key matches the public key
	reconstructedPubkey := beacon.ValidatorPubkey(privateKey.PublicKey().Marshal())
	if reconstructedPubkey != pubkey {
		return nil, fmt.Errorf("private keystore file %s claims to be for validator %s but it's for validator %s", keyFilePath, pubkey.HexWithPrefix(), reconstructedPubkey.HexWithPrefix())
	}

	return privateKey, nil
}

// Initializes the Stakewise keystore directory and saves a random password to it
func initializeKeystorePassword(passwordPath string) (string, error) {
	// Make a password
	password, err := nmc_utils.GenerateRandomPassword()
	if err != nil {
		return "", err
	}

	// Make the keystore dir
	keystoreDir := filepath.Dir(passwordPath)
	err = os.MkdirAll(keystoreDir, dirMode)
	if err != nil {
		return "", fmt.Errorf("error creating keystore directory [%s]: %w", keystoreDir, err)
	}

	err = os.WriteFile(passwordPath, []byte(password), fileMode)
	if err != nil {
		return "", fmt.Errorf("error saving password to file [%s]: %w", passwordPath, err)
	}
	return password, nil
}
