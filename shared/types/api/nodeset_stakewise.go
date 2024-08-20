package api

import (
	"github.com/nodeset-org/nodeset-client-go/common/stakewise"
	"github.com/rocket-pool/node-manager-core/beacon"
)

type NodeSetStakeWise_GetRegisteredValidatorsData struct {
	NotRegistered bool                        `json:"notRegistered"`
	Validators    []stakewise.ValidatorStatus `json:"validators"`
}

type NodeSetStakeWise_GetDepositDataSetVersionData struct {
	NotRegistered bool `json:"notRegistered"`
	Version       int  `json:"version"`
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
