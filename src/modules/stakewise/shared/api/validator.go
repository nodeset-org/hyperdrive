package swapi

import (
	"github.com/rocket-pool/node-manager-core/beacon"
)

type ValidatorExitInfo struct {
	Index     uint64                    `json:"index"`
	Signature beacon.ValidatorSignature `json:"signature"`
}

type ValidatorGetSignedExitMessagesData struct {
	Epoch     uint64                       `json:"epoch"`
	ExitInfos map[string]ValidatorExitInfo `json:"exitInfos"` // map[beacon.ValidatorPubkey]ValidatorExitInfo
}
