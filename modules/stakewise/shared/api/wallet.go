package swapi

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/node-manager-core/beacon"
)

type WalletInitializeData struct {
	AccountAddress common.Address `json:"accountAddress"`
}

type WalletGenerateKeysData struct {
	Pubkeys []beacon.ValidatorPubkey `json:"pubkeys"`
}
