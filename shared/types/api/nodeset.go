package api

// The registration status of the node with the NodeSet server
type NodeSetRegistrationStatus string

const (
	// The node has been registered with a user account on the NodeSet server
	NodeSetRegistrationStatus_Registered NodeSetRegistrationStatus = "registered"

	// The node has not been registered with a user account on the NodeSet server
	NodeSetRegistrationStatus_Unregistered NodeSetRegistrationStatus = "unregistered"

	// The node's registration status is unknown
	NodeSetRegistrationStatus_Unknown NodeSetRegistrationStatus = "unknown"

	// The node has no wallet yet
	NodeSetRegistrationStatus_NoWallet NodeSetRegistrationStatus = "no-wallet"
)

type NodeSetRegisterNodeData struct {
	Success           bool `json:"success"`
	AlreadyRegistered bool `json:"alreadyRegistered"`
	NotWhitelisted    bool `json:"notWhitelisted"`
}

type NodeSetGetRegistrationStatusData struct {
	Status       NodeSetRegistrationStatus `json:"status"`
	ErrorMessage string                    `json:"errorMessage"`
}
