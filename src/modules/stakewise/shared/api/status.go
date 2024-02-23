package swapi

import (
	"github.com/nodeset-org/eth-utils/beacon"
)

type ActiveValidatorsData struct {
	ActiveValidators []beacon.ValidatorPubkey `json:"pubkeys"`
}
