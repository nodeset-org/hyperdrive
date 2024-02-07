package swapi

import "github.com/nodeset-org/eth-utils/beacon"

type ValidatorExitInfo struct {
	Index     uint64                    `json:"index"`
	Signature beacon.ValidatorSignature `json:"signature"`
}

type ValidatorGetSignedExitMessagesData struct {
	Epoch     uint64                                       `json:"epoch"`
	ExitInfos map[beacon.ValidatorPubkey]ValidatorExitInfo `json:"exitInfos"`
}
