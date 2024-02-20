package swapi

import (
	"github.com/nodeset-org/eth-utils/beacon"
)

type ValidatorStatusData struct {
	ActiveValidators  []beacon.ValidatorPubkey `json:"activePubkeys"`
	ExitingValidators []beacon.ValidatorPubkey `json:"exitPubkeys"`
	ExitedValidators  []beacon.ValidatorPubkey `json:"exitedPubkeys"`
}
