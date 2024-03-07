package swapi

import (
	"github.com/ethereum/go-ethereum/common"
	nmc_beacon "github.com/rocket-pool/node-manager-core/beacon"
)

type WalletInitializeData struct {
	AccountAddress common.Address `json:"accountAddress"`
}

type WalletGenerateKeysData struct {
	Pubkeys []nmc_beacon.ValidatorPubkey `json:"pubkeys"`
}
