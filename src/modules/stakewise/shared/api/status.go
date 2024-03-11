package swapi

import (
	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

type NodesetStatus string

const (
	RegisteredToStakewise NodesetStatus = "RegisteredToStakewise"
	UploadedStakewise     NodesetStatus = "UploadedStakewise"
	UploadedToNodeset     NodesetStatus = "UploadedToNodeset"
	Generated             NodesetStatus = "Generated"
)

type ValidatorStatusData struct {
	BeaconStatus  map[beacon.ValidatorPubkey]types.ValidatorState `json:"beaconStatus"`
	NodesetStatus map[beacon.ValidatorPubkey]NodesetStatus        `json:"nodesetStatus"`
}
