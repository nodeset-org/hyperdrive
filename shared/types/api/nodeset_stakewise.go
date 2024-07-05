package api

import (
	apiv1 "github.com/nodeset-org/nodeset-client-go/api-v1"
	"github.com/rocket-pool/node-manager-core/beacon"
)

type NodeSetStakeWise_GetRegisteredValidatorsData struct {
	NotRegistered bool                    `json:"notRegistered"`
	Validators    []apiv1.ValidatorStatus `json:"validators"`
}

type NodeSetStakeWise_GetDepositDataSetData struct {
	NotRegistered bool                         `json:"notRegistered"`
	Version       int                          `json:"version"`
	DepositData   []beacon.ExtendedDepositData `json:"depositData"`
}

type NodeSetStakeWise_UploadDepositDataData struct {
	NotRegistered      bool `json:"notRegistered"`
	VaultNotFound      bool `json:"vaultNotFound"`
	InvalidPermissions bool `json:"invalidPermissions"`
}

type NodeSetStakeWise_UploadSignedExitsData struct {
	NotRegistered bool `json:"notRegistered"`
}
