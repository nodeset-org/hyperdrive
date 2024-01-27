package api

import "github.com/ethereum/go-ethereum/common"

type UtilsResolveEnsData struct {
	Address       common.Address `json:"address"`
	EnsName       string         `json:"ensName"`
	FormattedName string         `json:"formattedName"`
}
