package api

import (
	v2constellation "github.com/nodeset-org/nodeset-client-go/api-v2/constellation"
	nscommon "github.com/nodeset-org/nodeset-client-go/common"
)

type NodeSetConstellation_GetRegistrationSignatureData struct {
	NotRegistered bool   `json:"notRegistered"`
	NotAuthorized bool   `json:"notAuthorized"`
	Signature     []byte `json:"signature"`
}

type NodeSetConstellation_GetDepositSignatureData struct {
	NotRegistered      bool   `json:"notRegistered"`
	NotAuthorized      bool   `json:"notAuthorized"`
	LimitReached       bool   `json:"limitReached"`
	MissingExitMessage bool   `json:"missingExitMessage"`
	Signature          []byte `json:"signature"`
}

type NodeSetConstellation_GetValidatorsData struct {
	NotRegistered bool                              `json:"notRegistered"`
	NotAuthorized bool                              `json:"notAuthorized"`
	Validators    []v2constellation.ValidatorStatus `json:"validators"`
}

type NodeSetConstellation_UploadSignedExitsRequestBody struct {
	ExitMessages []nscommon.ExitData `json:"exitMessages"`
}

type NodeSetConstellation_UploadSignedExitsData struct {
	NotRegistered bool `json:"notRegistered"`
	NotAuthorized bool `json:"notAuthorized"`
}
