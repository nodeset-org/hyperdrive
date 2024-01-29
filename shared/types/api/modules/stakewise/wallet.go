package api

import "github.com/ethereum/go-ethereum/common"

type WalletInitializeData struct {
	AccountAddress common.Address `json:"accountAddress"`
}
