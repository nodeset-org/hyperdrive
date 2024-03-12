package swapi

import (
	"github.com/rocket-pool/node-manager-core/beacon"
)

type NodesetStatus string

const (
	RegisteredToStakewise NodesetStatus = "RegisteredToStakewise"
	UploadedStakewise     NodesetStatus = "UploadedStakewise"
	UploadedToNodeset     NodesetStatus = "UploadedToNodeset"
	Generated             NodesetStatus = "Generated"
)

type ValidatorStatusData struct {
	BeaconStatus  map[beacon.ValidatorPubkey]beacon.ValidatorState `json:"beaconStatus"`
	NodesetStatus map[beacon.ValidatorPubkey]NodesetStatus         `json:"nodesetStatus"`
}
