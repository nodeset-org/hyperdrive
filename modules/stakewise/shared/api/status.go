package swapi

import (
	"github.com/ethereum/go-ethereum/common"
)

type StatusActiveWalletsResponse struct {
	ActiveWallets []common.Address `json:"accountAddresses"`
}

type StatusActiveWalletsData struct {
}
