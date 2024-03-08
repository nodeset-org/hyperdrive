package swapi

import (
	"github.com/rocket-pool/node-manager-core/beacon"
)

type ActiveValidatorsData struct {
	ActiveValidators []beacon.ValidatorPubkey `json:"pubkeys"`
}
