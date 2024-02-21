package swapi

import (
	"github.com/nodeset-org/eth-utils/beacon"
)

type ValidatorStatusData struct {
	Generated                  []beacon.ValidatorPubkey `json:"generatedPubkeys"`
	UploadedToNodeset          []beacon.ValidatorPubkey `json:"uploadedNodesetPubkeys"`
	UploadToStakewise          []beacon.ValidatorPubkey `json:"uploadedStakewisePubkeys"`
	RegisteredToStakewise      []beacon.ValidatorPubkey `json:"registeredStakewisePubkeys"`
	WaitingDepositConfirmation []beacon.ValidatorPubkey `json:"waitingDepositConfirmationPubkeys"`
	Depositing                 []beacon.ValidatorPubkey `json:"depositingPubkeys"`
	Deposited                  []beacon.ValidatorPubkey `json:"depositedPubkeys"`
	Active                     []beacon.ValidatorPubkey `json:"activePubkeys"`
	Exiting                    []beacon.ValidatorPubkey `json:"exitingPubkeys"`
	Exited                     []beacon.ValidatorPubkey `json:"exitedPubkeys"`
}
