package swapi

import "github.com/nodeset-org/eth-utils/beacon"

type ValidatorExitInfo struct {
	Pubkey    beacon.ValidatorPubkey    `json:"pubkey"`
	Index     uint64                    `json:"index"`
	Signature beacon.ValidatorSignature `json:"signature"`
}

type ValidatorExitData struct {
	Epoch     uint64              `json:"epoch"`
	ExitInfos []ValidatorExitInfo `json:"exitInfos"`
}
