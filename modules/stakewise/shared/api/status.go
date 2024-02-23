package swapi

import (
	"github.com/nodeset-org/eth-utils/beacon"
)

type ValidatorStatusData struct {
	Generated             []beacon.ValidatorPubkey `json:"generatedPubkeys"`
	UploadedToNodeset     []beacon.ValidatorPubkey `json:"uploadedNodesetPubkeys"`
	UploadToStakewise     []beacon.ValidatorPubkey `json:"uploadedStakewisePubkeys"`
	RegisteredToStakewise []beacon.ValidatorPubkey `json:"registeredStakewisePubkeys"`
	PendingInitialized    []beacon.ValidatorPubkey `json:"pendingInitializedPubkeys"`
	PendingQueued         []beacon.ValidatorPubkey `json:"pendingQueuedPubkeys"`
	ActiveOngoing         []beacon.ValidatorPubkey `json:"activeOngoingPubkeys"`
	ActiveExited          []beacon.ValidatorPubkey `json:"activeExitedPubkeys"`
	ActiveSlashed         []beacon.ValidatorPubkey `json:"activeSlashedPubkeys"`
	ExitedUnslashed       []beacon.ValidatorPubkey `json:"exitedUnslashedPubkeys"`
	ExitedSlashed         []beacon.ValidatorPubkey `json:"exitedSlashedPubkeys"`
	WithdrawalPossible    []beacon.ValidatorPubkey `json:"withdrawalPossiblePubkeys"`
	WithdrawalDone        []beacon.ValidatorPubkey `json:"withdrawalDonePubkeys"`
}
