package api

import "github.com/rocket-pool/node-manager-core/beacon"

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

// Details of an exit message
type ExitMessageDetails struct {
	Epoch          string `json:"epoch"`
	ValidatorIndex string `json:"validator_index"`
}

// Voluntary exit message
type ExitMessage struct {
	Message   ExitMessageDetails `json:"message"`
	Signature string             `json:"signature"`
}

// Data for a pubkey's voluntary exit message
type ExitData struct {
	Pubkey      string      `json:"pubkey"`
	ExitMessage ExitMessage `json:"exit_message"`
}

// Validator status info
type ValidatorStatus struct {
	Pubkey              beacon.ValidatorPubkey `json:"pubkey"`
	Status              string                 `json:"status"`
	ExitMessageUploaded bool                   `json:"exitMessage"`
}
