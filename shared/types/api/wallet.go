package api

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/eth-utils/eth"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

type WalletStatusData struct {
	WalletStatus types.WalletStatus `json:"walletStatus"`
}

type WalletInitializeData struct {
	Mnemonic       string         `json:"mnemonic"`
	AccountAddress common.Address `json:"accountAddress"`
}

type WalletRecoverData struct {
	AccountAddress common.Address           `json:"accountAddress"`
	ValidatorKeys  []beacon.ValidatorPubkey `json:"validatorKeys"`
}

type WalletSearchAndRecoverData struct {
	FoundWallet    bool                     `json:"foundWallet"`
	AccountAddress common.Address           `json:"accountAddress"`
	DerivationPath string                   `json:"derivationPath"`
	Index          uint                     `json:"index"`
	ValidatorKeys  []beacon.ValidatorPubkey `json:"validatorKeys"`
}

type WalletRebuildData struct {
	ValidatorKeys []beacon.ValidatorPubkey `json:"validatorKeys"`
}

type WalletExportData struct {
	Password          string `json:"password"`
	Wallet            string `json:"wallet"`
	AccountPrivateKey []byte `json:"accountPrivateKey"`
}

type WalletSetEnsNameData struct {
	Address common.Address       `json:"address"`
	EnsName string               `json:"ensName"`
	TxInfo  *eth.TransactionInfo `json:"txInfo"`
}

type WalletTestMnemonicData struct {
	CurrentAddress   common.Address `json:"currentAddress"`
	RecoveredAddress common.Address `json:"recoveredAddress"`
}

type WalletSignMessageData struct {
	SignedMessage string `json:"signedMessage"`
}

type WalletExportEthKeyData struct {
	EthKeyJson []byte `json:"ethKeyJson"`
}

type WalletGenerateValidatorKeyData struct {
	PrivateKey []byte `json:"privateKey"`
}
