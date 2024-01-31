package keystore

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/nodeset-org/eth-utils/beacon"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
	eth2ks "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"

	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/nodeset-org/hyperdrive/shared/utils"
)

// Teku keystore manager
type TekuKeystoreManager struct {
	keystorePath  string
	encryptor     *eth2ks.Encryptor
	keystoreDir   string
	secretsDir    string
	validatorsDir string
}

// Create new teku keystore manager
func NewTekuKeystoreManager(keystorePath string) *TekuKeystoreManager {
	return &TekuKeystoreManager{
		keystorePath:  keystorePath,
		encryptor:     eth2ks.New(eth2ks.WithCipher("scrypt")),
		keystoreDir:   "teku",
		secretsDir:    "passwords",
		validatorsDir: "keys",
	}
}

// Get the keystore directory
func (ks *TekuKeystoreManager) GetKeystoreDir() string {
	return filepath.Join(ks.keystorePath, ks.keystoreDir)
}

// Store a validator key
func (ks *TekuKeystoreManager) StoreValidatorKey(key *eth2types.BLSPrivateKey, derivationPath string) error {

	// Get validator pubkey
	pubkey := beacon.ValidatorPubkey(key.PublicKey().Marshal())

	// Create a new password
	password, err := utils.GenerateRandomPassword()
	if err != nil {
		return fmt.Errorf("Could not generate random password: %w", err)
	}

	// Encrypt key
	encryptedKey, err := ks.encryptor.Encrypt(key.Marshal(), password)
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

	// Get secret file path
	secretFilePath := filepath.Join(ks.keystorePath, ks.keystoreDir, ks.secretsDir, pubkey.HexWithPrefix()+".txt")

	// Create secrets dir
	if err := os.MkdirAll(filepath.Dir(secretFilePath), DirMode); err != nil {
		return fmt.Errorf("Could not create validator secrets folder: %w", err)
	}

	// Write secret to disk
	if err := os.WriteFile(secretFilePath, []byte(password), FileMode); err != nil {
		return fmt.Errorf("Could not write validator secret to disk: %w", err)
	}

	// Get key file path
	keyFilePath := filepath.Join(ks.keystorePath, ks.keystoreDir, ks.validatorsDir, pubkey.HexWithPrefix()+".json")

	// Create key dir
	if err := os.MkdirAll(filepath.Dir(keyFilePath), DirMode); err != nil {
		return fmt.Errorf("Could not create validator key folder: %w", err)
	}

	// Write key store to disk
	if err := os.WriteFile(keyFilePath, keyStoreBytes, FileMode); err != nil {
		return fmt.Errorf("Could not write validator key to disk: %w", err)
	}

	// Return
	return nil

}

// Load a private key
func (ks *TekuKeystoreManager) LoadValidatorKey(pubkey beacon.ValidatorPubkey) (*eth2types.BLSPrivateKey, error) {

	// Get key file path
	keyFilePath := filepath.Join(ks.keystorePath, ks.keystoreDir, ks.validatorsDir, pubkey.HexWithPrefix()+".json")

	// Read the key file
	_, err := os.Stat(keyFilePath)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("couldn't open the Teku keystore for pubkey %s: %w", pubkey.HexWithPrefix(), err)
	}
	bytes, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read the Teku keystore for pubkey %s: %w", pubkey.HexWithPrefix(), err)
	}

	// Unmarshal the keystore
	var keystore types.ValidatorKeystore
	err = json.Unmarshal(bytes, &keystore)
	if err != nil {
		return nil, fmt.Errorf("error deserializing Teku keystore for pubkey %s: %w", pubkey.HexWithPrefix(), err)
	}

	// Get secret file path
	secretFilePath := filepath.Join(ks.keystorePath, ks.keystoreDir, ks.secretsDir, pubkey.HexWithPrefix()+".txt")

	// Read secret from disk
	_, err = os.Stat(secretFilePath)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("couldn't open the Teku secret for pubkey %s: %w", pubkey.HexWithPrefix(), err)
	}
	bytes, err = os.ReadFile(secretFilePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read the Teku secret for pubkey %s: %w", pubkey.HexWithPrefix(), err)
	}

	// Decrypt key
	password := string(bytes)
	decryptedKey, err := ks.encryptor.Decrypt(keystore.Crypto, password)
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
