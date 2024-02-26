package swapi

import "github.com/nodeset-org/eth-utils/beacon"

type NodesetUploadDepositDataData struct {
	ServerResponse []byte                   `json:"serverResponse"`
	NewPubkeys     []beacon.ValidatorPubkey `json:"newPubkeys"`
	TotalCount     uint64                   `json:"totalCount"`
}

type NodesetGetValidatorsData struct {
}
