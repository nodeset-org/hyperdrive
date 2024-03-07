package swapi

import (
	nmc_beacon "github.com/rocket-pool/node-manager-core/beacon"
)

type ValidatorExitInfo struct {
	Index     uint64                        `json:"index"`
	Signature nmc_beacon.ValidatorSignature `json:"signature"`
}

type ValidatorGetSignedExitMessagesData struct {
	Epoch     uint64                       `json:"epoch"`
	ExitInfos map[string]ValidatorExitInfo `json:"exitInfos"` // map[beacon.ValidatorPubkey]ValidatorExitInfo
}
