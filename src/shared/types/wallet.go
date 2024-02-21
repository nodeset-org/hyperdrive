package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type WalletStatus struct {
	Address struct {
		NodeAddress common.Address `json:"nodeAddress"`
		HasAddress  bool           `json:"hasAddress"`
	} `json:"address"`

	Wallet struct {
		Type          WalletType     `json:"type"`
		IsLoaded      bool           `json:"isLoaded"`
		IsOnDisk      bool           `json:"isOnDisk"`
		WalletAddress common.Address `json:"walletAddress"`
	} `json:"wallet"`

	Password struct {
		IsPasswordSaved bool `json:"isPasswordSaved"`
	} `json:"password"`
}
type DerivationPath string

const (
	DerivationPath_Default    DerivationPath = ""
	DerivationPath_LedgerLive DerivationPath = "ledger-live"
	DerivationPath_Mew        DerivationPath = "mew"
)

// An enum describing the type of wallet used by the node
type WalletType string

const (
	// Unset wallet type
	WalletType_Unknown WalletType = ""

	// Indicator for local wallets that have encrypted keystores saved to disk
	WalletType_Local WalletType = "local"

	// Indicator for hardware wallets that store the private key offline
	WalletType_Hardware WalletType = "hardware"
)

// Keystore for local node wallets - note that this is NOT an EIP-2335 keystore.
type LocalWalletData struct {
	// Encrypted seed information
	Crypto map[string]interface{} `json:"crypto"`

	// Name of the encryptor used to generate the keystore
	Name string `json:"name"`

	// Version of the encryptor used to generate the keystore
	Version uint `json:"version"`

	// Unique ID for this keystore
	UUID uuid.UUID `json:"uuid"`

	// The path that should be used to derive the target key; assumes there's only one index that can be iterated on
	DerivationPath string `json:"derivationPath,omitempty"`

	// The index of the target wallet, used to format DerivationPath
	WalletIndex uint `json:"walletIndex,omitempty"`
}

// Placeholder for hardware wallets
type HardwareWalletData struct {
	// NYI
}

// Data storage for node wallets
type WalletData struct {
	// The type of wallet
	Type WalletType `json:"type"`

	// Data about a local wallet
	LocalData LocalWalletData `json:"localData"`

	// Data about a hardware wallet
	HardwareData HardwareWalletData `json:"hardwareData"`
}
