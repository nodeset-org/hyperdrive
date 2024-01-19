package api

import "github.com/ethereum/go-ethereum/common"

// This is a wrapper for the EC status report
type ClientStatus struct {
	IsWorking    bool    `json:"isWorking"`
	IsSynced     bool    `json:"isSynced"`
	SyncProgress float64 `json:"syncProgress"`
	NetworkId    uint    `json:"networkId"`
	Error        string  `json:"error"`
}

// This is a wrapper for the manager's overall status report
type ClientManagerStatus struct {
	PrimaryClientStatus  ClientStatus `json:"primaryEcStatus"`
	FallbackEnabled      bool         `json:"fallbackEnabled"`
	FallbackClientStatus ClientStatus `json:"fallbackEcStatus"`
}

type ConfirmNodePrimaryWithdrawalAddressResponse struct {
	Status string      `json:"status"`
	Error  string      `json:"error"`
	TxHash common.Hash `json:"txHash"`
}
