package api

import (
	apiv1 "github.com/nodeset-org/nodeset-client-go/api-v1"
	"github.com/rocket-pool/node-manager-core/beacon"
)

type NodeSetStakeWise_GetRegisteredValidatorsData struct {
	Validators []apiv1.ValidatorStatus `json:"validators"`
}

type NodeSetStakeWise_GetDepositDataSetData struct {
	Version     int                          `json:"version"`
	DepositData []beacon.ExtendedDepositData `json:"depositData"`
}
