package swapi

import (
	"github.com/nodeset-org/eth-utils/beacon"
	swtypes "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/types"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

type ValidatorStateInfo struct {
	Pubkey        beacon.ValidatorPubkey `json:"pubkey"`
	Index         string                 `json:"index"`
	BeaconStatus  types.ValidatorState   `json:"beaconStatus"`
	NodesetStatus swtypes.NodesetStatus  `json:"nodesetStatus"`
}

type ValidatorStatusData struct {
	States []ValidatorStateInfo `json:"states"`
}
