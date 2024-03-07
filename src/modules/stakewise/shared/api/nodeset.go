package swapi

import (
	nmc_beacon "github.com/rocket-pool/node-manager-core/beacon"
)

type NodesetUploadDepositDataData struct {
	ServerResponse []byte                       `json:"serverResponse"`
	NewPubkeys     []nmc_beacon.ValidatorPubkey `json:"newPubkeys"`
	TotalCount     uint64                       `json:"totalCount"`
}
