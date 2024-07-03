package api

import "github.com/rocket-pool/node-manager-core/beacon"

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

type NodeSetStakeWise_GetRegisteredValidatorsData struct {
	Validators []ValidatorStatus `json:"validators"`
}

type NodeSetStakeWise_GetDepositDataSetData struct {
	Version     int                          `json:"version"`
	DepositData []beacon.ExtendedDepositData `json:"depositData"`
}
