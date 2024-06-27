package api

type NodesetStatus string

const (
	NodesetStatus_Unknown               NodesetStatus = ""
	NodesetStatus_RegisteredToStakewise NodesetStatus = "RegisteredToStakewise"
	NodesetStatus_UploadedStakewise     NodesetStatus = "UploadedStakewise"
	NodesetStatus_UploadedToNodeset     NodesetStatus = "UploadedToNodeset"
	NodesetStatus_Generated             NodesetStatus = "Generated"
)

type NodeSetRegistrationStatus string

const (
	NodeSetRegistrationStatus_Registered   NodeSetRegistrationStatus = "registered"
	NodeSetRegistrationStatus_Unregistered NodeSetRegistrationStatus = "unregistered"
	NodeSetRegistrationStatus_Unknown      NodeSetRegistrationStatus = "unknown"
	NodeSetRegistrationStatus_NoWallet     NodeSetRegistrationStatus = "no-wallet"
)

type NodeSetRegisterNodeData struct {
	Success           bool `json:"success"`
	AlreadyRegistered bool `json:"alreadyRegistered"`
	NotWhitelisted    bool `json:"notWhitelisted"`
}

type NodeSetRegistrationStatusData struct {
	Status       NodeSetRegistrationStatus `json:"status"`
	ErrorMessage string                    `json:"errorMessage"`
}
