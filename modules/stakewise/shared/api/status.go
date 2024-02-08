package swapi

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/eth-utils/beacon"
)

type StatusActiveWalletsResponse struct {
	ActiveWallets []common.Address `json:"accountAddresses"`
}

type ActiveValidatorsData struct {
	ActiveValidators []beacon.ValidatorPubkey `json:"pubkeys"`
}
