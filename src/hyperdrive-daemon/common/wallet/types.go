package wallet

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

// Interface for wallet managers
type IWalletManager interface {
	// The type of wallet
	GetType() types.WalletType

	// The wallet address
	GetAddress() (common.Address, error)

	// A transactor that can sign transactions
	GetTransactor() (*bind.TransactOpts, error)

	// Sign a message with the wallet's private key
	SignMessage(message []byte) ([]byte, error)

	// Sign a transaction with the wallet's private key
	SignTransaction(serializedTx []byte) ([]byte, error)

	// Serialize the wallet data as JSON
	SerializeData() (string, error)
}
