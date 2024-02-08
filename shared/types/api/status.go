package api

import "github.com/nodeset-org/eth-utils/beacon"

type ActiveValidatorsData struct {
	ActiveValidators []beacon.ValidatorPubkey `json:"active_validators"`
}
