package wallet

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

type IWalletManager interface {
	GetType() types.WalletType
	GetAddress() (common.Address, error)
	GetTransactor() (*bind.TransactOpts, error)
	SignMessage(message []byte) ([]byte, error)
	SignTransaction(serializedTx []byte) ([]byte, error)
}
