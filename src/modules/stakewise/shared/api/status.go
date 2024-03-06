package swapi

type BeaconStatus string
type NodesetStatus string

const (
	WithdrawalDone     BeaconStatus = "WithdrawalDone"
	WithdrawalPossible BeaconStatus = "WithdrawalPossible"
	ExitedSlashed      BeaconStatus = "ExitedSlashed"
	ExitedUnslashed    BeaconStatus = "ExitedUnslashed"
	ActiveSlashed      BeaconStatus = "ActiveSlashed"
	ActiveExited       BeaconStatus = "ActiveExited"
	ActiveOngoing      BeaconStatus = "ActiveOngoing"
	PendingQueued      BeaconStatus = "PendingQueued"
	PendingInitialized BeaconStatus = "PendingInitialized"
	NotAvailable       BeaconStatus = "NotAvailable"
)

const (
	RegisteredToStakewise NodesetStatus = "RegisteredToStakewise"
	UploadedStakewise     NodesetStatus = "UploadedStakewise"
	UploadedToNodeset     NodesetStatus = "UploadedToNodeset"
	Generated             NodesetStatus = "Generated"
)

type ValidatorStatusData struct {
	BeaconStatus  map[string]BeaconStatus  `json:"beaconStatus"`  // string => beacon.ValidatorPubkey
	NodesetStatus map[string]NodesetStatus `json:"nodesetStatus"` // string => beacon.ValidatorPubkey
}
