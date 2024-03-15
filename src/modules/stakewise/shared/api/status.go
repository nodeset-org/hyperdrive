package swapi

import (
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
	BeaconStatus  map[string]types.ValidatorState `json:"beaconStatus"`  // string => beacon.ValidatorPubkey
	NodesetStatus map[string]NodesetStatus        `json:"nodesetStatus"` // string => beacon.ValidatorPubkey
}
