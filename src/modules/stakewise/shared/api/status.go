package swapi

import (
	swtypes "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/types"
	"github.com/rocket-pool/node-manager-core/beacon"
)

type ValidatorStateInfo struct {
	Pubkey        beacon.ValidatorPubkey `json:"pubkey"`
	Index         string                 `json:"index"`
	BeaconStatus  beacon.ValidatorState  `json:"beaconStatus"`
	NodesetStatus swtypes.NodesetStatus  `json:"nodesetStatus"`
}

type ValidatorStatusData struct {
	States []ValidatorStateInfo `json:"states"`
}
