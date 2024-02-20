package swapi

import (
	"github.com/nodeset-org/eth-utils/beacon"
)

type ValidatorStatusData struct {
	ActiveValidators []beacon.ValidatorPubkey `json:"pubkeys"`
}
