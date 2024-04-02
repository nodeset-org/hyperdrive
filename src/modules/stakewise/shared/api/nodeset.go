package swapi

import (
	"math/big"

	"github.com/rocket-pool/node-manager-core/beacon"
)

type NodesetUploadDepositDataData struct {
	SufficientBalance   bool                     `json:"sufficientBalance"`
	Balance             *big.Int                 `json:"balance"`
	RequiredBalance     *big.Int                 `json:"requiredBalance"`
	ServerResponse      []byte                   `json:"serverResponse"`
	UnregisteredPubkeys []beacon.ValidatorPubkey `json:"newPubkeys"`
	TotalCount          uint64                   `json:"totalCount"`
}
