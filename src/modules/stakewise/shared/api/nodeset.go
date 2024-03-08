package swapi

import (
	"github.com/rocket-pool/node-manager-core/beacon"
)

type NodesetUploadDepositDataData struct {
	ServerResponse []byte                   `json:"serverResponse"`
	NewPubkeys     []beacon.ValidatorPubkey `json:"newPubkeys"`
	TotalCount     uint64                   `json:"totalCount"`
}
