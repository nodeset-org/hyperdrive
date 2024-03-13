package swapi

import (
	"math/big"

	"github.com/rocket-pool/node-manager-core/beacon"
)

type NodesetUploadDepositDataData struct {
	SufficientBalance    bool                     `json:"sufficientBalance"`
	Balance              *big.Int                 `json:"balance"`
	RequiredBalance      *big.Int                 `json:"requiredBalance"`
	UnregisteredKeyCount int                      `json:"unregisteredKeyCount"`
	ServerResponse       []byte                   `json:"serverResponse"`
	NewPubkeys           []beacon.ValidatorPubkey `json:"newPubkeys"`
	TotalCount           uint64                   `json:"totalCount"`
}
