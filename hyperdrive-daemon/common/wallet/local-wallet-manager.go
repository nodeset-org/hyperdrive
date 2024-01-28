package wallet

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethkeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	sharedtypes "github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/tyler-smith/go-bip39"
	eth2ks "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
)

// Simple class to wrap a node's local wallet keystore.
// Note that this does *not* manage the wallet data file on disk, though it does manage the
// legacy keystore used by some integrations.
type LocalWalletManager struct {
	// The ID of the execution layer chain currently being used
	chainID *big.Int

	// Encryptor used to encrypt the account seed
	encryptor *eth2ks.Encryptor

	// Root secret used for EL node wallet derivation and BLS validator key derivation
	// This is the "master" constructed directly from the mnemonic
	seed []byte

	// Derived node wallet private key on the EL
	nodePrivateKey *ecdsa.PrivateKey

	// Serialized data of the loaded wallet
	data *sharedtypes.LocalWalletData

	// Transactor for signing transactions
	transactor *bind.TransactOpts

	// The path to the keystore in Geth's account (v3) format, needed by some projects
	gethKeystorePath string

	// The node's wallet's seed data, serialized in Geth account (v3) format
	ethkey []byte
}

// Creates a new wallet manager for local wallets
func NewLocalWalletManager(legacyKeystorePath string, chainID uint) *LocalWalletManager {
	return &LocalWalletManager{
		chainID:          big.NewInt(int64(chainID)),
		gethKeystorePath: legacyKeystorePath,
		encryptor:        eth2ks.New(),
	}
}

// Get the type of this wallet manager
func (m *LocalWalletManager) GetType() sharedtypes.WalletType {
	return sharedtypes.WalletType_Local
}

// Get the pubkey of the loaded node wallet, or false if one isn't loaded yet
func (m *LocalWalletManager) GetAddress() (common.Address, error) {
	if m.nodePrivateKey == nil {
		return common.Address{}, fmt.Errorf("wallet is not initialized")
	}
	return crypto.PubkeyToAddress(m.nodePrivateKey.PublicKey), nil
}

// Get the private key if it's been loaded
func (m *LocalWalletManager) GetPrivateKey() *ecdsa.PrivateKey {
	return m.nodePrivateKey
}

// Get the legacy keystore in Geth format
func (m *LocalWalletManager) GetEthKeystore() ([]byte, error) {
	if m.ethkey == nil {
		return nil, fmt.Errorf("wallet is not initialized")
	}
	return m.ethkey, nil
}

// Get the transactor for the wallet
func (m *LocalWalletManager) GetTransactor() (*bind.TransactOpts, error) {
	if m.transactor == nil {
		return nil, fmt.Errorf("wallet is not initialized")
	}
	return m.transactor, nil
}

// Initialize a new keystore from a mnemonic and derivation info, derive the corresponding key, and load it all up
func (m *LocalWalletManager) InitializeKeystore(derivationPath string, walletIndex uint, mnemonic string, password string) (*sharedtypes.LocalWalletData, error) {
	// Generate the seed from the mnemonic
	seed := bip39.NewSeed(mnemonic, "")

	// Encrypt the seed with the password
	encryptedSeed, err := m.encryptor.Encrypt(seed, password)
	if err != nil {
		return nil, fmt.Errorf("error encrypting wallet seed: %w", err)
	}

	// Create a new wallet data
	data := &sharedtypes.LocalWalletData{
		Crypto:         encryptedSeed,
		Name:           m.encryptor.Name(),
		Version:        m.encryptor.Version(),
		UUID:           uuid.New(),
		DerivationPath: derivationPath,
		WalletIndex:    walletIndex,
	}

	// Load it
	err = m.LoadWallet(data, password)
	if err != nil {
		return nil, fmt.Errorf("error loading wallet after initialization: %w", err)
	}

	// Get the Geth key
	key := &gethkeystore.Key{
		Address:    crypto.PubkeyToAddress(m.nodePrivateKey.PublicKey),
		PrivateKey: m.nodePrivateKey,
		Id:         uuid.UUID(m.data.UUID),
	}

	// Serialize it
	m.ethkey, err = gethkeystore.EncryptKey(key, password, gethkeystore.StandardScryptN, gethkeystore.StandardScryptP)
	if err != nil {
		return nil, fmt.Errorf("error serializing legacy keystore: %w", err)
	}

	return data, nil
}

// Verifies that the provided password is correct for this wallet's keystore
func (m *LocalWalletManager) VerifyPassword(password string) (bool, error) {
	if m.data == nil {
		return false, fmt.Errorf("wallet is not initialized")
	}

	// Make a new local manager and load the data with the candidate password
	candidateMgr := NewLocalWalletManager("", 0)
	err := candidateMgr.LoadWallet(m.data, password)
	if err != nil {
		return false, fmt.Errorf("error verifying wallet with candidate password: %w", err)
	}
	trueBytes := crypto.FromECDSA(m.nodePrivateKey)
	candidateBytes := crypto.FromECDSA(candidateMgr.nodePrivateKey)

	return bytes.Equal(trueBytes, candidateBytes), nil
}

// Saves the legacy Geth keystore to disk
func (m *LocalWalletManager) SaveKeystore() error {
	// Write it to file
	err := os.WriteFile(m.gethKeystorePath, m.ethkey, walletFileMode)
	if err != nil {
		return fmt.Errorf("error writing legacy keystore to [%s]: %w", m.gethKeystorePath, err)
	}
	return nil
}

// Load the node wallet's private key from the keystore
func (m *LocalWalletManager) LoadWallet(data *sharedtypes.LocalWalletData, password string) error {
	// Decrypt the seed
	var err error
	seed, err := m.encryptor.Decrypt(data.Crypto, password)
	if err != nil {
		return fmt.Errorf("error decrypting wallet keystore: %w", err)
	}

	// Create the master key
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return fmt.Errorf("error creating wallet master key: %w", err)
	}

	// Handle an empty derivation path
	if data.DerivationPath == "" {
		data.DerivationPath = DefaultNodeKeyPath
	}

	// Get the derived key
	derivedKey, index, err := getDerivedKey(masterKey, data.DerivationPath, data.WalletIndex)
	if err != nil {
		return fmt.Errorf("error getting node wallet derived key: %w", err)
	}
	data.WalletIndex = index // Update the index in case of the ErrInvalidChild issue

	// Get the private key from it
	privateKey, err := derivedKey.ECPrivKey()
	if err != nil {
		return fmt.Errorf("error getting node wallet private key: %w", err)
	}
	privateKeyECDSA := privateKey.ToECDSA()

	// Make a transactor from it
	transactor, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, m.chainID)
	if err != nil {
		return fmt.Errorf("error creating transactor for node private key: %w", err)
	}
	transactor.Context = context.Background()

	// Store everything if there are no errors
	m.seed = seed
	m.nodePrivateKey = privateKeyECDSA
	m.data = data
	m.transactor = transactor
	return nil
}

// Signs a message with the node wallet's private key
func (m *LocalWalletManager) SignMessage(message []byte) ([]byte, error) {
	messageHash := accounts.TextHash(message)
	signedMessage, err := crypto.Sign(messageHash, m.nodePrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Error signing message: %w", err)
	}

	// fix the ECDSA 'v' (see https://medium.com/mycrypto/the-magic-of-digital-signatures-on-ethereum-98fe184dc9c7#:~:text=The%20version%20number,2%E2%80%9D%20was%20introduced)
	signedMessage[crypto.RecoveryIDOffset] += 27
	return signedMessage, nil
}

// Signs a transaction with the node wallet's private key
func (m *LocalWalletManager) SignTransaction(serializedTx []byte) ([]byte, error) {
	tx := types.Transaction{}
	err := tx.UnmarshalBinary(serializedTx)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling TX: %w", err)
	}

	signer := types.NewLondonSigner(m.chainID)
	signedTx, err := types.SignTx(&tx, signer, m.nodePrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error signing TX: %w", err)
	}

	signedData, err := signedTx.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("error marshalling signed TX to binary: %w", err)
	}

	return signedData, nil
}

// Serialize the wallet data as JSON
func (m *LocalWalletManager) SerializeData() (string, error) {
	if m.data == nil {
		return "", fmt.Errorf("wallet is not initialized")
	}

	bytes, err := json.Marshal(m.data)
	if err != nil {
		return "", fmt.Errorf("error serializing wallet data: %w", err)
	}
	return string(bytes), nil
}

// Get the derived key & derivation path for the account at the index
func getDerivedKey(masterKey *hdkeychain.ExtendedKey, derivationPath string, index uint) (*hdkeychain.ExtendedKey, uint, error) {
	formattedDerivationPath := fmt.Sprintf(derivationPath, index)

	// Parse derivation path
	path, err := accounts.ParseDerivationPath(formattedDerivationPath)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid node key derivation path '%s': %w", formattedDerivationPath, err)
	}

	// Follow derivation path
	key := masterKey
	for i, n := range path {
		key, err = key.Derive(n)
		if err == hdkeychain.ErrInvalidChild {
			// Start over with the next index
			return getDerivedKey(masterKey, derivationPath, index+1)
		} else if err != nil {
			return nil, 0, fmt.Errorf("invalid child key at depth %d: %w", i, err)
		}
	}

	// Return
	return key, index, nil
}
