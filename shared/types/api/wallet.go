package api

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/node-manager-core/wallet"
)

type WalletStatusData struct {
	WalletStatus wallet.WalletStatus `json:"walletStatus"`
}

type WalletBalanceData struct {
	Balance *big.Int `json:"balance"`
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
	SignedMessage []byte `json:"signedMessage"`
}

type WalletSignTxData struct {
	SignedTx []byte `json:"signedTx"`
}

type WalletExportEthKeyData struct {
	EthKeyJson []byte `json:"ethKeyJson"`
	Password   string `json:"password"`
}

type WalletGenerateValidatorKeyData struct {
	PrivateKey []byte `json:"privateKey"`
}

type WalletSendData struct {
	Balance             *big.Int             `json:"balance"`
	TokenName           string               `json:"name"`
	TokenSymbol         string               `json:"symbol"`
	CanSend             bool                 `json:"canSend"`
	InsufficientBalance bool                 `json:"insufficientBalance"`
	TxInfo              *eth.TransactionInfo `json:"txInfo"`
}
