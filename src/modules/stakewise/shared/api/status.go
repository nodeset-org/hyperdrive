package swapi

import (
	"github.com/nodeset-org/eth-utils/beacon"
)

type ValidatorStatus string

const (
	WithdrawalDone        ValidatorStatus = "WithdrawalDone"
	WithdrawalPossible    ValidatorStatus = "WithdrawalPossible"
	ExitedSlashed         ValidatorStatus = "ExitedSlashed"
	ExitedUnslashed       ValidatorStatus = "ExitedUnslashed"
	ActiveSlashed         ValidatorStatus = "ActiveSlashed"
	ActiveExited          ValidatorStatus = "ActiveExited"
	ActiveOngoing         ValidatorStatus = "ActiveOngoing"
	PendingQueued         ValidatorStatus = "PendingQueued"
	PendingInitialized    ValidatorStatus = "PendingInitialized"
	RegisteredToStakewise ValidatorStatus = "RegisteredToStakewise"
	UploadedStakewise     ValidatorStatus = "UploadedStakewise"
	UploadedToNodeset     ValidatorStatus = "UploadedToNodeset"
	Generated             ValidatorStatus = "Generated"
)

type ValidatorStatusData struct {
	ValidatorStatus map[beacon.ValidatorPubkey]ValidatorStatus `json:"validatorStatus"`
}
