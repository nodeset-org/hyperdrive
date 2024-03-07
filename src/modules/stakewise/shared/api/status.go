package swapi

import (
	nmc_beacon "github.com/rocket-pool/node-manager-core/beacon"
)

type ActiveValidatorsData struct {
	ActiveValidators []nmc_beacon.ValidatorPubkey `json:"pubkeys"`
}
