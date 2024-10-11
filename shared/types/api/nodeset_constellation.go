package api

import (
	"github.com/ethereum/go-ethereum/common"
	v2constellation "github.com/nodeset-org/nodeset-client-go/api-v2/constellation"
	nscommon "github.com/nodeset-org/nodeset-client-go/common"
)

type NodeSetConstellation_GetRegisteredAddressData struct {
	NotRegisteredWithNodeSet       bool           `json:"notRegisteredWithNodeSet"`
	NotRegisteredWithConstellation bool           `json:"notRegisteredWithConstellation"`
	InvalidPermissions             bool           `json:"invalidPermissions"`
	RegisteredAddress              common.Address `json:"registeredAddress"`
}

type NodeSetConstellation_GetRegistrationSignatureData struct {
	NotRegistered        bool   `json:"notRegistered"`
	NotAuthorized        bool   `json:"notAuthorized"`
	InvalidPermissions   bool   `json:"invalidPermissions"`
	IncorrectNodeAddress bool   `json:"incorrectNodeAddress"`
	Signature            []byte `json:"signature"`
}

type NodeSetConstellation_GetDepositSignatureData struct {
	NotRegistered            bool   `json:"notRegistered"`
	NotWhitelisted           bool   `json:"notWhitelisted"`
	IncorrectNodeAddress     bool   `json:"incorrectNodeAddress"`
	LimitReached             bool   `json:"limitReached"`
	MissingExitMessage       bool   `json:"missingExitMessage"`
	AddressAlreadyRegistered bool   `json:"addressAlreadyRegistered"`
	InvalidPermissions       bool   `json:"invalidPermissions"`
	Signature                []byte `json:"signature"`
}

type NodeSetConstellation_GetValidatorsData struct {
	NotRegistered        bool                              `json:"notRegistered"`
	NotWhitelisted       bool                              `json:"notWhitelisted"`
	IncorrectNodeAddress bool                              `json:"incorrectNodeAddress"`
	InvalidPermissions   bool                              `json:"invalidPermissions"`
	Validators           []v2constellation.ValidatorStatus `json:"validators"`
}

type NodeSetConstellation_UploadSignedExitsRequestBody struct {
	Deployment   string                       `json:"deployment"`
	ExitMessages []nscommon.EncryptedExitData `json:"exitMessages"`
}

type NodeSetConstellation_UploadSignedExitsData struct {
	NotRegistered            bool `json:"notRegistered"`
	NotWhitelisted           bool `json:"notWhitelisted"`
	IncorrectNodeAddress     bool `json:"incorrectNodeAddress"`
	InvalidValidatorOwner    bool `json:"invalidValidatorOwner"`
	ExitMessageAlreadyExists bool `json:"exitMessageAlreadyExists"`
	InvalidExitMessage       bool `json:"invalidExitMessage"`
	InvalidPermissions       bool `json:"invalidPermissions"`
}
